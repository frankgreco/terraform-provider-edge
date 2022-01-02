package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/mattbaird/jsonpatch"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tfftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type resourceFirewallRulesetAttachmentType struct{}

func (r resourceFirewallRulesetAttachmentType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Attach a firewall ruleset to inbound, outbound, and local traffic.",
		Attributes: map[string]tfsdk.Attribute{
			"interface": {
				Type:          tfftypes.StringType,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Description:   "The interface to attach firewall rules to.",
			},
			"in": {
				Type:        tfftypes.StringType,
				Optional:    true,
				Description: "Match inbound packets.",
			},
			"out": {
				Type:        tfftypes.StringType,
				Optional:    true,
				Description: "Match outbound packets.",
			},
			"local": {
				Type:        tfftypes.StringType,
				Optional:    true,
				Description: "Match local packets.",
			},
		},
	}, nil
}

// New resource instance
func (r resourceFirewallRulesetAttachmentType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFirewallRulesetAttachment{
		p: *(p.(*provider)),
	}, nil
}

type resourceFirewallRulesetAttachment struct {
	p provider
}

// Create a new resource
func (r resourceFirewallRulesetAttachment) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var plan types.FirewallAttachment
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if _, err := r.p.client.Interfaces.Ethernet.AttachFirewallRuleset(ctx, plan.Interface, &plan); err != nil {
		resp.Diagnostics.AddError(
			"There was an issue creating the attachment.",
			err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceFirewallRulesetAttachment) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var id string
	{
		diagnostics := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("interface"), &id)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	ethernet, err := r.p.client.Interfaces.Ethernet.Get(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue retrieving the ethernet interface %s.", id),
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.Set(ctx, *ethernet.Firewall)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceFirewallRulesetAttachment) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var current types.FirewallAttachment
	{
		diagnostics := req.State.Get(ctx, &current)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}

		log.Printf("[TRACE] current firewall attachment struct: %+v", current)
	}

	var desired types.FirewallAttachment
	{
		diags := req.Plan.Get(ctx, &desired)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		log.Printf("[TRACE] desired firewall attachment struct: %+v", desired)
	}

	var patches []jsonpatch.JsonPatchOperation
	{
		cData, err := json.Marshal(&current)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall attachment from state",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] current firewall attachment json: %s", string(cData))

		dData, err := json.Marshal(&desired)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not marshal firewall attachment from plan",
				err.Error(),
			)
			return
		}
		log.Printf("[TRACE] desired firewall attachment json: %s", string(dData))

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

	updated, err := r.p.client.Interfaces.Ethernet.UpdateFirewallRulesetAttachment(ctx, &current, patches)
	if err != nil {
		resp.Diagnostics.AddError(
			"There was an issue updating the firewall attachment.",
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

// Delete resource
func (r resourceFirewallRulesetAttachment) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var id string
	{
		diagnostics := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("interface"), &id)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if err := r.p.client.Interfaces.Ethernet.DetachFirewallRuleset(ctx, id); err != nil {
		resp.Diagnostics.AddError(
			"There was an issue deleting the ruleset.",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceFirewallRulesetAttachment) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	ethernet, err := r.p.client.Interfaces.Ethernet.Get(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("There was an issue retrieving the ethernet interface %s.", req.ID),
			err.Error(),
		)
		return
	}

	diagnostics := resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("interface"), req.ID)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagnostics = resp.State.Set(ctx, *ethernet.Firewall)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
}
