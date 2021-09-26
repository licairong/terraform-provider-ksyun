package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunLines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunLinesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"lines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunLinesRead(d *schema.ResourceData, meta interface{}) error {
	eipService := EipService{meta.(*KsyunClient)}
	return eipService.ReadAndSetLines(d, dataSourceKsyunLines())
}
