---
layout: "ksyun"
page_title: "Ksyun: ksyun_network_acl_entry"
sidebar_current: "docs-ksyun-resource-network-acl-entry"
description: |-
  Provides a Network ACL Entry resource under Network ACL resource.
---

# ksyun_network_acl_entry

Provides a Network ACL Entry resource under Network ACL resource.

## Example Usage

```hcl
resource "ksyun_network_acl_entry" "test" {
  description = "测试1"
  cidr_block = "10.0.16.0/24"
  rule_number = 16
  direction = "in"
  rule_action = "deny"
  protocol = "ip"
  network_acl_id = "679b6a88-67dd-4e17-a80a-985d9673050e"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the network acl entry.
* `network_acl_id` - (Required, ForceNew) The id of the network acl.
* `cidr_block` - (Required, ForceNew) The cidr_block of the network acl entry.
* `rule_number` - (Required, ForceNew) The rule_number of the network acl entry.
* `direction` - (Required, ForceNew) The direction of the network acl entry.Valid Value: 'in','out'.
* `rule_action` - (Required, ForceNew) The rule_action of the network acl entry.Valid Value: 'allow','deny'.
* `protocol` - (Required, ForceNew) The protocol of the network acl entry.Valid Value: 'ip','icmp','tcp','udp'.
* `icmp_type` - (Optional, ForceNew) The icmp_type of the network acl entry.If protocol is icmp,Required.
* `icmp_code` - (Optional, ForceNew) The icmp_code of the network acl entry.If protocol is icmp,Required.
* `port_range_from` - (Optional, ForceNew) The port_range_from of the network acl entry.If protocol is tcp or udp,Required.
* `port_range_to` - (Optional, ForceNew) The port_range_to of the network acl entry.If protocol is tcp or udp,Required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_acl_entry_id` - The id of creation of network acl entry.

## Import

Network ACL Entry can be imported using the `network_acl_id:rule_number:direction`, e.g.

```
$ terraform import ksyun_network_acl_entry.default $network_acl_id:$rule_number$direction
```