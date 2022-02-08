package provider

import (
	"github.com/frankgreco/terraform-helpers/validators"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func schemaFirewallPortGroup() tfsdk.Schema {
	return tfsdk.Schema{
		Description: "A logical grouping of ports.",
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
				Description:   "A unique, human readable name for this port group.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoWhitespace(),
				},
			},
			"description": {
				Type:        types.StringType,
				Optional:    true,
				Description: "A human readable description for this port group.",
				Validators: []tfsdk.AttributeValidator{
					validators.MinLength(1),
				},
			},
			"port_ranges": {
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"from": {
						Type:     types.NumberType,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							validators.Range(float64(1), float64(65535.0)),
							validators.Compare(validators.ComparatorLessThan, "to"),
						},
					},
					"to": {
						Type:     types.NumberType,
						Required: true,
						Validators: []tfsdk.AttributeValidator{
							validators.Range(float64(1), float64(65535.0)),
						},
					},
				}, tfsdk.ListNestedAttributesOptions{}),
				Optional:    true,
				Description: "A list of port ranges.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoOverlap(),
				},
			},
			"ports": {
				Type:        types.ListType{ElemType: types.NumberType},
				Optional:    true,
				Description: "A list of port numbers.",
				Validators: []tfsdk.AttributeValidator{
					validators.NoOverlap(),
				},
			},
		},
	}
}
