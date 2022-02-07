package provider

import (
	"context"

	"terraform-provider-edge/internal/utils"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type resourceFirewallRulesetType struct{}

func (r resourceFirewallRulesetType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemaFirewallRuleset(), nil
}

func (r resourceFirewallRulesetType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return utils.Resource{
		Name:         "firewall ruleset",
		Attribute:    "name",
		IsConfigured: (p.(*provider)).configured,
		Api:          resourceFirewallRuleset{p: *(p.(*provider))},
		Type:         types.Ruleset{},
	}, nil
}

type resourceFirewallRuleset struct {
	p provider
}

func (r resourceFirewallRuleset) Read(ctx context.Context, id string) (interface{}, error) {
	return r.p.client.Firewall.GetRuleset(ctx, id)
}

func (r resourceFirewallRuleset) Create(ctx context.Context, plan interface{}) (interface{}, error) {
	desiredRuleset := plan.(types.Ruleset)

	createdRuleset, err := r.p.client.Firewall.CreateRuleset(ctx, &desiredRuleset)
	if err != nil {
		return nil, err
	}

	return normalize(&desiredRuleset, createdRuleset), nil
}

func (r resourceFirewallRuleset) Update(ctx context.Context, current, desired interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	currentRuleset := current.(types.Ruleset)
	desiredRuleset := desired.(types.Ruleset)

	updatedRuleset, err := r.p.client.Firewall.UpdateRuleset(ctx, &currentRuleset, patches)
	if err != nil {
		return nil, err
	}

	return normalize(&desiredRuleset, updatedRuleset), nil
}

func (r resourceFirewallRuleset) Delete(ctx context.Context, id string) error {
	return r.p.client.Firewall.DeleteRuleset(ctx, id)
}

// We need to tell the difference between null and false for certain booleans.
// If we don't, we'll get errors of the following type:
//
//	providerproduced an unexpected new value: .default_logging: was null, but now cty.False.
//
// To fix this, if we get a null or false value back and the desired state was either
// null or false, we're going to set the state to whatever was in the plan.
func normalize(desired, actual *types.Ruleset) *types.Ruleset {
	// ruleset.default_logging
	if l := actual.DefaultLogging; l == nil || !*l {
		actual.DefaultLogging = desired.DefaultLogging
	}

	// ruleset.rules[*].log
	desiredIndexed := indexPriority(desired.Rules)

	for _, rule := range actual.Rules {
		if l := rule.Log; l == nil || !*l {
			rule.Log = desired.Rules[desiredIndexed[rule.Priority]].Log
		}
	}

	return actual
}

func indexPriority(rules []*types.Rule) map[int]int {
	m := map[int]int{}

	for i, rule := range rules {
		m[rule.Priority] = i
	}

	return m
}
