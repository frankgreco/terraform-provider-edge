resource "edge_firewall_ruleset" "example" {
  name            = "example"
  description     = "drop all traffic on its way to 192.168.2.1/24 over port 80"
  default_action  = "accept"
  default_logging = true

  rule {
    priority    = 10
    description = "ssh"
    action      = "drop" 
    protocol    = "tcp"
    log         = true

    destination = {
      address = "192.168.2.1/24"
      port    = {
          from = 80
          to   = 80
      }
    }
  }
}
