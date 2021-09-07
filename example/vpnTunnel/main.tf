resource "ksyun_vpn_tunnel" "default" {
  vpn_tunnel_name   = "ksyun_vpn_tunnel_tf_1"
  type = "Ipsec"
  vpn_gateway_id = "9b3d361e-f65b-464b-947a-fafb5cfb10d2"
  customer_gateway_id = "7f5a5c91-4814-41bf-b9d6-d9d811f4df0f"
  ike_dh_group = 2
  pre_shared_key = "123456789abcd"
}