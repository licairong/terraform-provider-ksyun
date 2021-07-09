package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	"github.com/KscSDK/ksc-sdk-go/service/ebs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

func resourceKsyunVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeCreate,
		Update: resourceKsyunVolumeUpdate,
		Read:   resourceKsyunVolumeRead,
		Delete: resourceKsyunVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"online_resize": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HourlyInstantSettlement",
					"Daily",
				}, false),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"volume_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).ebsconn
	transform := map[string]SdkReqTransform{
		"online_resize": {Ignore: true},
	}
	createReq, err := SdkRequestAutoMapping(d, resourceKsyunVolume(), false, transform, nil, SdkReqParameter{
		false,
	})
	action := "CreateVolume"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.CreateVolume(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating volume: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	id, ok := (*resp)["VolumeId"]
	if !ok {
		return fmt.Errorf("error on creating volume : no id found")
	}
	idRes, ok := id.(string)
	if !ok {
		return fmt.Errorf("error on creating volume : no id found")
	}
	d.SetId(idRes)
	err = checkKsyunVolumeStatus(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error on waiting for volume %q complete creating, %s", d.Id(), err)
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func checkKsyunVolumeStatus(d *schema.ResourceData, meta interface{}, timeOut time.Duration) error {
	var (
		err error
	)
	conn := meta.(*KsyunClient).ebsconn
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"available", "in-use"},
		Refresh:    resourceKsyunVolumeStatusRefresh(conn, d.Id(), []string{"available", "in-use"}),
		Timeout:    timeOut,
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()
	return err
}

func resourceKsyunVolumeStatusRefresh(conn *ebs.Ebs, volumeId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"VolumeId.1": volumeId}
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.DescribeVolumes(&req)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		volumeList, ok := (*resp)["Volumes"]
		if !ok {
			return nil, "", fmt.Errorf("no volume get")
		}
		volumes, ok1 := volumeList.([]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no volume get")
		}
		if volumes == nil || len(volumes) < 1 {
			return nil, "", fmt.Errorf("no volume get")
		}
		volume, ok2 := volumes[0].(map[string]interface{})
		if !ok2 {
			return nil, "", fmt.Errorf("no volume get")
		}
		status, ok3 := volume["VolumeStatus"]
		if !ok3 {
			return nil, "", fmt.Errorf("no volume status get")
		}
		logger.Debug(logger.RespFormat, action, status)
		if status == "error" {
			return nil, "", fmt.Errorf("volume error")
		}
		for k, v := range target {
			if v == status {
				return resp, status.(string), nil
			}
			if k == len(target)-1 {
				status = statusPending
			}
		}
		return resp, status.(string), nil
	}
}

func resourceKsyunVolumeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).ebsconn
	readReq := make(map[string]interface{})
	readReq["VolumeId.1"] = d.Id()
	action := "DescribeVolumes"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err = conn.DescribeVolumes(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	volumeList, ok := (*resp)["Volumes"]
	if !ok {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	volumes, ok1 := volumeList.([]interface{})
	if !ok1 {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	if volumes == nil || len(volumes) != 1 {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	SdkResponseAutoResourceData(d, resourceKsyunVolume(), volumes[0], nil)
	if _, ok = d.GetOk("online_resize"); !ok {
		if volumes[0].(map[string]interface{})["VolumeStatus"].(string) == "available" {
			err = d.Set("online_resize", false)
		} else {
			err = d.Set("online_resize", true)
		}
	}
	return err
}

func modifyKsyunVolumeInfo(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).ebsconn
	transform := map[string]SdkReqTransform{
		"volume_name": {},
		"volume_desc": {},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunVolume(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on update volume info %q, %s", d.Id(), err)
	}
	if len(req) > 0 {
		err = checkKsyunVolumeStatus(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req["VolumeId"] = d.Id()
		action := "ModifyVolume"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.ModifyVolume(&req)
		if err != nil {
			return fmt.Errorf("error on update volume name %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func modifyKsyunVolumeSize(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	conn := meta.(*KsyunClient).ebsconn
	transform := map[string]SdkReqTransform{
		"size": {},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunVolume(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on update volume size %q, %s", d.Id(), err)
	}
	if len(req) > 0 {
		err = checkKsyunVolumeStatus(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		if status, ok := d.GetOk("volume_status"); ok && status.(string) == "available" && d.Get("online_resize").(bool) {
			return fmt.Errorf("error on resize volume %q, status is available not support online_resize", d.Id())
		}
		o, n := d.GetChange("size")
		if o.(int) > n.(int) {
			return fmt.Errorf("error on resize volume %q, resize not support decrease", d.Id())
		}

		req["OnlineResize"] = d.Get("online_resize")
		req["VolumeId"] = d.Id()
		action := "ResizeVolume"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = conn.ResizeVolume(&req)
		if err != nil {
			return fmt.Errorf("error on resize volume %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		err = checkKsyunVolumeStatus(d, meta, d.Timeout(schema.TimeoutUpdate))
	}
	return err
}

func resourceKsyunVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	err = modifyKsyunVolumeInfo(d, meta)
	if err != nil {
		return err
	}
	err = modifyKsyunVolumeSize(d, meta)
	if err != nil {
		return err
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func resourceKsyunVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).ebsconn
	deleteReq := make(map[string]interface{})
	deleteReq["VolumeId"] = d.Id()
	deleteReq["ForceDelete"] = "true"
	action := "DeleteVolume"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	resp, err := conn.DeleteVolume(&deleteReq)
	if err != nil {
		return fmt.Errorf("error on delete volume %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		readReq := make(map[string]interface{})
		readReq["VolumeId.1"] = d.Id()
		action := "DescribeVolumes"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err := conn.DescribeVolumes(&readReq)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		logger.Debug(logger.RespFormat, action, readReq, *resp)
		volumeList, ok := (*resp)["Volumes"]
		if !ok {
			return nil
		}
		volumes, ok1 := volumeList.([]interface{})
		if !ok1 {
			return nil
		}
		if volumes == nil || len(volumes) < 1 {
			return nil
		}
		volume, ok2 := volumes[0].(map[string]interface{})
		if !ok2 {
			return nil
		}
		status, ok3 := volume["VolumeStatus"]
		if !ok3 {
			return nil
		}
		if status == "recycling" {
			return nil
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}
