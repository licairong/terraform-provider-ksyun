package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBareMetalRaidAttributes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBareMetalRaidAttributesRead,
		Schema: map[string]*schema.Schema{
			"host_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"raid_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raid_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"raid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_space": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"space": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"system_disk_space": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunBareMetalRaidAttributesRead(d *schema.ResourceData, meta interface{}) error {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	return bareMetalService.ReadAndSetRaidAttributes(d, dataSourceKsyunBareMetalRaidAttributes())
}
