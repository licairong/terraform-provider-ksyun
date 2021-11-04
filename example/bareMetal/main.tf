resource "ksyun_bare_metal" "default" {
  host_name = "xym-test-勿删"
  host_type = "MI-I2"
  image_id = "eb8c0428-476e-49af-8ccb-9fad2455a54c"
  key_id = "9c45b560-e51d-4aee-9e99-0e292476692d"
  network_interface_mode = "single"
#  raid_id = "e754e372-fb27-4489-8043-d9a9c33958c0"
  raid = "Raid1"
  availability_zone = "cn-beijing-6b"
  security_agent = "classic"
  cloud_monitor_agent = "classic"
  subnet_id = "d2fdc1b5-0280-4ca7-920b-0bd0453c130c"
  security_group_ids = ["7e2f45b5-e79d-4612-a7fc-fe74a50b639a"]
  private_ip_address = "10.0.12.71"
  dns1 = "198.18.254.30"
  dns2 = "198.18.254.31"
  system_file_type = "EXT4"
  container_agent = "supported"
  //  private_ip_address = "10.0.80.14"
  //  dns1 = "198.18.224.10"
  //  dns2 = "198.18.224.11"
  force_re_install = false
}


