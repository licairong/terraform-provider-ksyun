package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

func kecNetworkInterfaceCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.Id() != "" && (d.HasChange("private_ip_address") || d.HasChange("subnet_id") || d.HasChange("security_group_ids")) {
		var data []interface{}
		vpcService := VpcService{meta.(*KsyunClient)}
		condition := map[string]interface{}{
			"NetworkInterfaceId.1": d.Id(),
		}
		data, err = vpcService.ReadNetworkInterfaces(condition)
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

func networkAclEntryCustomizeDiff(d *schema.ResourceDiff, meta interface{}) (err error) {
	if d.HasChange("network_acl_entries") {
		m := make(map[string]interface{})
		for _, v := range d.Get("network_acl_entries").(*schema.Set).List() {
			num := v.(map[string]interface{})["rule_number"].(int)
			direction := v.(map[string]interface{})["direction"].(string)
			if _, ok := m[strconv.Itoa(num)+"-"+direction]; ok {
				return fmt.Errorf("rule_number must unique")
			} else {
				m[strconv.Itoa(num)] = ""
			}
		}
	}
	return err
}
