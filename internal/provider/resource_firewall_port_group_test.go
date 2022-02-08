package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/frankgreco/edge-sdk-go/types"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEdgeFirewallPortGroup(t *testing.T) {
	group := &types.PortGroup{
		Name:        "acc_test",
		Description: strptr("description"),
		Ports:       []int{443, 22},
		Ranges: []*types.PortRange{
			{
				From: 80,
				To:   90,
			},
			{
				From: 30000,
				To:   32000,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactory,
		Steps: []resource.TestStep{
			{
				Config:      toPortGroupResource(group.Clone().WithPorts([]int{80, 80})),
				ExpectError: regexp.MustCompile("The element 80 is supplied by more than one range."),
			},
			{
				Config:      toPortGroupResource(group.Clone().WithRanges(80, 90, 83, 84)),
				ExpectError: regexp.MustCompile("The elements between 83 and 90 are supplied by more than one range."),
			},
			{
				Config: toPortGroupResource(group),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "id", "acc_test"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "name", "acc_test"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "description", "description"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "ports.#", "2"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "ports.0", "443"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "ports.1", "22"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "port_ranges.#", "2"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "port_ranges.0.from", "80"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "port_ranges.0.to", "90"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "port_ranges.1.from", "30000"),
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "port_ranges.1.to", "32000"),
				),
			},
			{
				Config: toPortGroupResource(group.Clone().WithDescription(nil)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("edge_firewall_port_group.acc_test", "description"),
				),
			},
			{
				Config: toPortGroupResource(group.Clone().WithDescription(strptr(""))),
				ExpectError: regexp.MustCompile("String must be at least 1 characters long."),
			},
			{
				Config: toPortGroupResource(group.Clone().WithDescription(strptr("updated"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("edge_firewall_port_group.acc_test", "description", "updated"),
				),
			},
		},
	})
}

// toPortGroupResourcce converts the go representation into the terraform representation.
func toPortGroupResource(group *types.PortGroup) string {
	var optionalDescription string
	{
		if group.Description != nil {
			optionalDescription = fmt.Sprintf("description = \"%s\"", *group.Description)
		}
	}

	var optionalPorts string
	{
		if group.Ports != nil {
			optionalPorts = fmt.Sprintf("ports = [%s]", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(group.Ports)), ", "), "[]"))
		}
	}

	var optionalPortRanges string
	{
		if group.Ranges != nil {
			ranges := []string{}
			for _, r := range group.Ranges {
				ranges = append(ranges, r.String())
			}
			optionalPortRanges = fmt.Sprintf("port_ranges = [%s]", strings.Join(ranges, ", "))
		}
	}

	return fmt.Sprintf(`
resource "edge_firewall_port_group" "acc_test" {
	name = "%s"
	%s
	%s
	%s
}`, group.Name, optionalDescription, optionalPortRanges, optionalPorts)
}
