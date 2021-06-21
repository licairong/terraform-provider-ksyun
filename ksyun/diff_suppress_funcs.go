package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func purchaseTimeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("charge_type"); ok && (v.(string) == "Monthly" || v.(string) == "PrePaidByMonth") {
		return false
	}
	return true
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
