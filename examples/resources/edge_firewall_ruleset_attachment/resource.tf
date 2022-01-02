data "edge_interface_ethernet" "eth3" {
  id = "eth3"
}

data "edge_firewall_ruleset" "ssh" {
  name = "ssh"
}

resource "edge_firewall_ruleset_attachment" "foo" {
  interface = data.edge_interface_ethernet.eth3.id 
  in        = data.edge_firewall_ruleset.ssh.name
}
