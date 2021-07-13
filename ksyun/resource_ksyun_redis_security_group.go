package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
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

func updateRedisSecurityGroupRules(d *schema.ResourceData, meta interface{}, sgId string) (err error) {
	//rule
	if d.HasChange("rules") {
		var (
			oldArray []string
			newArray []string
			add      []interface{}
			del      []interface{}
			resp     *map[string]interface{}
		)
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}
		data := (*resp)["Data"].(map[string]interface{})
		rulesMap := make(map[string]interface{})
		//get rule id for del
		if rules, ok := data["rules"]; ok {
			for _, r := range rules.([]interface{}) {
				rule := r.(map[string]interface{})
				rulesMap[rule["cidr"].(string)] = rule["id"]
			}
		}
		o, n := d.GetChange("rules")
		for _, v := range o.(*schema.Set).List() {
			oldArray = append(oldArray, v.(string))
		}
		for _, v := range n.(*schema.Set).List() {
			newArray = append(newArray, v.(string))
		}
		for _, a := range oldArray {
			if _, ok := rulesMap[a]; !ok {
				continue
			}
			exist := false
			for _, b := range newArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				del = append(del, rulesMap[a])
			}
		}
		for _, a := range newArray {
			exist := false
			for _, b := range oldArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				add = append(add, a)
			}
		}
		transformAdd := map[string]SdkReqTransform{
			"rules": {
				mapping: "Cidrs",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(add) > 0 {
						return add, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformAdd, false, sgId)
		if err != nil {
			return err
		}
		transformDel := map[string]SdkReqTransform{
			"rules": {
				mapping: "SecurityGroupRuleId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(del) > 0 {
						return del, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformDel, true, sgId)
		if err != nil {
			return err
		}

	}
	return err
}

func updateRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, sgId string) (err error) {
	//cache_ids
	if d.HasChange("cache_ids") {
		var (
			err      error
			oldArray []string
			newArray []string
			add      []interface{}
			del      []interface{}
		)
		_, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}

		o, n := d.GetChange("cache_ids")
		for _, v := range o.(*schema.Set).List() {
			oldArray = append(oldArray, v.(string))
		}
		for _, v := range n.(*schema.Set).List() {
			newArray = append(newArray, v.(string))
		}
		for _, a := range oldArray {
			exist := false
			for _, b := range newArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				del = append(del, a)
			}
		}
		for _, a := range newArray {
			exist := false
			for _, b := range oldArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				add = append(add, a)
			}
		}
		transformAdd := map[string]SdkReqTransform{
			"cache_ids": {
				mapping: "CacheId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(add) > 0 {
						return add, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupAllocate(d, meta, transformAdd, false, sgId)
		if err != nil {
			return err
		}
		transformDel := map[string]SdkReqTransform{
			"cache_ids": {
				mapping: "CacheId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(del) > 0 {
						return del, true
					}
					return nil, true
				},
			},
		}
		err = processRedisSecurityGroupAllocate(d, meta, transformDel, true, sgId)
		if err != nil {
			return err
		}

	}
	return err
}

func resourceRedisSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		updateReq map[string]interface{}
		resp      *map[string]interface{}
		err       error
	)
	transform := map[string]SdkReqTransform{
		"rules": {Ignore: true},
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

func readRedisSecurityGroup(d *schema.ResourceData, meta interface{}, securityGroupId string) (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	req := make(map[string]interface{})
	if securityGroupId == "" {
		securityGroupId = d.Id()
	}
	req["SecurityGroupId"] = securityGroupId
	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &req,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeSecurityGroup(&req)
		},
		existFn: func(i *map[string]interface{}) bool {
			v, _ := getSdkValue("Data", *i)
			if v == nil || len(v.(map[string]interface{})) == 0 {
				return false
			}
			return true
		},
	}
	action := "DescribeSecurityGroup"
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, fmt.Errorf("error on reading redis security group %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	return resp, err
}

func readRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, sgId string) (map[string]interface{}, error) {
	var (
		resp      *map[string]interface{}
		err       error
		instances []interface{}
	)
	currentCount := int64(0)
	total := int64(0)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	if sgId == "" {
		sgId = d.Id()
	}
	readReq["SecurityGroupId"] = sgId
	readReq["Limit"] = 1000
	readReq["FilterCache"] = true
	for {
		readReq["Offset"] = currentCount
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &readReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				return conn.DescribeInstances(&readReq)
			},
		}

		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return nil, fmt.Errorf("error on reading redis security group allocate instances %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, readReq, *resp)
		data := (*resp)["Data"].(map[string]interface{})
		total = int64(data["total"].(float64))
		lists := data["list"].([]interface{})
		for _, v := range lists {
			instances = append(instances, v)
		}
		currentCount = int64(len(instances))
		if currentCount == total {
			data["list"] = instances
			return data, err
		}
	}
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
			var cidr []string
			for _, v := range i.([]interface{}) {
				cidr = append(cidr, v.(map[string]interface{})["cidr"].(string))
			}
			return cidr
		},
	}
	allocate, err := readRedisSecurityGroupAllocate(d, meta, "")
	if err != nil {
		return err
	}
	data["list"] = allocate["list"]
	extra["list"] = SdkResponseMapping{
		Field:         "cache_ids",
		FieldRespFunc: redisSgAllocateFieldRespFunc(),
	}
	SdkResponseAutoResourceData(d, resourceRedisSecurityGroup(), data, extra)
	return nil
}

func redisSgAllocateFieldRespFunc() FieldRespFunc {
	return func(i interface{}) interface{} {
		var caches []string
		for _, v := range i.([]interface{}) {
			caches = append(caches, v.(map[string]interface{})["id"].(string))
		}
		return caches
	}
}

func validateRedisSgExists(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "not exist")
}

func processRedisSecurityGroupRule(d *schema.ResourceData, meta interface{}, transform map[string]SdkReqTransform, isUpdate bool, sgId string) error {
	var (
		req    map[string]interface{}
		resp   *map[string]interface{}
		err    error
		action string
	)
	req, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), false, transform, nil)
	if len(req) > 0 {
		conn := meta.(*KsyunClient).kcsv1conn
		if sgId == "" {
			sgId = d.Id()
		}
		//read one time and merge available_zone
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		req["SecurityGroupId"] = sgId
		req["AvailableZone"] = d.Get("available_zone")
		if !isUpdate {
			action = "CreateSecurityGroupRule"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.CreateSecurityGroupRule(&req)
			if err != nil {
				return fmt.Errorf("error on add redis security group rules: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		} else {
			action = "DeleteSecurityGroupRule"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.DeleteSecurityGroupRule(&req)
			if err != nil {
				return fmt.Errorf("error on delete redis security group rules: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		}

	}
	return err
}

func processRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}, transform map[string]SdkReqTransform, isUpdate bool, sgId string) error {
	var (
		req    map[string]interface{}
		resp   *map[string]interface{}
		err    error
		action string
	)
	req, err = SdkRequestAutoMapping(d, resourceRedisSecurityGroup(), false, transform, nil)
	if len(req) > 0 {
		conn := meta.(*KsyunClient).kcsv1conn
		if sgId == "" {
			sgId = d.Id()
		}
		//read one time and merge available_zone
		resp, err = readRedisSecurityGroup(d, meta, sgId)
		if err != nil {
			return err
		}
		req["AvailableZone"] = d.Get("available_zone")
		if !isUpdate {
			req["SecurityGroupId.1"] = sgId
			action = "AllocateSecurityGroup"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.AllocateSecurityGroup(&req)
			if err != nil {
				return fmt.Errorf("error on allocate security group to redis: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		} else {
			req["SecurityGroupId"] = sgId
			action = "DeallocateSecurityGroup"
			logger.Debug(logger.ReqFormat, action, req)
			resp, err = conn.DeallocateSecurityGroup(&req)
			if err != nil {
				return fmt.Errorf("error on deallocateSecurityGroup  security group to redis: %s", err)
			}
			logger.Debug(logger.RespFormat, action, req, *resp)
		}

	}
	return err
}

func deallocateSecurityGroup(d *schema.ResourceData, meta interface{}, sgId string, cacheIds []string, all bool) error {
	var (
		err  error
		resp *map[string]interface{}
	)
	_, err = readRedisSecurityGroup(d, meta, sgId)
	if err != nil {
		return err
	}
	if sgId == "" {
		sgId = d.Id()
	}
	deallocateReq := make(map[string]interface{})
	deallocateReq["SecurityGroupId"] = sgId
	if !all {
		for i, id := range cacheIds {
			deallocateReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = id
		}
	} else {
		allocate, err := readRedisSecurityGroupAllocate(d, meta, "")
		if err != nil {
			return err
		}
		for i, a := range allocate["list"].([]interface{}) {
			deallocateReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = a.(map[string]interface{})["id"]
		}
	}

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &deallocateReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DeallocateSecurityGroup(&deallocateReq)
		},
	}
	action := "DeallocateSecurityGroup"
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deallocateReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err == nil {
			data := (*resp)["Data"].([]interface{})
			if len(data) == 0 {
				return nil
			}
		}
		if err != nil && inUseError(err) {
			return resource.RetryableError(err)
		}
		return nil
	})
}
