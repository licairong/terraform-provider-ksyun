package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunBareMetals() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunBareMetalsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"project_id": {
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
			"host_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"vpc_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"subnet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"cabinet_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"host_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"epc_host_status": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"os_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"product_type": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"lease",
						"customer",
						"lending",
					}, false),
				},
				Set: schema.HashString,
			},
			"bare_metals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cabinet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"raid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"frequence": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"core_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"disk_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"raid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"space": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"network_interface_attribute_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_interface_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subnet_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns1": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dns2": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mac": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_set": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"security_group_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"network_interface_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_monitor_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_bond": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunBareMetalsRead(d *schema.ResourceData, meta interface{}) error {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	return bareMetalService.ReadAndSetBareMetals(d, dataSourceKsyunBareMetals())
}
