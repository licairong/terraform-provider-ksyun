package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func importKecNetworkInterfaceAttachment(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("network_interface_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("instance_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
