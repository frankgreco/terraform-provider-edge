resource "edge_firewall_address_group" "router" {
    name        = "router"
    description = "router interface addresses"

    cidrs = [
        "192.168.2.1",
        "192.168.3.1",
        "192.168.4.1",
    ]
}


resource "edge_firewall_ruleset" "example" {
  name            = "example"
  description     = "drop all ssh traffic to the router"
  default_action  = "accept"

  rule {
    priority    = 10
    description = "ssh"
    action      = "drop" 
    protocol    = "tcp"

    destination = {
      from_port = 22
      to_port   = 22

      address_group = edge_firewall_address_group.router.name
    }
  }

}
