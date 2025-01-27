package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

// redis security group
func resourceRedisSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupCreate,
		Delete: resourceRedisSecurityGroupDelete,
		Update: resourceRedisSecurityGroupUpdate,
		Read:   resourceRedisSecurityGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"cache_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceRedisSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		createReq map[string]interface{}
		resp      *map[string]interface{}
		err       error
	)
	transform := map[string]SdkReqTransform{
		"rules":     {Ignore: true},
		"cache_ids": {Ignore: true},
	}

	conn := meta.(*KsyunClient).kcsv1conn
	createReq, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), false, transform, nil, SdkReqParameter{onlyTransform: false})
	action := "CreateSecurityGroup"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.CreateSecurityGroup(&createReq)
	if err != nil {
		return fmt.Errorf("error on create redis security group: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["securityGroupId"].(string))
	}
	//rule
	transform = map[string]SdkReqTransform{
		"rules": {
			mapping: "Cidrs",
			Type:    TransformWithN,
		},
	}
	err = processRedisSecurityGroupRule(d, meta, transform, false, "")
	if err != nil {
		return err
	}
	//cache_id
	transform = map[string]SdkReqTransform{
		"cache_ids": {
			mapping: "CacheId",
			Type:    TransformWithN,
		},
	}
	err = processRedisSecurityGroupAllocate(d, meta, transform, false, "")
	if err != nil {
		return err
	}
	return resourceRedisSecurityGroupRead(d, meta)
}

func resourceRedisSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		err error
	)
	//deallocate cache instance
	err = deallocateSecurityGroup(d, meta, "", nil, true)
	if err != nil {
		return err
	}
	// delete redis security group
	deleteReq := make(map[string]interface{})
	deleteReq["SecurityGroupId.1"] = d.Id()
	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &deleteReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DeleteSecurityGroup(&deleteReq)
			},
		}
		action := "DeleteSecurityGroup"
		logger.Debug(logger.ReqFormat, action, deleteReq)
		_, err = integrationAzConf.integrationRedisAz()
		if err == nil {
			return nil
		}
		_, err = readRedisSecurityGroup(d, meta, "")
		if err != nil {
			if validateRedisSgExists(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}

func resourceRedisSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		updateReq map[string]interface{}
		resp      *map[string]interface{}
		err       error
	)
	transform := map[string]SdkReqTransform{
		"rules":     {Ignore: true},
		"cache_ids": {Ignore: true},
	}
	updateReq, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), true, transform, nil, SdkReqParameter{
		false,
	})

	if len(updateReq) > 0 {
		updateReq["SecurityGroupId"] = d.Id()
		if _, ok := updateReq["Name"]; !ok {
			updateReq["Name"] = d.Get("name")
		}
		if _, ok := updateReq["Description"]; !ok {
			updateReq["Description"] = d.Get("description")
		}
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &updateReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.ModifySecurityGroup(&updateReq)
			},
		}
		action := "ModifySecurityGroup"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on modify redis security group: %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq, *resp)
	}
	err = updateRedisSecurityGroupRules(d, meta, "")
	if err != nil {
		return fmt.Errorf("error on modify redis security group: %s", err)
	}
	err = updateRedisSecurityGroupAllocate(d, meta, "")
	if err != nil {
		return fmt.Errorf("error on modify redis security group: %s", err)
	}
	return resourceRedisSecurityGroupRead(d, meta)
}

func resourceRedisSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = readRedisSecurityGroup(d, meta, "")
	if err != nil {
		return err
	}
	extra := make(map[string]SdkResponseMapping)
	data := (*resp)["Data"].(map[string]interface{})
	extra["rules"] = SdkResponseMapping{
		Field: "rules",
		FieldRespFunc: func(i interface{}) interface{} {
			return setRedisSgCidrs(i.([]interface{}), d)
		},
	}
	allocate, err := readRedisSecurityGroupAllocate(d, meta, "")
	if err != nil {
		return err
	}
	data["list"] = allocate["list"]
	extra["list"] = SdkResponseMapping{
		Field:         "cache_ids",
		FieldRespFunc: redisSgAllocateFieldRespFunc(d),
	}
	SdkResponseAutoResourceData(d, resourceRedisSecurityGroup(), data, extra)
	return nil
}
