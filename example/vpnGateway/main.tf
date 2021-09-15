resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name   = "ksyun_vpn_gw_tf1"
  band_width = 10
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type = "Daily"
}