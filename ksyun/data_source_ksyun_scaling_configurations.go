package ksyun

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"regexp"
)

func dataSourceKsyunScalingConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunScalingConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
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
			"project_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"scaling_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"scaling_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"mem": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"data_disk_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"gpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"need_monitor_agent": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"need_security_agent": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_name_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"keep_image_login": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"system_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"key_id": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"data_disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"delete_with_instance": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"instance_name_time_suffix": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"address_band_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"band_width_share_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"line_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"address_project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	var result []map[string]interface{}
	var allScalingConfigurations []interface{}
	var err error
	r := dataSourceKsyunScalingConfigurations()

	limit := 10
	offset := 1

	client := meta.(*KsyunClient)
	allScalingConfigurations = []interface{}{}
	result = []map[string]interface{}{}
	readScalingConfiguration := make(map[string]interface{})

	var only map[string]SdkReqTransform
	only = map[string]SdkReqTransform{
		"ids":                        {mapping: "ScalingConfigurationId", Type: TransformWithN},
		"project_ids":                {mapping: "ProjectId", Type: TransformWithN},
		"scaling_configuration_name": {},
	}

	readScalingConfiguration, err = SdkRequestAutoMapping(d, r, false, only, nil)

	if err != nil {
		return fmt.Errorf("error on reading ScalingConfiguration list, %s", err)
	}

	for {
		readScalingConfiguration["MaxResults"] = limit
		readScalingConfiguration["Marker"] = offset

		logger.Debug(logger.ReqFormat, "DescribeScalingConfiguration", readScalingConfiguration)
		resp, err := client.kecconn.DescribeScalingConfiguration(&readScalingConfiguration)
		if err != nil {
			return fmt.Errorf("error on reading ScalingConfiguration list req(%v):%v", readScalingConfiguration, err)
		}
		l := (*resp)["ScalingConfigurationSet"].([]interface{})
		allScalingConfigurations = append(allScalingConfigurations, l...)
		if len(l) < limit {
			break
		}

		offset = offset + limit
	}

	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, v := range allScalingConfigurations {
			item := v.(map[string]interface{})
			if r != nil && !r.MatchString(item["ScalingConfigurationName"].(string)) {
				continue
			}
			result = append(result, item)
		}
	} else {
		merageResultDirect(&result, allScalingConfigurations)
	}

	err = dataSourceKsyunScalingConfigurationsSave(d, result)
	if err != nil {
		return fmt.Errorf("error on reading ScalingConfigurationName list, %s", err)
	}

	return nil
}

func scalingConfigurationSpecialMapping() map[string]SdkResponseMapping {
	specialMapping := make(map[string]SdkResponseMapping)
	specialMapping["StorageSize"] = SdkResponseMapping{Field: "data_disk_gb"}
	specialMapping["DataDiskEbsDetail"] = SdkResponseMapping{
		Field: "data_disks",
		FieldRespFunc: func(i interface{}) interface{} {
			var result []map[string]interface{}
			result = []map[string]interface{}{}
			v := i.(string)
			var dat []interface{}
			if err := json.Unmarshal([]byte(v), &dat); err == nil {
				for _, v := range dat {
					d := v.(map[string]interface{})
					r := make(map[string]interface{})
					r["delete_with_instance"] = d["deleteWithInstance"]
					r["disk_size"] = d["size"]
					r["disk_type"] = d["type"]
					result = append(result, r)
				}
			}
			return result
		},
	}
	specialMapping["InstanceTypeSet"] = SdkResponseMapping{
		Field: "instance_type",
		FieldRespFunc: func(i interface{}) interface{} {
			var result string
			for _, v := range i.([]interface{}) {
				result = result + v.(string) + ","
			}
			if len(result) > 0 {
				result = result[0 : len(result)-1]
			}
			return result
		},
	}
	specialMapping["KeyIds"] = SdkResponseMapping{
		Field: "key_id",
		FieldRespFunc: func(i interface{}) interface{} {
			var (
				value []string
			)
			_ = json.Unmarshal([]byte(i.(string)), &value)
			return value

		},
	}
	return specialMapping
}

func dataSourceKsyunScalingConfigurationsSave(d *schema.ResourceData, result []map[string]interface{}) error {
	resource := dataSourceKsyunScalingConfigurations()
	targetName := "scaling_configurations"
	_, _, err := SdkSliceMapping(d, result, SdkSliceData{
		IdField: "ScalingConfigurationId",
		IdMappingFunc: func(idField string, item map[string]interface{}) string {
			return item["ScalingConfigurationId"].(string)
		},
		SliceMappingFunc: func(item map[string]interface{}) map[string]interface{} {
			delete(item, "InstanceType")
			return SdkResponseAutoMapping(resource, targetName, item, nil, scalingConfigurationSpecialMapping())
		},
		TargetName: targetName,
	})
	return err
}
