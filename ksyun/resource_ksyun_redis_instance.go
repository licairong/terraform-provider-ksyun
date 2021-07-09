package ksyun

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// instance
func resourceRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceCreate,
		Delete: resourceRedisInstanceDelete,
		Update: resourceRedisInstanceUpdate,
		Read:   resourceRedisInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
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
			"mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"slave_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 8),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bill_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 5, 87}),
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("bill_type"); ok && v == 1 {
						return false
					}
					return true
				},
			},
			"pass_word": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"iam_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"2.8",
					"3.0",
					"4.0",
					"5.0",
				}, false),
			},
			"backup_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("security_group_ids") != nil {
						if sgs, ok := d.Get("security_group_ids").(*schema.Set); ok {
							if (*sgs).Contains(new) {
								err := d.Set("security_group_id", new)
								if err == nil {
									return true
								}
							}
						}
					}
					return false
				},
			},
			"reset_all_parameters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"net_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timing_switch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"sub_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_use": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kcsv1conn
	var (
		resp *map[string]interface{}
		err  error
	)
	// valid parameters ...
	createParam, err := resourceRedisInstanceParameterCheckAndPrepare(d, meta, false)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	r := resourceRedisInstance()
	transform := map[string]SdkReqTransform{
		"reset_all_parameters": {Ignore: true},
		"parameters":           {Ignore: true},
		"protocol": {ValueFunc: func(d *schema.ResourceData) (interface{}, bool) {
			v, ok := d.GetOk("protocol")
			if ok {
				return v, ok
			} else {
				mode := d.Get("mode").(int)
				switch mode {
				case 1:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				case 2:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				default:
					_ = d.Set("protocol", "4.0")
					return "4.0", true
				}
			}
		}},
	}
	//generate req
	createReq, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	//create redis instance
	action := "CreateCacheCluster"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.CreateCacheCluster(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		d.SetId((*resp)["Data"].(map[string]interface{})["CacheId"].(string))
	}
	err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutCreate),"")
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	if len(*createParam) > 0 {
		err = setResourceRedisInstanceParameter(d, meta, createParam)
		if err != nil {
			return fmt.Errorf("error on create Instance: %s", err)
		}
	}

	return resourceRedisInstanceRead(d, meta)
}

func setResourceRedisInstanceParameter(d *schema.ResourceData, meta interface{}, createReq *map[string]interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	(*createReq)["CacheId"] = d.Id()
	(*createReq)["Protocol"] = d.Get("protocol")

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          createReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.SetCacheParameters(createReq)
		},
	}

	action := "SetCacheParameters"
	logger.Debug(logger.ReqFormat, action, *createReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return fmt.Errorf("error on set instance parameter: %s", err)
	}
	logger.Debug(logger.RespFormat, action, *createReq, *resp)
	err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutUpdate),"")
	if err != nil {
		return fmt.Errorf("error on create Instance: %s", err)
	}
	return err
}

func resourceRedisInstanceParameterCheckAndPrepare(d *schema.ResourceData, meta interface{}, isUpdate bool) (*map[string]interface{}, error) {
	var (
		reset bool
		resp  *map[string]interface{}
		err   error
		index int
	)
	conn := meta.(*KsyunClient).kcsv1conn
	req := make(map[string]interface{})

	parameters := make(map[string]string)
	if !isUpdate || (isUpdate && d.HasChange("parameters")) {
		if data, ok := d.GetOk("parameters"); ok {
			for k, v := range data.(map[string]interface{}) {
				parameters[k] = v.(string)
			}
		}
	}

	//reset_all_parameters and parameters Conflict
	if r, ok := d.GetOk("reset_all_parameters"); ok && !isUpdate && r.(bool) && len(parameters) > 0 {
		err = fmt.Errorf("parameters is not empty,can not set reset_all_parameters true")
		return &req, err
	}
	if r, ok := d.GetOk("reset_all_parameters"); ok && isUpdate && r.(bool) {
		if data, ok := d.GetOk("parameters"); ok {
			if len(data.(map[string]interface{})) > 0 {
				err = fmt.Errorf("parameters is not empty,can not set reset_all_parameters true")
				return &req, err
			}
		}
		if d.HasChange("reset_all_parameters") {
			reset = true
		}

	}

	// condition on reset_all_parameters
	if isUpdate && reset {
		logger.Debug(logger.ReqFormat, "DemoTest", reset)
		req["ResetAllParameters"] = reset

		return &req, d.Set("reset_all_parameters", reset)
	}

	//condition on set parameters, check parameter key and value valid
	action := "DescribeCacheDefaultParameters"
	logger.Debug(logger.ReqFormat, action, nil)
	resp, err = conn.DescribeCacheDefaultParameters(&map[string]interface{}{})
	logger.Debug(logger.RespFormat, action, nil, resp)
	if err != nil {
		return &req, fmt.Errorf("error on DescribeCacheDefaultParameters: %s", err)
	}
	data, err := getSdkValue("Data", *resp)
	if err != nil {
		return &req, fmt.Errorf("error on DescribeCacheDefaultParameters: %s", err)
	}
	defaultParameters := make(map[string]interface{})
	for _, item := range data.([]interface{}) {
		name, err := getSdkValue("name", item)
		if err != nil {
			continue
		}
		defaultParameters[name.(string)] = item
	}
	// query current parameter
	cacheParameters := make(map[string]string)
	if d.Id() != "" {
		reqParam := make(map[string]interface{})
		reqParam["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &reqParam,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DescribeCacheParameters(&reqParam)
			},
		}
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return &req, fmt.Errorf("error on DescribeCacheParameters: %s", err)
		}
		data, err = getSdkValue("Data", *resp)
		for _, item := range data.([]interface{}) {
			name, err := getSdkValue("name", item)
			if err != nil {
				continue
			}
			currentValue, err := getSdkValue("currentValue", item)
			cacheParameters[name.(string)] = currentValue.(string)
		}
	}
	for k, v := range parameters {
		if _, ok := defaultParameters[k]; !ok {
			return &req, fmt.Errorf("error on paramerter %v not support", k)
		}
		paramType, err := getSdkValue("validity.type", defaultParameters[k])
		if err != nil {
			continue
		}
		switch paramType.(string) {
		case "enum":
			values, err := getSdkValue("validity.values", defaultParameters[k])
			if err != nil {
				break
			}
			valid := false
			for _, v1 := range values.([]interface{}) {
				if v1.(string) == v {
					valid = true
				}
			}
			if !valid {
				return &req, fmt.Errorf("error on paramerter %v value must in  %v ", k, values)
			}
		case "range":
			minStr, err := getSdkValue("validity.min", defaultParameters[k])
			if err != nil {
				break
			}
			maxStr, err := getSdkValue("validity.max", defaultParameters[k])
			if err != nil {
				break
			}
			min, err := strconv.Atoi(minStr.(string))
			if err != nil {
				break
			}
			max, err := strconv.Atoi(maxStr.(string))
			if err != nil {
				break
			}
			current, err := strconv.Atoi(v)
			if err != nil {
				return &req, fmt.Errorf("error on paramerter %v value must number ", k)
			}
			if current > max || current < min {
				return &req, fmt.Errorf("error on paramerter %v value must in %v,%v ", k, minStr, maxStr)
			}
		case "regexp":
			value, err := getSdkValue("validity.value", defaultParameters[k])
			if err != nil {
				break
			}
			reg := regexp.MustCompile(value.(string))
			if reg == nil {
				continue
			}
			if !reg.MatchString(v) {
				return &req, fmt.Errorf("error on paramerter %v value must match %v ", k, value)
			}
		default:
			break
		}

		if cv, ok := cacheParameters[k]; ok && cv == v {
			continue
		}
		index = index + 1
		req[fmt.Sprintf("%v%v", "Parameters.ParameterName.", index)] = k
		req[fmt.Sprintf("%v%v", "Parameters.ParameterValue.", index)] = v
	}

	return &req, d.Set("reset_all_parameters", reset)
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Id()

	return resource.Retry(20*time.Minute, func() *resource.RetryError {
		var (
			resp *map[string]interface{}
			err error
		)
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &deleteReq,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DeleteCacheCluster(&deleteReq)
			},
		}
		action := "DeleteCacheCluster"
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err = integrationAzConf.integrationRedisAz()
		if err == nil {
			return nil
		}
		logger.Debug(logger.RespFormat, action, deleteReq, resp)
		_, err = describeRedisInstance(d, meta,"")
		if err != nil {
			if validateExists(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(errors.New("deleting"))
	})
}

func validateExists(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "cannot be found") || strings.Contains(strings.ToLower(err.Error()), "invalidaction")
}

func modifyRedisInstanceNameAndProject(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"name":           {},
		"iam_project_id": {mapping: "ProjectId"},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}
	//modify project
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return fmt.Errorf("error on updating instance iam_project_id , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.RenameCacheCluster(&req)
			},
		}
		action := "RenameCacheCluster"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on rename instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func modifyRedisInstancePassword(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"pass_word": {mapping: "Password"},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.UpdatePassword(&req)
			},
		}
		action := "UpdatePassword"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on UpdatePassword instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func modifyRedisInstanceSpec(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"capacity": {},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		req["CacheId"] = d.Id()
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.ResizeCacheCluster(&req)
			},
		}
		action := "ResizeCacheCluster"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on ResizeCacheCluster instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		err = checkRedisInstanceStatus(d, meta, d.Timeout(schema.TimeoutUpdate),"")
		if err != nil {
			return fmt.Errorf("error on ResizeCacheCluster instance %q, %s", d.Id(), err)
		}
	}
	return err
}

func modifyRedisInstanceSg(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
		req  map[string]interface{}
		resp *map[string]interface{}
	)

	transform := map[string]SdkReqTransform{
		"security_group_id": {},
	}
	req, err = SdkRequestAutoMapping(d, resourceRedisInstance(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}

	if len(req) > 0 {
		oldSg, newSg := d.GetChange("security_group_id")
		req["CacheId.1"] = d.Id()
		req["SecurityGroupId"] = oldSg
		integrationAzConf := &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.DeallocateSecurityGroup(&req)
			},
		}
		action := "DeallocateSecurityGroup"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on DeallocateSecurityGroup instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)

		req["SecurityGroupId"] = newSg
		integrationAzConf = &IntegrationRedisAzConf{
			resourceData: d,
			client:       meta.(*KsyunClient),
			req:          &req,
			field:        "available_zone",
			requestFunc: func() (*map[string]interface{}, error) {
				conn := meta.(*KsyunClient).kcsv1conn
				return conn.AllocateSecurityGroup(&req)
			},
		}
		action = "AllocateSecurityGroup"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err = integrationAzConf.integrationRedisAz()
		if err != nil {
			return fmt.Errorf("error on AllocateSecurityGroup instance %q, %s", d.Id(), err)
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
	}
	return err
}

func resourceRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer func(d *schema.ResourceData, meta interface{}) {
		_err := resourceRedisInstanceRead(d, meta)
		if err == nil {
			err = _err
		} else {
			if _err != nil {
				err = fmt.Errorf(err.Error()+" %s", _err)
			}
		}

	}(d, meta)
	// valid parameters ...
	createParam, err := resourceRedisInstanceParameterCheckAndPrepare(d, meta, true)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}

	// rename
	err = modifyRedisInstanceNameAndProject(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	// update password
	err = modifyRedisInstancePassword(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	//sg
	err = modifyRedisInstanceSg(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}
	// resize mem
	err = modifyRedisInstanceSpec(d, meta)
	if err != nil {
		return fmt.Errorf("error on update instance: %s", err)
	}

	// update parameter
	if len(*createParam) > 0 {
		err = setResourceRedisInstanceParameter(d, meta, createParam)
		if err != nil {
			return fmt.Errorf("error on create Instance: %s", err)
		}
	}
	err = d.Set("reset_all_parameters", d.Get("reset_all_parameters"))
	return err
}

func describeRedisInstance(d *schema.ResourceData, meta interface{},id string) (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	queryReq := make(map[string]interface{})
	if id == ""{
		id = d.Id()
	}
	queryReq["CacheId"] = id

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &queryReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheCluster(&queryReq)
		},
	}
	action := "DescribeCacheCluster"
	logger.Debug(logger.ReqFormat, action, queryReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return resp, err
	}
	logger.Debug(logger.RespFormat, action, queryReq, *resp)
	return resp, err
}

func resourceRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		item map[string]interface{}
		resp *map[string]interface{}
		ok   bool
		err  error
	)
	resp, err = describeRedisInstance(d, meta,"")
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
		return nil
	}
	// merge some field
	add := make(map[string]interface{})
	for k, v := range item {
		if k == "az" {
			add["AvailableZone"] = v
		}
	}
	if _, ok = item["slaveNum"]; !ok {
		item["slaveNum"] = 0
	}
	for k, v := range add {
		item[k] = v
	}
	extra := make(map[string]SdkResponseMapping)
	extra["protocol"] = SdkResponseMapping{
		Field: "protocol",
		FieldRespFunc: func(i interface{}) interface{} {
			return strings.Replace(i.(string), "redis ", "", -1)
		},
	}
	extra["size"] = SdkResponseMapping{
		Field: "capacity",
	}
	SdkResponseAutoResourceData(d, resourceRedisInstance(), item, extra)

	//merge securityGroupIds
	err = resourceRedisInstanceSgRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	//merge parameters
	err = resourceRedisInstanceParamRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	return nil
}

func resourceRedisInstanceSgRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)

	querySg := make(map[string]interface{})
	querySg["CacheId"] = d.Id()

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &querySg,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeSecurityGroups(&querySg)
		},
	}

	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return err
	}
	if item, ok := (*resp)["Data"].(map[string]interface{}); ok {
		var itemSetSlice []string
		if sgs, ok := item["list"].([]interface{}); ok {
			for _, sg := range sgs {
				if info, ok := sg.(map[string]interface{}); ok {
					itemSetSlice = append(itemSetSlice, info["securityGroupId"].(string))
				}
			}
		}
		if len(itemSetSlice) > 0 && d.Get("security_group_id") == nil || d.Get("security_group_id").(string) == "" {
			err = d.Set("security_group_id", itemSetSlice[0])
			if err != nil {
				return err
			}
		}
		return d.Set("security_group_ids", itemSetSlice)
	}
	return nil
}

func resourceRedisInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Id()

	integrationAzConf := &IntegrationRedisAzConf{
		resourceData: d,
		client:       meta.(*KsyunClient),
		req:          &readReq,
		field:        "available_zone",
		requestFunc: func() (*map[string]interface{}, error) {
			conn := meta.(*KsyunClient).kcsv1conn
			return conn.DescribeCacheParameters(&readReq)
		},
	}

	action := "DescribeCacheParameters"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err = integrationAzConf.integrationRedisAz()
	if err != nil {
		return fmt.Errorf("error on reading instance parameter %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data) == 0 {
		return nil
	}
	result := make(map[string]interface{})
	parameter := make(map[string]interface{})
	for _, d := range data {
		param := d.(map[string]interface{})
		result[param["name"].(string)] = fmt.Sprintf("%v", param["currentValue"])
	}
	if local, ok := d.GetOk("parameters"); ok {
		for k, v := range local.(map[string]interface{}) {
			if _, ok = result[k]; ok {
				parameter[k] = v
			}
		}
	}

	if err := d.Set("parameters", parameter); err != nil {
		return fmt.Errorf("error set data %v :%v", result, err)
	}
	return nil
}

func checkRedisInstanceStatus(d *schema.ResourceData, meta interface{}, timeout time.Duration,id string) error {
	var err error
	if id == ""{
		id = d.Id()
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"2"},
		Refresh:    stateRefreshForRedis(d, meta, []string{"2"},id),
		Timeout:    timeout,
		Delay:      20 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func stateRefreshForRedis(d *schema.ResourceData, meta interface{}, target []string,id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			resp *map[string]interface{}
			item map[string]interface{}
			ok   bool
			err  error
		)

		resp, err = describeRedisInstance(d, meta,id)
		if err != nil {
			return nil, "", err
		}
		if item, ok = (*resp)["Data"].(map[string]interface{}); !ok {
			return nil, "", fmt.Errorf("no instance information was queried.%s", "")
		}
		status := int(item["status"].(float64))
		serviceStatus := int(item["serviceStatus"].(float64))
		// instance status error
		if status == 0 || status == 99 {
			return nil, "", fmt.Errorf("instance create error,status:%v", status)
		}
		// trade instance status error
		if serviceStatus == 3 {
			return nil, "", fmt.Errorf("instance create error,serviceStatus:%v", serviceStatus)
		}
		state := strconv.Itoa(status)
		for k, v := range target {
			if v == state && serviceStatus == 2 {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}
