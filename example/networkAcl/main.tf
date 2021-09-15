resource "ksyun_network_acl" "test" {
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  network_acl_name = "ceshi"
  network_acl_entries {
    description = "232323"
    cidr_block = "10.0.3.0/24"
    rule_number = 3
    direction = "in"
    rule_action = "allow"
    protocol = "ip"
  }
  network_acl_entries {
    cidr_block = "10.0.1.0/24"
    rule_number = 1
    direction = "out"
    rule_action = "allow"
    protocol = "tcp"
    port_range_from = 2
    port_range_to = 80
  }
  network_acl_entries {
    description = "111111"
    cidr_block = "10.0.2.0/24"
    rule_number = 2
    direction = "in"
    rule_action = "deny"
    protocol = "ip"
  }



}