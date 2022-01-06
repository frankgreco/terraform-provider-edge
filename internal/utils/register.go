package utils

import (
	"context"

	"github.com/frankgreco/edge-sdk-go/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type terraformRetriever interface {
	Get(context.Context, interface{}) diag.Diagnostics
}

// don't know how to reflect ... into interface {}
func retrieve(ctx context.Context, r terraformRetriever, target interface{}) (interface{}, diag.Diagnostics) {
	switch target.(type) {
	case types.AddressGroup:
		var tmp types.AddressGroup
		diags := r.Get(ctx, &tmp)
		return tmp, diags
	case types.PortGroup:
		var tmp types.PortGroup
		diags := r.Get(ctx, &tmp)
		return tmp, diags
	case types.FirewallAttachment:
		var tmp types.FirewallAttachment
		diags := r.Get(ctx, &tmp)
		return tmp, diags
	case types.Ruleset:
		var tmp types.Ruleset
		diags := r.Get(ctx, &tmp)
		(&tmp).SetCodecMode(types.CodecModeLocal)
		return tmp, diags
	default:
		var diags diag.Diagnostics
		diags.AddError("Could not unmarshal terraform plan", "Unknown go type")
		return nil, diags
	}
}
