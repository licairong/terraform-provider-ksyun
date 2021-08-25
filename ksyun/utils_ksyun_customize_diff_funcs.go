package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func kecNetworkInterfaceCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.Id() != "" && (d.HasChange("private_ip_address") || d.HasChange("subnet_id") || d.HasChange("security_group_ids")) {
		var data []interface{}
		vpcService := VpcService{meta.(*KsyunClient)}
		condition := map[string]interface{}{
			"NetworkInterfaceId.1": d.Id(),
		}
		data, err = vpcService.readNetworkInterfaces(condition)
		if err != nil {
			return err
		}
		if len(data) != 1 {
			return fmt.Errorf("NetworkInterface %s not exist ", d.Id())
		}
		if _, ok := data[0].(map[string]interface{})["InstanceId"]; ok {
			return err
		}
		if d.HasChange("private_ip_address") {
			err = d.ForceNew("private_ip_address")
			if err != nil {
				return err
			}
		}
		if d.HasChange("subnet_id") {
			err = d.ForceNew("subnet_id")
			if err != nil {
				return err
			}
		}
		if d.HasChange("security_group_ids") {
			err = d.ForceNew("security_group_ids")
			if err != nil {
				return err
			}
		}
	}
	return err
}
