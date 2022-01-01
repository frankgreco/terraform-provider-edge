package provider

import (
	"github.com/frankgreco/terraform-attribute-validators"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceFirewallRulesetSchema(isDataSource bool) tfsdk.Schema {
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
				Optional:    !isDataSource,
				Computed:    isDataSource,
			},
			"default_action": {
				Description: "The default action to take if traffic is not matched by one of the rules in the ruleset. Must be one of `reject`, `drop`, `accept`.",
				Type:        tftypes.StringType,
				Required:    !isDataSource,
				Computed:    isDataSource,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSlice(true, "reject", "drop", "accept"),
				},
			},
		},
		Blocks: map[string]tfsdk.Block{
			"rule": {
				NestingMode: tfsdk.BlockNestingModeSet,
				Attributes: map[string]tfsdk.Attribute{
					"priority": {
						Type:        tftypes.NumberType,
						Required:    !isDataSource,
						Computed:    isDataSource,
						Description: "The priority of this rule. The higher the priority, the higher the precedence.",
					},
					"description": {
						Type:        tftypes.StringType,
						Optional:    !isDataSource,
						Computed:    isDataSource,
						Description: "A human readable description for this rule.",
					},
					"action": {
						Type:        tftypes.StringType,
						Required:    !isDataSource,
						Computed:    isDataSource,
						Description: "The action to take on traffic that matches this rule. Must be one of `reject`, `drop`, `accept`.",
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, "drop", "reject", "accept"),
						},
					},
					"protocol": {
						Type:        tftypes.StringType,
						Optional:    !isDataSource,
						Computed:    isDataSource,
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
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "Match packets that are part of a two-way connection.",
							},
							"new": {
								Type:        tftypes.BoolType,
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "Match packets creating a new connection.",
							},
							"related": {
								Type:        tftypes.BoolType,
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "Match packets related to established connections.",
							},
							"invalid": {
								Type:        tftypes.BoolType,
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "Match packets that cannot be identified.",
							},
						}),
						Optional: !isDataSource,
						Computed: isDataSource,
					},
					"destination": {
						Description: "Details about the traffic's destination. If not specified, all sources will be evaluated.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"from_port": {
								Type:        tftypes.NumberType,
								Required:    !isDataSource,
								Computed:    isDataSource,
								Description: "The first destination port in the port range this rule will apply to.",
								Validators: []tfsdk.AttributeValidator{
									validators.Range(float64(1), float64(65535.0)),
								},
							},
							"to_port": {
								Type:        tftypes.NumberType,
								Required:    !isDataSource,
								Computed:    isDataSource,
								Description: "The first destination port in the port range this rule will apply to. If only one port is desired, set to the same value in `from_port`.",
								Validators: []tfsdk.AttributeValidator{
									validators.Range(float64(1), float64(65535.0)),
								},
							},
							"address": {
								Type:        tftypes.StringType,
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "The cidr this rule applies to. If not provided, it is treated as 0.0.0.0/0.",
								Validators: []tfsdk.AttributeValidator{
									validators.Cidr(),
								},
							},
						}),
						Optional: !isDataSource,
						Computed: isDataSource,
					},
					"source": {
						Description: "Details about the traffic's source. If not specified, all sources will be evaluated.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"from_port": {
								Type:        tftypes.NumberType,
								Required:    !isDataSource,
								Computed:    isDataSource,
								Description: "The first destination port in the port range this rule will apply to.",
								Validators: []tfsdk.AttributeValidator{
									validators.Range(float64(1), float64(65535.0)),
								},
							},
							"to_port": {
								Type:        tftypes.NumberType,
								Required:    !isDataSource,
								Computed:    isDataSource,
								Description: "The first destination port in the port range this rule will apply to. If only one port is desired, set to the same value in `from_port`.",
								Validators: []tfsdk.AttributeValidator{
									validators.Range(float64(1), float64(65535.0)),
								},
							},
							"address": {
								Type:        tftypes.StringType,
								Optional:    !isDataSource,
								Computed:    isDataSource,
								Description: "The cidr this rule applies to. If not provided, it is treated as `0.0.0.0/0`.",
								Validators: []tfsdk.AttributeValidator{
									validators.Cidr(),
								},
							},
							"mac": {
								Type:     tftypes.StringType,
								Optional: !isDataSource,
								Computed: isDataSource,
							},
						}),
						Optional: !isDataSource,
						Computed: isDataSource,
					},
				},
			},
		},
	}
}
