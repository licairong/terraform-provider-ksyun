# Specify the provider and access details
provider "ksyun" {
}

data "ksyun_vpn_tunnels" "default" {
  output_file="output_result"
}

