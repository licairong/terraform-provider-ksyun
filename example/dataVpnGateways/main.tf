# Specify the provider and access details
provider "ksyun" {
}

data "ksyun_vpn_gateways" "default" {
  output_file="output_result"
}

