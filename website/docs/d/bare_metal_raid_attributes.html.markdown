---
layout: "ksyun"
page_title: "Ksyun: ksyun_bare_metal_raid_attributes"
sidebar_current: "docs-ksyun-datasource-bare-metal-raid-attributes"
description: |-
Provides a list of Bare Metal Raid Attributes resources in the current region.
---

# ksyun_bare_metal_raid_attributes

This data source provides a list of Bare Metal Raid Attributes resources according to their Bare Metal Raid Attribute ID.

## Example Usage

```hcl
# Get  bare metal_raid_attributes
data "ksyun_bare_metal_raid_attributes" "default" {
  output_file="output_result"
}

```

## Argument Reference

The following arguments are supported:

* `host_type` - (Optional) A list of Bare Metal Raid Attribute Host Types.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `raid_attributes` - It is a nested type which documented below.
* `total_count` - Total number of Bare Metal Raid Attributes that satisfy the condition.

