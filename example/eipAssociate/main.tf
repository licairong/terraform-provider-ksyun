//# Specify the provider and access details
//provider "ksyun" {
//  region = "eu-east-1"
//}
//
resource "ksyun_eip_associate" "slb" {
  allocation_id="43b04c7a-174e-458e-aa2d-03b6d55d92ea"
  instance_type="Slb"
  instance_id="48957fac-ed38-4dba-9180-dfbffcd5a20d"
//  network_interface_id=""
}
//resource "ksyun_eip_associate" "server" {
//  allocation_id="419782b7-6766-4743-afb7-7c7081214092"
//  instance_type="Ipfwd"
//  instance_id="566567677-6766-4743-afb7-7c7081214092"
//  network_interface_id="87945980-59659-04548-759045803"
//}
