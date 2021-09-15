# Specify the provider and access details
provider "ksyun" {
}

# Get  routes
data "ksyun_network_acls" "default" {
  output_file="output_result"
//  vpc_ids = ["769c780b-acbd-41ca-9a06-4960e2423c7e"]
}
