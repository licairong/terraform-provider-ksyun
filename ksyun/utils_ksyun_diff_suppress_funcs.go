package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func purchaseTimeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("charge_type"); ok && (v.(string) == "Monthly" || v.(string) == "PrePaidByMonth") {
		return false
	}
	return true
}

func kecImportDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	//由于一些字段暂时无法支持从查询中返回 所以现在设立做特殊处理拦截变更 用来适配导入的场景 后续支持后在对导入场景做优化
	if !d.IsNewResource() {
		if !d.Get("has_init_info").(bool) {
			if k == "local_volume_snapshot_id" {
				return true
			}
			if k == "user_data" {
				return true
			}
		}
	}
	if (k == "keep_image_login" || k == "key_id") && !d.IsNewResource() && !d.HasChange("image_id") {
		return true
	}

	return false
}

func kcsParameterDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if old != "" && new == "" {
		return true
	}
	if k == "parameters.notify-keyspace-events" && old == "" && new == "" {
		return true
	}
	return false
}

func rdsParameterDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if k == "parameters.#" {
		logger.Debug(logger.RespFormat, "DemoTest", d.ConnInfo())
		logger.Debug(logger.RespFormat, "DemoTest", d.Get("parameters"))
		return false
	}
	return true
}

func kcsSecurityGroupDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("security_group_ids") != nil && old == "" && new != "" {
		if sgs, ok := d.Get("security_group_ids").(*schema.Set); ok {
			if (*sgs).Contains(new) {
				err := d.Set("security_group_id", new)
				if err == nil {
					return true
				}
			}
		}
	}
	return false
}

func networkAclEntryDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("protocol") != "icmp" && (k == "icmp_type" || k == "icmp_code") {
		return true
	}
	if d.Get("protocol") != "tcp" && d.Get("protocol") != "udp" && (k == "port_range_from" || k == "port_range_to") {
		return true
	}
	return false
}

func securityGroupEntryDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("protocol") != "icmp" && (k == "icmp_type" || k == "icmp_code") {
		return true
	}
	if d.Get("protocol") != "tcp" && d.Get("protocol") != "udp" && (k == "port_range_from" || k == "port_range_to") {
		return true
	}
	return false
}
