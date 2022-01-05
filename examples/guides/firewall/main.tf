data "edge_interface_ethernet" "eth1" {
  id = "eth1"
}

data "edge_interface_ethernet" "eth2" {
  id = "eth2"
}

data "edge_interface_ethernet" "eth3" {
  id = "eth3"
}

resource "edge_firewall_address_group" "router" {
    name        = "router"
    description = "router interface addresses"

    cidrs = [
        "192.168.2.1",
        "192.168.3.1",
        "192.168.4.1"
    ]
}

resource "edge_firewall_port_group" "whitelist" {
    name        = "whitelist"
    description = "common known ports"

    port_ranges = [{
        from = 8080
        to   = 8081
    }]

    ports = [80, 443, 22]
}

resource "edge_firewall_ruleset" "router" {
  name            = "router"
  description     = "router bound traffic"
  default_action  = "drop"

  rule {
    priority    = 10
    description = "common known ports"
    action      = "accept" 
    protocol    = "all"

    destination = {
      address_group = edge_firewall_address_group.router.name
      port_group    = edge_firewall_port_group.whitelist.name
    }
  }

  rule {
    priority    = 20
    description = "no ssh from a specific address in the previous range"
    action      = "drop" 
    protocol    = "tcp"

    destination = {
      address_group = edge_firewall_address_group.router.name
      port = {
          from = 23
          to   = 22
      }
    }

    source = {
        address = "192.168.2.44/32"
    }
  }
}

// resource "edge_firewall_ruleset_attachment" "eth1" {
//   interface = data.edge_interface_ethernet.eth1.id 
//   in        = edge_firewall_ruleset.router.name
// }

resource "edge_firewall_ruleset_attachment" "eth2" {
  interface = data.edge_interface_ethernet.eth2.id 
  in        = edge_firewall_ruleset.router.name
}

resource "edge_firewall_ruleset_attachment" "eth3" {
  interface = data.edge_interface_ethernet.eth3.id 
  in        = edge_firewall_ruleset.router.name
}
