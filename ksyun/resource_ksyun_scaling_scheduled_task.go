package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

var ksyunScalingScheduledTaskRepeatUnit = []string{
	"Day",
	"Month",
	"Week",
}

func resourceKsyunScalingScheduledTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunScalingScheduledTaskCreate,
		Read:   resourceKsyunScalingScheduledTaskRead,
		Delete: resourceKsyunScalingScheduledTaskDelete,
		Update: resourceKsyunScalingScheduledTaskUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				items := strings.Split(d.Id(), ":")
				if len(items) != 2 {
					return nil, fmt.Errorf("id must split with %s and size %v", ":", 2)
				}
				err = d.Set("scaling_group_id", items[0])
				if err != nil {
					return nil, err
				}
				err = d.Set("scaling_scheduled_task_id", items[1])
				if err != nil {
					return nil, err
				}
				d.SetId(items[1] + ":" + items[0])
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{

			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scaling_scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf-scaling-scheduled_task",
			},

			"readjust_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"readjust_min_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"readjust_expect_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},

			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"recurrence": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"repeat_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(ksyunScalingScheduledTaskRepeatUnit, false),
			},

			"repeat_cycle": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scaling_scheduled_task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunScalingScheduledTaskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingScheduledTask()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating ScalingScheduledTask, %s", err)
	}
	//zero process
	if _, ok := req["ReadjustMaxSize"]; !ok {
		req["ReadjustMaxSize"] = 0
	}
	if _, ok := req["ReadjustMinSize"]; !ok {
		req["ReadjustMinSize"] = 0
	}
	if _, ok := req["ReadjustExpectSize"]; !ok {
		req["ReadjustExpectSize"] = 0
	}

	action := "CreateScalingScheduledTask"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.CreateScheduledTask(&req)
	if err != nil {
		return fmt.Errorf("error on creating ScalingScheduledTask, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["ReturnSet"].(map[string]interface{})["ScalingScheduleTaskId"].(string) + ":" + req["ScalingGroupId"].(string))
	}
	return resourceKsyunScalingScheduledTaskRead(d, meta)
}

func resourceKsyunScalingScheduledTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	r := resourceKsyunScalingScheduledTask()

	var err error

	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return fmt.Errorf("error on modifying ScalingScheduledTask, %s", err)
	}
	if len(req) > 0 {
		req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
		req["ScalingScheduledTaskId"] = strings.Split(d.Id(), ":")[0]
		action := "ModifyScheduledTask"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyScheduledTask(&req)
		if err != nil {
			return fmt.Errorf("error on modifying ScalingScheduledTask, %s", err)
		}
	}
	return resourceKsyunScalingScheduledTaskRead(d, meta)
}

func resourceKsyunScalingScheduledTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn

	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingScheduledTaskId.1"] = strings.Split(d.Id(), ":")[0]
	action := "DescribeScheduledTask"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeScheduledTask(&req)
	if err != nil {
		return fmt.Errorf("error on reading ScalingScheduledTask %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ScalingScheduleTaskSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunScalingScheduledTask(), items[0], nil)
	}
	return nil
}

func resourceKsyunScalingScheduledTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.kecconn
	req := make(map[string]interface{})
	req["ScalingGroupId"] = strings.Split(d.Id(), ":")[1]
	req["ScalingScheduledTaskId"] = strings.Split(d.Id(), ":")[0]
	action := "DeleteScalingScheduledTask"
	otherErrorRetry := 10

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DeleteScheduledTask(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return OtherErrorProcess(&otherErrorRetry, fmt.Errorf("error on  deleting ScalingScheduledTask %q, %s", d.Id(), err1))
		}
	})

}
