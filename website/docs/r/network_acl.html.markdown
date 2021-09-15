---
layout: "ksyun"
page_title: "Ksyun: ksyun_network_acl"
sidebar_current: "docs-ksyun-resource-network-acl"
description: |-
  Provides a Network ACL resource under VPC resource.
---

# ksyun_network_acl

Provides a Network ACL resource under VPC resource.

## Example Usage

```hcl
resource "ksyun_network_acl" "default" {
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  network_acl_name = "ceshi"
  network_acl_entries {
    description = "232323"
    cidr_block = "10.0.3.0/24"
    rule_number = 3
    direction = "in"
    rule_action = "allow"
    protocol = "ip"
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_name` - (Optional) The name of the network acl.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `network_acl_entries` -  (Optional)  The entry set of the vpn network acl.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - The id of creation of network acl

## Import

Network ACL can be imported using the `id`, e.g.

```
$ terraform import ksyun_network_acl.default $id
```