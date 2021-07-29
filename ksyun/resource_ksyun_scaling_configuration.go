package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingConfigurationCreate,
		Read:   resourceKsyunScalingConfigurationRead,
		Delete: resourceKsyunScalingConfigurationDelete,
		Update: resourceKsyunScalingConfigurationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{

			"scaling_configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-scaling-config",
				ForceNew: false,
			},

			"image_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Required: true,
			},

			"instance_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "I1.1A",
				ForceNew:         false,
				ValidateFunc:     stringSplitSchemaValidateFunc(","),
				DiffSuppressFunc: stringSplitDiffSuppressFunc(","),
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"system_disk_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskType,
			},

			"system_disk_size": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateKecSystemDiskSize,
			},

			"data_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},

			"data_disks": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskType,
						},
						"disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateKecDataDiskSize,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"key_id": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_name_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_name_time_suffix": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"need_monitor_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
			},

			"need_security_agent": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateKecInstanceAgent,
			},

			"address_band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"band_width_share_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"line_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"address_project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"gpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"mem": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunScalingConfigurationExtra(d *schema.ResourceData, forceGet bool) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	var r map[string]SdkReqTransform

	r = map[string]SdkReqTransform{
		"instance_type":    {Type: TransformWithN},
		"key_id":           {Type: TransformWithN},
		"system_disk_type": {mapping: "SystemDisk.DiskType"},
		"system_disk_size": {mapping: "SystemDisk.DiskSize"},
		"data_disks": {mappings: map[string]string{
			"data_disks": "DataDisk",
			"disk_size":  "Size",
			"disk_type":  "Type",
		}, Type: TransformListN},
	}
	extra = SdkRequestAutoExtra(r, d, forceGet)
	return extra
}

func resourceKsyunScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var resp *map[string]interface{}
	var err error

	createScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, false, nil,
		resourceKsyunScalingConfigurationExtra(d, false))
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}

	action := "CreateScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, createScalingConfiguration)
	resp, err = conn.CreateScalingConfiguration(&createScalingConfiguration)
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfiguration, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ScalingConfigurationId"].(string))
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	scalingConfiguration := resourceKsyunScalingConfiguration()

	var err error

	modifyScalingConfiguration, err := SdkRequestAutoMapping(d, scalingConfiguration, true, nil, resourceKsyunScalingConfigurationExtra(d, true))
	if err != nil {
		return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
	}

	if len(modifyScalingConfiguration) > 0 {
		modifyScalingConfiguration["ScalingConfigurationId"] = d.Id()
		action := "ModifyScalingConfiguration"
		logger.Debug(logger.ReqFormat, action, modifyScalingConfiguration)
		_, err = conn.ModifyScalingConfiguration(&modifyScalingConfiguration)
		if err != nil {
			return fmt.Errorf("error on modifying ScalingConfiguration, %s", err)
		}
	}
	return resourceKsyunScalingConfigurationRead(d, meta)
}

func resourceKsyunScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	readScalingConfiguration := make(map[string]interface{})
	readScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	projectErr := AddProjectInfo(d, &readScalingConfiguration, client)
	if projectErr != nil {
		return projectErr
	}
	action := "DescribeScalingConfiguration"
	logger.Debug(logger.ReqFormat, action, readScalingConfiguration)
	resp, err := conn.DescribeScalingConfiguration(&readScalingConfiguration)
	if err != nil {
		return fmt.Errorf("error on reading ScalingConfiguration %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingConfigurationSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		delete(items[0].(map[string]interface{}), "InstanceType")
		SdkResponseAutoResourceData(d, resourceKsyunScalingConfiguration(), items[0], scalingConfigurationSpecialMapping())
	}
	return nil
}

func resourceKsyunScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	deleteScalingConfiguration := make(map[string]interface{})
	deleteScalingConfiguration["ScalingConfigurationId.1"] = d.Id()
	action := "DeleteScalingConfiguration"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteScalingConfiguration)
		resp, err1 := conn.DeleteScalingConfiguration(&deleteScalingConfiguration)
		logger.Debug(logger.AllFormat, action, deleteScalingConfiguration, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingConfiguration %q, %s", d.Id(), err1))
		}
	})

}
