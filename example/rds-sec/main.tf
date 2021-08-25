resource "ksyun_krds_security_group" "default" {
  security_group_name = "SecGroup_d24b58bb-a311-47ef-b5f4-77d93dbb1938"
  security_group_description = "Security Group for d24b58bb-a311-47ef-b5f4-77d93dbb1938"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "wtf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "wtf"
  }

}