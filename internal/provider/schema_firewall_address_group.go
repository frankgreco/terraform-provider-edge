package provider

import (
	"github.com/frankgreco/terraform-helpers/validators"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func schemaFirewallAddressGroup() tfsdk.Schema {
	return tfsdk.Schema{
		Description: "A logical grouping of addresses.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The identifier of the resource. This will always be the name. It is present only for legacy purposes.",
				Type:        types.StringType,
				Computed:    true,
			},
			"name": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: []tfsdk.AttributePlanModifier{tfsdk.RequiresReplace()},
				Description:   "A unique, human readable name for this address group.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoWhitespace(),
				},
			},
			"description": {
				Type:        types.StringType,
				Optional:    true,
				Description: "A human readable description for this address group.",
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"cidrs": {
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Description: "A non-overlapping list of cidrs.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoOverlappingCIDRs(),
				},
			},
		},
	}
}
