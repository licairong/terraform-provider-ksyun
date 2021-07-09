package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

// redis security group allocate
func resourceRedisSecurityGroupAllocate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupAllocateCreate,
		Delete: resourceRedisSecurityGroupAllocateDelete,
		Read:   resourceRedisSecurityGroupAllocateRead,
		Update: resourceRedisSecurityGroupAllocateUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				err = d.Set("security_group_id", d.Id())
				if err != nil {
					return nil, err
				}
				d.SetId(d.Id() + "-allocate")
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cache_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceRedisSecurityGroupAllocateCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		err error
	)

	transform := map[string]SdkReqTransform{
		"cache_ids": {
			mapping: "CacheId",
			Type:    TransformWithN,
		},
	}
	err = processRedisSecurityGroupRuleAllocate(d, meta, transform, false, d.Get("security_group_id").(string))
	if err != nil {
		return fmt.Errorf("error on allocate redis security group: %s", err)
	}

	d.SetId(d.Get("security_group_id").(string) + "-allocate")

	return resourceRedisSecurityGroupAllocateRead(d, meta)
}

func resourceRedisSecurityGroupAllocateUpdate(d *schema.ResourceData, meta interface{}) error {
	//cache_ids
	if d.HasChange("cache_ids") {
		var (
			err      error
			oldArray []string
			newArray []string
			add      []interface{}
			del      []interface{}
		)
		_, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
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
		err = processRedisSecurityGroupRuleAllocate(d, meta, transformAdd, false, d.Get("security_group_id").(string))
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
		err = processRedisSecurityGroupRuleAllocate(d, meta, transformDel, true, d.Get("security_group_id").(string))
		if err != nil {
			return err
		}

	}

	return resourceRedisSecurityGroupAllocateRead(d, meta)
}

func resourceRedisSecurityGroupAllocateDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		err error
	)
	_, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	createReq["AvailableZone"] = d.Get("available_zone")
	createReq["SecurityGroupId"] = d.Get("security_group_id")
	ids := SchemaSetToStringSlice(d.Get("cache_ids"))
	for i, id := range ids {
		createReq[fmt.Sprintf("%v%v", "CacheId.", i+1)] = id
	}
	action := "DeallocateSecurityGroup"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, createReq)
		resp, err := conn.DeallocateSecurityGroup(&createReq)
		logger.Debug(logger.RespFormat, action, createReq, *resp, err)
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

func readRedisSecurityGroupAllocate(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	var (
		resp      *map[string]interface{}
		err       error
		instances []interface{}
	)
	currentCount := int64(0)
	total := int64(0)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	readReq["SecurityGroupId"] = d.Get("security_group_id")
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

func resourceRedisSecurityGroupAllocateRead(d *schema.ResourceData, meta interface{}) error {
	data, err := readRedisSecurityGroupAllocate(d, meta)
	if err != nil {
		return err
	}
	extra := map[string]SdkResponseMapping{
		"list": {
			Field: "cache_ids",
			FieldRespFunc: func(i interface{}) interface{} {
				var caches []string
				for _, v := range i.([]interface{}) {
					caches = append(caches, v.(map[string]interface{})["id"].(string))
				}
				return caches
			},
		},
	}
	SdkResponseAutoResourceData(d, resourceRedisSecurityGroupAllocate(), data, extra)
	return nil
}

func processRedisSecurityGroupRuleAllocate(d *schema.ResourceData, meta interface{}, transform map[string]SdkReqTransform, isUpdate bool, sgId string) error {
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
