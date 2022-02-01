package provider

import (
	"context"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"terraform-provider-edge/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type resourceFirewallRulesetAttachmentType struct{}

func (r resourceFirewallRulesetAttachmentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemaFirewallRulesetAttachment(), nil
}

func (r resourceFirewallRulesetAttachmentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return utils.Resource{
		Name:         "firewall ruleset attachment",
		Attribute:    "interface",
		IsConfigured: (p.(*provider)).configured,
		Api:          resourceFirewallRulesetAttachment{p: *(p.(*provider))},
		Type:         types.FirewallAttachment{},
	}, nil
}

type resourceFirewallRulesetAttachment struct {
	p provider
}

func (r resourceFirewallRulesetAttachment) Read(ctx context.Context, id string) (interface{}, error) {
	return r.p.client.Interfaces.Ethernet.GetFirewallRulesetAttachment(ctx, id)
}

func (r resourceFirewallRulesetAttachment) Create(ctx context.Context, plan interface{}) (interface{}, error) {
	attachment := plan.(types.FirewallAttachment)
	return r.p.client.Interfaces.Ethernet.AttachFirewallRuleset(ctx, attachment.Interface, &attachment)
}

func (r resourceFirewallRulesetAttachment) Update(ctx context.Context, current interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	attachment := current.(types.FirewallAttachment)
	return r.p.client.Interfaces.Ethernet.UpdateFirewallRulesetAttachment(ctx, &attachment, patches)
}

func (r resourceFirewallRulesetAttachment) Delete(ctx context.Context, id string) error {
	return r.p.client.Interfaces.Ethernet.DetachFirewallRuleset(ctx, id)
}
