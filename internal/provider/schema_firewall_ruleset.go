package provider

import (
	"github.com/frankgreco/terraform-helpers/validators"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-edge/internal/utils"
)

var (
	protocols = []string{
		"tcp",
		"udp",
		"tcp_udp",
		"ah",
		"ax.25",
		"dccp",
		"ddp",
		"egp",
		"eigrp",
		"encap",
		"esp",
		"etherip",
		"fc",
		"ggp",
		"gre",
		"hip",
		"hmp",
		"icmp",
		"idpr-cmtp",
		"idpr",
		"igmp",
		"igp",
		"ip",
		"ipcomp",
		"ipencap",
		"ipip",
		"isis",
		"iso-tp4",
		"l2tp",
		"manet",
		"mpls-in-ip",
		"ospf",
		"pim",
		"pup",
		"rdp",
		"rohc",
		"rspf",
		"rsvp",
		"sctp",
		"skip",
		"st",
		"udplite",
		"vmtp",
		"vrrp",
		"wesp",
		"xns-idp",
		"xtp",
	}
)

func schemaFirewallRuleset() tfsdk.Schema {
	port := tfsdk.Attribute{
		Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
			"from": {
				Type:     types.NumberType,
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					validators.Range(float64(1), float64(65535.0)),
					validators.Compare(validators.ComparatorLessThanEqual, "to"),
				},
			},
			"to": {
				Type:     types.NumberType,
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					validators.Range(float64(1), float64(65535.0)),
				},
			},
		}),
		Optional:    true,
		Description: "A port range. Conflicts with `port_group`.",
		Validators: []tfsdk.AttributeValidator{
			validators.ConflictsWith("port_group"),
		},
	}

	portGroup := tfsdk.Attribute{
		Type:        types.StringType,
		Optional:    true,
		Description: "The port group this rule applies to. If not provided, all ports will be matched. Conflicts with `port`.",
		Validators: []tfsdk.AttributeValidator{
			validators.ConflictsWith("port"),
		},
	}

	address := tfsdk.Attribute{
		Type:        types.StringType,
		Optional:    true,
		Description: "The cidr this rule applies to. If not provided, it is treated as `0.0.0.0/0`. Conflicts with `address_group`.",
		Validators: []tfsdk.AttributeValidator{
			validators.Cidr(),
			validators.ConflictsWith("address_group"),
		},
	}

	addressGroup := tfsdk.Attribute{
		Type:        types.StringType,
		Optional:    true,
		Description: "The address group this rule applies to. If not provided, all addresses will be matched. Conflicts with `address`.",
		Validators: []tfsdk.AttributeValidator{
			validators.ConflictsWith("address"),
		},
	}

	return tfsdk.Schema{
		Description: "A grouping of firewall rules. The firewall is not enforced unless attached to an interface which can be done with the `firewall_ruleset_attachment` resource.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The identifier of the resource. This will always be the name. It is present only for legacy purposes.",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Description: "A unique, human readable name for this ruleset.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.NoWhitespace(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
			},
			"description": {
				Description: "A human readable description for this ruleset.",
				Type:        types.StringType,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"default_action": {
				Description: "The default action to take if traffic is not matched by one of the rules in the ruleset. Must be one of `reject`, `drop`, `accept`.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSlice(true, "reject", "drop", "accept"),
				},
			},
			"default_logging": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Turn on logging for this rule. These rotated logs can be found in /var/log/messages on your router.",
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
						Type:        types.NumberType,
						Required:    true,
						Description: "The priority of this rule. The higher the priority, the higher the precedence.",
					},
					"description": {
						Type:        types.StringType,
						Optional:    true,
						Description: "A human readable description for this rule.",
						Validators: []tfsdk.AttributeValidator{
							validators.MinLength(1),
						},
					},
					"log": {
						Type:        types.BoolType,
						Optional:    true,
						Description: "Turn on logging for this rule. These rotated logs can be found in /var/log/messages on your router.",
					},
					"action": {
						Type:        types.StringType,
						Required:    true,
						Description: "The action to take on traffic that matches this rule. Must be one of `reject`, `drop`, `accept`.",
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, "drop", "reject", "accept"),
						},
					},
					"protocol": {
						Type:        types.StringType,
						Optional:    true,
						Description: "The protocol this rule applies to. If not specified, this rule applies to all protcols. Values prefixed with `!` specifies a _not_ behavior. If `!` is provided, this rule applies to all protocols except this one.",
						Validators: []tfsdk.AttributeValidator{
							validators.StringInSlice(true, append(
								append(
									protocols,
									utils.WithPrefix("!", protocols)...,
								),
								"all", "*",
							)...),
						},
					},
					"state": {
						Description: "This describes the connection state of a packet.",
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"established": {
								Type:        types.BoolType,
								Optional:    true,
								Description: "Match packets that are part of a two-way connection.",
							},
							"new": {
								Type:        types.BoolType,
								Optional:    true,
								Description: "Match packets creating a new connection.",
							},
							"related": {
								Type:        types.BoolType,
								Optional:    true,
								Description: "Match packets related to established connections.",
							},
							"invalid": {
								Type:        types.BoolType,
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
								Type:     types.StringType,
								Optional: true,
							},
						}),
						Optional: true,
						// Need a validator to ensure address conflicts with address_group and port conflicts with port_group.
					},
				},
			},
		},
	}
}
