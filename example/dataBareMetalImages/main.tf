# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  bare metal images
data "ksyun_bare_metal_images" "default" {
  output_file="output_result"
}

