resource "edge_firewall_address_group" "example" {
    name        = "example"
    description = "example address group"

    cidrs = [
        "192.168.2.1",
        "192.168.3.0/24"
    ]
}
