terraform {
  required_providers {
    edge = {
      source  = "frankgreco/ubiquiti/edge"
      version = "0.0.1"
    }
  }

  required_version = "~> 1.1.2"
}

provider "edge" {
  username = "<ubnt>"
  password = "<password>"
  host     = "<host>"
}

data "edge_firewall_ruleset" "foo" {
  name = "WAN_IN"
}

output "greco" {
  value = data.edge_firewall_ruleset.foo
}

data "edge_interface_ethernet" "eth3" {
  id = "eth3"
}

resource "edge_firewall_ruleset_attachment" "foo" {
  interface = data.edge_interface_ethernet.eth3.id 
  in        = edge_firewall_ruleset.test.name
  local     = edge_firewall_ruleset.test.name
  out       = edge_firewall_ruleset.test.name
}

resource "edge_firewall_ruleset" "test" {
  name            = "test_firewall_2"
  description     = "frank"
  default_action  = "accept"

  rule {
    priority    = 10
    description = "rule description 5"
    action      = "drop" 
    protocol    = "tcp"

    destination = {
      from_port = 80
      to_port   = 80
      address   = "0.0.0.0/0"
    }
  }

  rule {
    priority    = 30
    description = "rule description 2"
    action      = "accept" 
    protocol    = "tcp"

    destination = {
      from_port = 80
      to_port   = 80
      address   = "0.0.0.0/0"
    }
  }

  rule {
    priority    = 20
    description = "rule description 3"
    action      = "accept"
    protocol    = "tcp"

    destination = {
      from_port = 80
      to_port   = 80
      address   = "0.0.0.0/0"
    }
  }

  rule {
    priority    = 100
    description = "rule description 10"
    action      = "accept"
    protocol    = "tcp"

    destination = {
      from_port = 80
      to_port   = 80
      address   = "0.0.0.0/0"
    }
  }

}