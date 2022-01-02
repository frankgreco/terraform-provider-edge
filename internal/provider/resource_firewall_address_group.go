package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tfftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

type resourceFirewallAddressGroupType struct{}

func (r resourceFirewallAddressGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "A logical grouping of addresses.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:          tfftypes.StringType,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Description:   "A unique, human readable name for this address group.",
			},
			"description": {
				Type:        tfftypes.StringType,
				Optional:    true,
				Description: "A human readable description for this address group.",
			},
			"cidrs": {
				Type:        tfftypes.ListType{ElemType: tfftypes.StringType},
				Optional:    true,
				Description: "A non-overlapping list of cidrs.",
			},
		},
	}, nil
}

// New resource instance
func (r resourceFirewallAddressGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFirewallAddressGroup{
		p: *(p.(*provider)),
	}, nil
}

type resourceFirewallAddressGroup struct {
	p provider
}

func (r resourceFirewallAddressGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var plan types.AddressGroup
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	actual, err := r.p.client.Firewall.CreateAddressGroup(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue creating the address group.",
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

func (r resourceFirewallAddressGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.AddressGroup
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	group, err := r.p.client.Firewall.GetAddressGroup(ctx, state.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the address group.",
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

func (r resourceFirewallAddressGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var current types.AddressGroup
	{
		diagnostics := req.State.Get(ctx, &current)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}

		log.Printf("[TRACE] current firewall address group struct: %+v", current)
	}

	var desired types.AddressGroup
	{
		diags := req.Plan.Get(ctx, &desired)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		log.Printf("[TRACE] desired firewall address group struct: %+v", desired)
	}

	var patches []jsonpatch.JsonPatchOperation
	{
		cData, err := json.Marshal(&current)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall address group from state",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] current firewall address group json: %s", string(cData))

		dData, err := json.Marshal(&desired)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall address group from plan",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] desired firewall address group json: %s", string(dData))

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

	updated, err := r.p.client.Firewall.UpdateAddressGroup(ctx, &current, patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue updating the firewall address group.",
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

func (r resourceFirewallAddressGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.AddressGroup
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if err := r.p.client.Firewall.DeleteAddressGroup(ctx, state.Name); err != nil {
		resp.Diagnostics.AddError(
			"There was an issue deleting the address group.",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceFirewallAddressGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	group, err := r.p.client.Firewall.GetAddressGroup(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the address group.",
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
