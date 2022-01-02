package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type resourceFirewallRulesetType struct{}

func (r resourceFirewallRulesetType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return resourceFirewallRulesetSchema(false), nil
}

func (r resourceFirewallRulesetType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFirewallRuleset{
		p: *(p.(*provider)),
	}, nil
}

type resourceFirewallRuleset struct {
	p provider
}

func (r resourceFirewallRuleset) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var plan types.Ruleset
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	actual, err := r.p.client.Firewall.CreateRuleset(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue creating the ruleset.",
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

func (r resourceFirewallRuleset) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.Ruleset
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	ruleset, err := r.p.client.Firewall.GetRuleset(ctx, state.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the ruleset.",
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *ruleset)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceFirewallRuleset) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var current types.Ruleset
	{
		diagnostics := req.State.Get(ctx, &current)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}

		(&current).SetCodecMode(types.CodecModeLocal)
		log.Printf("[TRACE] current ruleset struct: %+v", current)
	}

	var desired types.Ruleset
	{
		diags := req.Plan.Get(ctx, &desired)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		(&desired).SetCodecMode(types.CodecModeLocal)
		log.Printf("[TRACE] desired ruleset struct: %+v", desired)
	}

	var patches []jsonpatch.JsonPatchOperation
	{
		cData, err := json.Marshal(&current)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal ruleset from state",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] current ruleset json: %s", string(cData))

		dData, err := json.Marshal(&desired)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal ruleset from plan",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] desired ruleset json: %s", string(dData))

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

	updated, err := r.p.client.Firewall.UpdateRuleset(ctx, &current, patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue updating the ruleset.",
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

func (r resourceFirewallRuleset) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var state types.Ruleset
	{
		diagnostics := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if err := r.p.client.Firewall.DeleteRuleset(ctx, state.Name); err != nil {
		resp.Diagnostics.AddError(
			"There was an issue deleting the ruleset.",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceFirewallRuleset) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	ruleset, err := r.p.client.Firewall.GetRuleset(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue reading the ruleset.",
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *ruleset)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}
