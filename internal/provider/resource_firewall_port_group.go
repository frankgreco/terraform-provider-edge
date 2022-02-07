package provider

import (
	"context"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"terraform-provider-edge/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type resourceFirewallPortGroupType struct{}

func (r resourceFirewallPortGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemaFirewallPortGroup(), nil
}

func (r resourceFirewallPortGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return utils.Resource{
		Name:         "firewall port group",
		Attribute:    "name",
		IsConfigured: (p.(*provider)).configured,
		Api:          resourceFirewallPortGroup{p: *(p.(*provider))},
		Type:         types.PortGroup{},
	}, nil
}

type resourceFirewallPortGroup struct {
	p provider
}

func (r resourceFirewallPortGroup) Read(ctx context.Context, id string) (interface{}, error) {
	return r.p.client.Firewall.GetPortGroup(ctx, id)
}

func (r resourceFirewallPortGroup) Create(ctx context.Context, plan interface{}) (interface{}, error) {
	group := plan.(types.PortGroup)
	return r.p.client.Firewall.CreatePortGroup(ctx, &group)
}

func (r resourceFirewallPortGroup) Update(ctx context.Context, current, desired interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	group := current.(types.PortGroup)
	return r.p.client.Firewall.UpdatePortGroup(ctx, &group, patches)
}

func (r resourceFirewallPortGroup) Delete(ctx context.Context, id string) error {
	return r.p.client.Firewall.DeletePortGroup(ctx, id)
}
