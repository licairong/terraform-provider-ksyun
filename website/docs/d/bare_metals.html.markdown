---
layout: "ksyun"
page_title: "Ksyun: ksyun_bare_metals"
sidebar_current: "docs-ksyun-datasource-bare-metals"
description: |-
Provides a list of Bare Metal resources in the current region.
---

# ksyun_bare_metals

This data source provides a list of Bare Metal resources according to their Bare Metal ID.

## Example Usage

```hcl
# Get  bare metals
data "ksyun_bare_metals" "default" {
  output_file="output_result"
  ids = []
  vpc_id = ["bfec0f43-9e5a-4f06-b7a1-df4768c1cd6f"]
  project_id = []
  host_name = []
  subnet_id = []
  cabinet_id = []
  epc_host_status = []
  os_name = []
  product_type = []
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of Bare Metal IDs, all the Bare Metals belong to this region will be retrieved if the ID is `""`.
* `vpc_id` - (Optional) One or more vpc IDs.
* `project_id` - (Optional) One or more project IDs.
* `host_name` - (Optional) One or more Bare Metal host names.
* `subnet_id` - (Optional) One or more subnet IDs.
* `cabinet_id` - (Optional) One or more Bare Metal cabinet IDs.
* `epc_host_status` - (Optional) One or more Bare Metal status.
* `os_name` - (Optional) One or more Bare Metal operating system names.
* `product_type` - (Optional) One or more Bare Metal product types,Valid is lease or customer or lending.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bare_metals` - It is a nested type which documented below.
* `total_count` - Total number of Bare Metals that satisfy the condition.

