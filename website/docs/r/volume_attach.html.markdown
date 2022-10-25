---
layout: "ksyun"
page_title: "Ksyun: ksyun_volume_attach"
sidebar_current: "docs-ksyun-resource-volume-attach"
description: |-
Provides an EBS attachment resource.
---

# ksyun_volume

Provides an EBS attachment resource.

## Example Usage

```h
resource "ksyun_volume_attach" "default" {
  volume_id   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  instance_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  delete_with_instance = true
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required, ForceNew) The ID of the EBS volume to be attached to a KEC instance.
* `instance_id` - (Required, ForceNew) The ID of the KEC instance to which the EBS volume is to be attached.
* `delete_with_instance` - (Optional) Specifies whether to delete the EBS volume when the KEC instance to which it is attached is deleted. Default value: false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `volume_name` - The name of the EBS volume.
* `volume_type` - The type of the EBS volume. Valid values:ESSD_PL1/ESSD_PL2/ESSD_PL3/SSD3.0/EHDD, default is `SSD3.0`
* `volume_desc` - The description of the EBS volume.
* `size` - The capacity of the EBS volume, in GB. Default is 10.
* `availability_zone` - The availability zone in which the EBS volume resides. For more information
* `project_id` - The ID of the project to which the EBS volume belongs
* `volume_status` - The status of the EBS volume
* `create_time` - The time when the EBS volume was created.
* `volume_category` - The category to which the EBS volume belongs. Valid values: system and data.

## Import

Instance can be imported using the `id`, e.g.

```
$ terraform import ksyun_volume.default <volume_id>:<instance_id>
```