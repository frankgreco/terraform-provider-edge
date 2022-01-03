package provider

import (
	"context"
	"os"

	"github.com/frankgreco/edge-sdk-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var stderr = os.Stderr

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     *edge.Client
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "The Edge provider provides the ability to configure a Ubiquiti Edge device.",
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Edge router URL. Can be set with `EDGE_HOST`.",
			},
			"username": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Admin username. Can be set with `EDGE_USERNAME`.",
			},
			"password": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: "Admin password. Can be set with `EDGE_PASSWORD`.",
			},
		},
	}, nil
}

type providerData struct {
	Username types.String `tfsdk:"username"`
	Host     types.String `tfsdk:"host"`
	Password types.String `tfsdk:"password"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	{
		diags := req.Config.Get(ctx, &config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var username string
	{
		if config.Username.Unknown {
			resp.Diagnostics.AddWarning(
				"Unable to create client",
				"Cannot use unknown value as username",
			)
			return
		}
		if config.Username.Null {
			username = os.Getenv("EDGE_USERNAME")
		} else {
			username = config.Username.Value
		}

		if username == "" {
			resp.Diagnostics.AddError(
				"Unable to find username",
				"Username cannot be an empty string",
			)
			return
		}
	}

	// User must provide a password to the provider
	var password string
	if config.Password.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if config.Password.Null {
		password = os.Getenv("EDGE_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if password == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find password",
			"password cannot be an empty string",
		)
		return
	}

	// User must specify a host
	var host string
	if config.Host.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
		return
	}

	if config.Host.Null {
		host = os.Getenv("EDGE_HOST")
	} else {
		host = config.Host.Value
	}

	if host == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Host cannot be an empty string",
		)
		return
	}

	c, err := edge.Login(host, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Unable to create edge client:\n\n"+err.Error(),
		)
		return
	}

	p.client = c
	p.configured = true
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"edge_firewall_ruleset":            resourceFirewallRulesetType{},
		"edge_firewall_ruleset_attachment": resourceFirewallRulesetAttachmentType{},
		"edge_firewall_address_group":      resourceFirewallAddressGroupType{},
		"edge_firewall_port_group":         resourceFirewallPortGroupType{},
	}, nil
}

func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"edge_interface_ethernet": dataSourceInterfaceEthernetType{},
	}, nil
}
