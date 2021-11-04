---
layout: "ksyun"
page_title: "Ksyun: ksyun_ssh_keys"
sidebar_current: "docs-ksyun-datasource-ssh-keys"
description: |-
Provides a list of SSH Key resources in the current region.
---

# ksyun_ssh_keys

This data source provides a list of SSH Key resources according to their SSH Key ID.

## Example Usage

```hcl
# Get  ssh key
data "ksyun_ssh_keys" "default" {
  output_file="output_result"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of SSH Key IDs, all the SSH Key belong to this region will be retrieved if the ID is `""`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - It is a nested type which documented below.
* `total_count` - Total number of SSH Keys that satisfy the condition.

