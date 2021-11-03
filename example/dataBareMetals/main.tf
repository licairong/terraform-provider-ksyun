# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  bare metals
data "ksyun_bare_metals" "default" {
  output_file="output_result"
  vpc_id = ["bfec0f43-9e5a-4f06-b7a1-df4768c1cd6f"]
}

