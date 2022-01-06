package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func schemaFirewallRulesetAttachment() tfsdk.Schema {
	return tfsdk.Schema{
		Description: "Attach a firewall ruleset to inbound, outbound, and local traffic.",
		Attributes: map[string]tfsdk.Attribute{
			"interface": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Description:   "The interface to attach firewall rules to.",
			},
			"in": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Match inbound packets.",
			},
			"out": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Match outbound packets.",
			},
			"local": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Match local packets.",
			},
		},
	}
}
