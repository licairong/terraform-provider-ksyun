

variable "available_zone" {
  default = "cn-shanghai-2a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-terraform"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds_security_group" "krds_sec_group_237" {
  security_group_name = "terraform_security_group_237"
  security_group_description = "terraform-security-group-237"
  security_group_rule{
    security_group_rule_protocol = "182.133.0.0/16"
    security_group_rule_name = "asdf"
  }
  security_group_rule{
    security_group_rule_protocol = "182.134.0.0/16"
    security_group_rule_name = "asdf2"
  }
}


resource "ksyun_krds" "houbin_terraform_4"{
  db_instance_class= "db.ram.1|db.disk.20"
  db_instance_name = "houbin_terraform_1-1"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.6"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  security_group_id = "${ksyun_krds_security_group.krds_sec_group_237.id}"
  preferred_backup_time = "02:00-03:00"
  port=3306
  parameters {
    name = "auto_increment_increment"
    value = "10"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

}

//resource "ksyun_krds_rr" "rds-rr-1"{
//  source_db_instance_identifier= "${ksyun_krds.houbin_terraform_4.id}"
//  db_instance_class= "db.ram.2|db.disk.30"
//  db_instance_name = "houbin_terraform_777_rr_1"
//  bill_type = "DAY"
//}


