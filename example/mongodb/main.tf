resource "ksyun_mongodb_instance" "default" {
  name = "ymq-mongodb001"
  instance_account = "root"
  instance_password = "admin"
  instance_class = "1C2G"
  storage = 50
  node_num = 5
  vpc_id = "c7f060c0-6d0d-4adb-987f-fda1fc988ffe"
  vnet_id = "af6ed8fb-1246-47f5-bce0-e8e61eaeeb22"
  db_version = "3.6"
  pay_type = "byDay"
  iam_project_id = "101812"
  availability_zone = "cn-beijing-6a"
  cidrs = "10.34.51.13/32"
}