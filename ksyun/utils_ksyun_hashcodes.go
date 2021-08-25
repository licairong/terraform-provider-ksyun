package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"

func kecNetworkInterfaceHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	return hashcode.String(m["network_interface_id"].(string))
}
