package provider

import (
	"context"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"terraform-provider-edge/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type resourceFirewallAddressGroupType struct{}

func (r resourceFirewallAddressGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemaFirewallAddressGroup(), nil
}

func (r resourceFirewallAddressGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return utils.Resource{
		Name:         "firewall address group",
		Attribute:    "name",
		IsConfigured: (p.(*provider)).configured,
		Api:          resourceFirewallAddressGroup{p: *(p.(*provider))},
		Type:         types.AddressGroup{},
	}, nil
}

type resourceFirewallAddressGroup struct {
	p provider
}

func (r resourceFirewallAddressGroup) Read(ctx context.Context, id string) (interface{}, error) {
	return r.p.client.Firewall.GetAddressGroup(ctx, id)
}

func (r resourceFirewallAddressGroup) Create(ctx context.Context, plan interface{}) (interface{}, error) {
	group := plan.(types.AddressGroup)
	return r.p.client.Firewall.CreateAddressGroup(ctx, &group)
}

func (r resourceFirewallAddressGroup) Update(ctx context.Context, current interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	return r.p.client.Firewall.UpdateAddressGroup(ctx, current.(*types.AddressGroup), patches)
}

func (r resourceFirewallAddressGroup) Delete(ctx context.Context, id string) error {
	return r.p.client.Firewall.DeleteAddressGroup(ctx, id)
}
