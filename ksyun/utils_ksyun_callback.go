package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type ksyunApiCallFunc func(d *schema.ResourceData, meta interface{}) error

func ksyunApiCall(api []ksyunApiCallFunc, d *schema.ResourceData, meta interface{}) (err error) {
	if api != nil {
		for _, f := range api {
			if f != nil {
				err = f(d, meta)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}
