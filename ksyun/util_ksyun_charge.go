package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func durationSchemaDiffSuppressFunc(field string,value interface{}) schema.SchemaDiffSuppressFunc{
	return func(k, old, new string, d *schema.ResourceData) bool {
		if v, ok := d.GetOk(field); ok && v == value {
			return false
		}
		return true
	}
}
