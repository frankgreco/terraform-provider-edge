package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type dataSourceFirewallRulesetType struct{}

func (r dataSourceFirewallRulesetType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Computed: true,
				Optional: true,
			},
			"default_action": {
				Type:     types.StringType,
				Computed: true,
			},
		},
		Blocks: map[string]tfsdk.Block{
			"rule": {
				NestingMode: tfsdk.BlockNestingModeSet,
				Attributes: map[string]tfsdk.Attribute{
					"priority": {
						Type:     types.NumberType,
						Computed: true,
					},
					"description": {
						Type:     types.StringType,
						Computed: true,
					},
					"action": {
						Type:     types.StringType,
						Computed: true,
					},
					"protocol": {
						Type:     types.StringType,
						Computed: true,
					},
					"state": {
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"established": {
								Type:     types.BoolType,
								Optional: true,
							},
							"new": {
								Type:     types.BoolType,
								Optional: true,
							},
							"related": {
								Type:     types.BoolType,
								Optional: true,
							},
							"invalid": {
								Type:     types.BoolType,
								Optional: true,
							},
						}),
						Computed: true,
						Optional: true,
					},
					"destination": {
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"from_port": {
								Type:     types.NumberType,
								Required: true,
							},
							"to_port": {
								Type:     types.NumberType,
								Required: true,
							},
							"address": {
								Type:     types.StringType,
								Optional: true,
							},
						}),
						Computed: true,
						Optional: true,
					},
					"source": {
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"from_port": {
								Type:     types.NumberType,
								Required: true,
							},
							"to_port": {
								Type:     types.NumberType,
								Required: true,
							},
							"address": {
								Type:     types.StringType,
								Optional: true,
							},
							"mac": {
								Type:     types.StringType,
								Optional: true,
							},
						}),
						Computed: true,
						Optional: true,
					},
				},
			},
		},
	}, nil
}

func (r dataSourceFirewallRulesetType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceFirewallRuleset{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceFirewallRuleset struct {
	p provider
}

func (r dataSourceFirewallRuleset) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"The provider has not been configured!",
			"Please configure the provider.",
		)
		return
	}

	var name string
	{
		diagnostics := req.Config.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("name"), &name)
		resp.Diagnostics.Append(diagnostics...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	ruleset, err := r.p.client.Firewall.GetRuleset(ctx, name)
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
