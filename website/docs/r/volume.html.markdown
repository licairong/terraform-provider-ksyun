---
layout: "ksyun"
page_title: "Ksyun: ksyun_volume"
sidebar_current: "docs-ksyun-resource-volume"
description: |-
  Provides an EBS resource.
---


# ksyun_volume

Provides a EBS resource.

## Example Usage

```h
resource "ksyun_volume" "default" {
  volume_name       = "test"
  volume_type       = "SSD3.0"
  size              = 15
  charge_type       = "Daily"
  availability_zone = "cn-shanghai-3a"
  volume_desc       = "test"

  ## 传入快照ID，用快照创建EBS盘
  ## 注意：如果使用的整机镜像创建主机，API会自动根据镜像中包含的快照创建数据盘，不需在tf配置中定义数据盘
  # snapshot_id = "snapshot_id"
}
```

## Argument Reference

The following arguments are supported:

* `volume_name` - (Optional) The name of the EBS volume.
* `volume_type` - (Optional, ForceNew) The type of the EBS volume. Valid values:ESSD_PL1/ESSD_PL2/ESSD_PL3/SSD3.0/EHDD, default is `SSD3.0`
* `volume_desc` - (Optional) The description of the EBS volume.
* `size` - (Optional) The capacity of the EBS volume, in GB. Default is 10.
* `availability_zone` - (Required, ForceNew) The availability zone in which the EBS volume resides. For more information
* `charge_type` - (Required, ForceNew) The billing mode of the EBS volume.
* `project_id` - (Optional) The ID of the project to which the EBS volume belongs
* `online_resize` - (Optional) Specifies whether to expand the capacity of the EBS volume online, default is true.
* `snapshot_id` - (Optional, ForceNew) When the cloud disk snapshot opens, the snapshot id is entered


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `volume_status` - The status of the EBS volume
* `create_time` - The time when the EBS volume was created.
* `volume_category` - The category to which the EBS volume belongs. Valid values: system and data.



## Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_volume.default xxxxxx
```