package provider

import (
	"context"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/frankgreco/terraform-helpers/validators"
	"github.com/mattbaird/jsonpatch"

	"terraform-provider-edge/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tfftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceFirewallPortGroupType struct{}

func (r resourceFirewallPortGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "A logical grouping of ports.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:          tfftypes.StringType,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Description:   "A unique, human readable name for this port group.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoWhitespace(),
				},
			},
			"description": {
				Type:        tfftypes.StringType,
				Optional:    true,
				Description: "A human readable description for this port group.",
			},
			"port_ranges": {
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"from": {
						Type:     tfftypes.NumberType,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							validators.Range(float64(1), float64(65535.0)),
							validators.Compare(validators.ComparatorLessThan, "to"),
						},
					},
					"to": {
						Type:     tfftypes.NumberType,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							validators.Range(float64(1), float64(65535.0)),
						},
					},
				}, tfsdk.ListNestedAttributesOptions{}),
				Optional:    true,
				Description: "A list of port ranges.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoOverlap(),
				},
			},
			"ports": {
				Type:        tfftypes.ListType{ElemType: tfftypes.NumberType},
				Optional:    true,
				Description: "A list of port numbers.",
				Validators:  []tfsdk.AttributeValidator{
					validators.NoOverlap(),
				},
			},
		},
	}, nil
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

func (r resourceFirewallPortGroup) Update(ctx context.Context, current interface{}, patches []jsonpatch.JsonPatchOperation) (interface{}, error) {
	return r.p.client.Firewall.UpdatePortGroup(ctx, current.(*types.PortGroup), patches)
}

func (r resourceFirewallPortGroup) Delete(ctx context.Context, id string) error {
	return r.p.client.Firewall.DeletePortGroup(ctx, id)
}
