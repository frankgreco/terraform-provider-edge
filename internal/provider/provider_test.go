package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	providerFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"edge": func() (tfprotov6.ProviderServer, error) {
			return tfsdk.NewProtocol6Server(New()), nil
		},
	}
)
