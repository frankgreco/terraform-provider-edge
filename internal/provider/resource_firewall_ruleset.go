package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/frankgreco/edge-sdk-go/types"
	"github.com/frankgreco/terraform-attribute-validators"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mattbaird/jsonpatch"
)

type resourceFirewallRulesetType struct{}

func (r resourceFirewallRulesetType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	port := tfsdk.Attribute{
		Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
			"from": {
				Type:     tftypes.NumberType,
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					validators.Range(float64(1), float64(65535.0)),
				},
			},
			"to": {
				Type:     tftypes.NumberType,
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					validators.Range(float64(1), float64(65535.0)),
				},
			},
		}),
		Optional:    true,
		Description: "A port range.",
	}

	address := tfsdk.Attribute{
		Type:        tftypes.StringType,
		Optional:    true,
		Description: "The cidr this rule applies to. If not provided, it is treated as `0.0.0.0/0`.",
		Validators: []tfsdk.AttributeValidator{
			validators.Cidr(),
		},
	}

	addressGroup := tfsdk.Attribute{
		Type:        tftypes.StringType,
		Optional:    true,
		Description: "The address group this rule applies to. If not provided, all addresses will be matched.",
	}

	portGroup := tfsdk.Attribute{
		Type:        tftypes.StringType,
		Optional:    true,
		Description: "The port group this rule applies to. If not provided, all ports will be matched.",
	}

	return tfsdk.Schema{
		Description: "A grouping of firewall rules. The firewall is not enforced unless attached to an interface which can be done with the `firewall_ruleset_attachment` resource.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "A unique, human readable name for this ruleset.",
				Type:        tftypes.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.NoWhitespace(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
			},
			"description": {
				Description: "A human readable description for this ruleset.",
				Type:        tftypes.StringType,
				Optional:    true,
			},
			"default_action": {
				Description: "The default action to take if traffic is not matched by one of the rules in the ruleset. Must be one of `reject`, `drop`, `accept`.",
				Type:        tftypes.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSlice(true, "reject", "drop", "accept"),
				},
			},
		},
		Blocks: map[string]tfsdk.Block{
			"rule": {
				Validators: []tfsdk.AttributeValidator{
					validators.Unique("priority"),
				},
				NestingMode: tfsdk.BlockNestingModeSet,
				Attributes: map[string]tfsdk.Attribute{
					"priority": {
						Type:        tftypes.NumberType,
						Required:    true,
						Description: "The priority of this rule. The higher the priority, the higher the precedence.",
					},
					"description": {
						Type:        tftypes.StringType,
						Optional:    true,
						Description: "A human readable description for this rule.",
					},
					"action": {
						Type:        tftypes.StringType,
						Required:    true,
						Description: "The action to take on traffic that matches this rule. Must be one of `reject`, `drop`, `accept`.",
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, "drop", "reject", "accept"),
						},
					},
					"protocol": {
						Type:        tftypes.StringType,
						Optional:    true,
						Description: "The protocol this rule applies to. If not specified, this rule applies to all protcols. Must be one of `tcp`, `udp`, `tcp_udp`.",
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, "tcp", "udp", "tcp_udp"),
						},
					},
					"state": {
						Description: "This describes the connection state of a packet.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"established": {
								Type:        tftypes.BoolType,
								Optional:    true,
								Description: "Match packets that are part of a two-way connection.",
							},
							"new": {
								Type:        tftypes.BoolType,
								Optional:    true,
								Description: "Match packets creating a new connection.",
							},
							"related": {
								Type:        tftypes.BoolType,
								Optional:    true,
								Description: "Match packets related to established connections.",
							},
							"invalid": {
								Type:        tftypes.BoolType,
								Optional:    true,
								Description: "Match packets that cannot be identified.",
							},
						}),
						Optional: true,
					},
					"destination": {
						Description: "Details about the traffic's destination. If not specified, all sources will be evaluated.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"address":       address,
							"port":          port,
							"address_group": addressGroup,
							"port_group":    portGroup,
						}),
						Optional: true,
						// Need a validator to ensure address conflicts with address_group and port conflicts with port_group.
					},
					"source": {
						Description: "Details about the traffic's source. If not specified, all sources will be evaluated.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"address":       address,
							"port":          port,
							"address_group": addressGroup,
							"port_group":    portGroup,
							"mac": {
								Type:     tftypes.StringType,
								Optional: true,
							},
						}),
						Optional: true,
						// Need a validator to ensure address conflicts with address_group and port conflicts with port_group.
					},
				},
			},
		},
	}, nil
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
