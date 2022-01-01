package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type dataSourceFirewallRulesetType struct{}

func (r dataSourceFirewallRulesetType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return resourceFirewallRulesetSchema(true), nil
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
