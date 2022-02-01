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
	ruleset := plan.(types.Ruleset)
	return r.p.client.Firewall.CreateRuleset(ctx, &ruleset)
}

func (r resourceFirewallRuleset) Update(ctx context.Context, current interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	ruleset := current.(types.Ruleset)
	return r.p.client.Firewall.UpdateRuleset(ctx, &ruleset, patches)
}

func (r resourceFirewallRuleset) Delete(ctx context.Context, id string) error {
	return r.p.client.Firewall.DeleteRuleset(ctx, id)
}
