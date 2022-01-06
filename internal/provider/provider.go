package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		MarkdownDescription: `
The Edge provider provides the ability to configure a Ubiquiti Edge device.

~> You must use ` + "`-parallelism=1` on all `destroy` and `apply` " + `operations because the EdgeOS configuration API is not safe for concurrent use.
`,
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
			"insecure": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Specify if the connection to the Edge configuration API should be insecure. Can be set with `EDGE_INSECURE`.",
			},
		},
	}, nil
}

type providerData struct {
	Username types.String `tfsdk:"username"`
	Host     types.String `tfsdk:"host"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	{
		resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	username, err := requiredString(config.Username, "username", "EDGE_USERNAME")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			err.Error(),
		)
	}

	password, err := requiredString(config.Password, "password", "EDGE_PASSWORD")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			err.Error(),
		)
	}

	host, err := requiredString(config.Host, "host", "EDGE_HOST")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			err.Error(),
		)
	}

	var insecure bool
	{
		if !config.Insecure.Null && !config.Insecure.Unknown {
			insecure = config.Insecure.Value
		}
		if strings.ToUpper(os.Getenv("EDGE_INSECURE")) == "TRUE" {
			insecure = true
		} else if strings.ToUpper(os.Getenv("EDGE_INSECURE")) == "FALSE" {
			insecure = false
		}
	}

	c, err := edge.Login(host, insecure, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure provider",
			"Unable to create edge client: "+err.Error(),
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

func requiredString(str types.String, name, env string) (string, error) {
	if str.Unknown {
		return "", fmt.Errorf("Cannot use unknown value for %s.", name)
	}

	val := str.Value

	if str.Null {
		val = os.Getenv(env)
	}

	if val == "" {
		return "", fmt.Errorf("The provider attribute %s must be defined.", name)
	}

	return val, nil
}
