package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/frankgreco/terraform-attribute-validators"
	"github.com/mattbaird/jsonpatch"

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
				// Need a validator that ensures from != to.
			},
			"ports": {
				Type:        tfftypes.ListType{ElemType: tfftypes.NumberType},
				Optional:    true,
				Description: "A list of port numbers.",
				// Need a validator that ensures no duplicates.
			},
			// Need a validator that ensures no duplicates exist.
		},
	}, nil
}

// New resource instance
func (r resourceFirewallPortGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFirewallPortGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourceFirewallPortGroup struct {
	p provider
}

func (r resourceFirewallPortGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var plan types.PortGroup
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	actual, err := r.p.client.Firewall.CreatePortGroup(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue creating the port group.",
			err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, *actual)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceFirewallPortGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.PortGroup
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	group, err := r.p.client.Firewall.GetPortGroup(ctx, state.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the port group.",
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *group)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceFirewallPortGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var current types.PortGroup
	{
		diagnostics := req.State.Get(ctx, &current)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}

		log.Printf("[TRACE] current firewall port group struct: %+v", current)
	}

	var desired types.PortGroup
	{
		diags := req.Plan.Get(ctx, &desired)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		log.Printf("[TRACE] desired firewall port group struct: %+v", desired)
	}

	var patches []jsonpatch.JsonPatchOperation
	{
		cData, err := json.Marshal(&current)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall port group from state",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] current firewall port group json: %s", string(cData))

		dData, err := json.Marshal(&desired)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall port group from plan",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] desired firewall port group json: %s", string(dData))

		p, err := jsonpatch.CreatePatch(cData, dData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not create patch document.",
				err.Error(),
			)
			return
		}
		patches = p
		log.Printf("[DEBUG] patch document: %+v", patches)
	}

	updated, err := r.p.client.Firewall.UpdatePortGroup(ctx, &current, patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue updating the firewall port group.",
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *updated)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceFirewallPortGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.PortGroup
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if err := r.p.client.Firewall.DeletePortGroup(ctx, state.Name); err != nil {
		resp.Diagnostics.AddError(
			"There was an issue deleting the port group.",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceFirewallPortGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	group, err := r.p.client.Firewall.GetPortGroup(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the port group.",
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *group)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}
