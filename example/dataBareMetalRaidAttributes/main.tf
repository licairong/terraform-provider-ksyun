# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  bare metal raid attributes
data "ksyun_bare_metal_raid_attributes" "default" {
  output_file="output_result"
}

