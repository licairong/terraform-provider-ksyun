resource "ksyun_network_acl_entry" "test" {
  description = "测试1"
  cidr_block = "10.0.16.0/24"
  rule_number = 16
  direction = "in"
  rule_action = "deny"
  protocol = "ip"
  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
}