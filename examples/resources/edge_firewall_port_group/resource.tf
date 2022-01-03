resource "edge_firewall_port_group" "example" {
  name        = "example"
  description = "example port group"

  port_ranges = [{
      from = 80
      to   = 90
    }, {
      from = 30000
      to   = 32000
  }]

  ports = [443, 22, 636]
}