# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  keys
data "ksyun_ssh_keys" "default" {
  output_file="output_result"
  ids=[]
  key_name=""
}

