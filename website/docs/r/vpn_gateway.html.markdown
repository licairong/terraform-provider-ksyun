---
layout: "ksyun"
page_title: "Ksyun: ksyun_vpn_gateway"
sidebar_current: "docs-ksyun-resource-vpn-gateway"
description: |-
  Provides a Vpn Gateway resource under VPC resource.
---

# ksyun_vpn_gateway

Provides a Vpn Gateway resource under VPC resource.

## Example Usage

```hcl
resource "ksyun_vpn_gateway" "default" {
  vpn_gateway_name   = "ksyun_vpn_gw_tf1"
  band_width = 10
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  charge_type = "Daily"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_name` - (Optional) The name of the vpn gateway.
* `band_width` - (Required) The bandWidth of the vpn gateway.Valid Values:5,10,20,50,100,200.
* `vpc_id` - (Required, ForceNew) The id of the vpc.
* `charge_type` -  (Required, ForceNew)  The charge type of the vpn gateway.Valid Values:'Monthly','Daily'
* `purchase_time` - (Optional, ForceNew) The purchase time of the vpn gateway.
* `project_id` - (Optional) The project id  of the vpn gateway.Default is 0 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpn_gateway_id` - The id of creation of vpn gateway

## Import

Vpn Gateway can be imported using the `id`, e.g.

```
$ terraform import ksyun_vpn_gateway.default $id
```