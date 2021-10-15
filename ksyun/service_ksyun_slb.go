package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"strings"
	"time"
)

type SlbService struct {
	client *KsyunClient
}

//start slb

func (s *SlbService) ReadLoadBalancers(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeLoadBalancers"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeLoadBalancers(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeLoadBalancers(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("LoadBalancerDescriptions", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadLoadBalancer(d *schema.ResourceData, loadBalancerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if loadBalancerId == "" {
		loadBalancerId = d.Id()
	}
	req := map[string]interface{}{
		"LoadBalancerId.1": loadBalancerId,
	}
	err = addProjectInfo(d, &req, s.client)
	if err != nil {
		return data, err
	}
	results, err = s.ReadLoadBalancers(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("LoadBalancer %s not exist ", loadBalancerId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		data, callErr := s.ReadLoadBalancer(d, "")
		if callErr != nil {
			if !d.IsNewResource() {
				return resource.NonRetryableError(callErr)
			}
			if notFoundError(callErr) {
				return resource.RetryableError(callErr)
			} else {
				return resource.NonRetryableError(fmt.Errorf("error on  reading lb %q, %s", d.Id(), callErr))
			}
		} else {
			extra := map[string]SdkResponseMapping{
				"ProjectId": {
					Field: "project_id",
					FieldRespFunc: func(i interface{}) interface{} {
						v, _ := strconv.Atoi(i.(string))
						return v
					},
				},
			}
			SdkResponseAutoResourceData(d, r, data, extra)
			return nil
		}
	})
}

func (s *SlbService) ReadAndSetLoadBalancers(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "LoadBalancerId",
			Type:    TransformWithN,
		},
		"project_id": {
			mapping: "ProjectId",
			Type:    TransformWithN,
		},
		"vpc_id": {
			mapping: "vpc-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLoadBalancers(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "LoadBalancerName",
		idFiled:     "LoadBalancerId",
		targetField: "lbs",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateLoadBalancerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"load_balancer_state": {
			mapping: "AdminStateUp",
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if data.Get("load_balancer_state") == "start" {
					return true, true
				} else {
					return false, true
				}
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	if req["Type"] == "internal" {
		if _, ok := req["SubnetId"]; !ok {
			return callback, fmt.Errorf(" Lb type is internal,must set SubnetId")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateLoadBalancer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateLoadBalancer(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("LoadBalancerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateLoadBalancerCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) ModifyLoadBalancerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {Ignore: true},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["LoadBalancerId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyLoadBalancer",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyLoadBalancer(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (s *SlbService) ModifyLoadBalancerProjectCall(d *schema.ResourceData, resource *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"project_id": {},
	}
	updateReq, err := SdkRequestAutoMapping(d, resource, true, transform, nil)
	if err != nil {
		return callback, err
	}
	if len(updateReq) > 0 {
		callback = ApiCall{
			param: &updateReq,
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				return resp, ModifyProjectInstanceNew(d.Id(), call.param, client)
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				return err
			},
		}
	}
	return callback, err
}

func (s *SlbService) ModifyLoadBalancer(d *schema.ResourceData, r *schema.Resource) (err error) {
	projectCall, err := s.ModifyLoadBalancerProjectCall(d, r)
	if err != nil {
		return err
	}
	call, err := s.ModifyLoadBalancerCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{projectCall, call}, d, s.client, true)
}

func (s *SlbService) RemoveLoadBalancerCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"LoadBalancerId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteLoadBalancer",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteLoadBalancer(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadBalancer(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *SlbService) RemoveLoadBalancer(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLoadBalancerCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

//start listener

func (s *SlbService) ReadListeners(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeListeners"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeListeners(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeListeners(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("ListenerSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadListener(d *schema.ResourceData, listenerId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if listenerId == "" {
		listenerId = d.Id()
	}
	req := map[string]interface{}{
		"ListenerId.1": listenerId,
	}
	results, err = s.ReadListeners(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Lb listener %s not exist ", listenerId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadListener(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *SlbService) ReadAndSetListeners(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "ListenerId",
			Type:    TransformWithN,
		},
		"load_balancer_id": {
			mapping: "load-balancer-id",
			Type:    TransformWithFilter,
		},
		"certificate_id": {
			mapping: "certificate-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadListeners(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "ListenerName",
		idFiled:     "ListenerId",
		targetField: "listeners",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateListenerCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"enable_http2": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return data.Get("enable_http2"), true
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if req["listener_protocol"] != "HTTPS" {
		delete(req, "EnableHttp2")
	}

	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateListeners",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateListeners(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("ListenerId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateListener(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateListenerCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

//start healthCheck

func (s *SlbService) ReadLoadHealthChecks(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeHealthChecks"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeHealthChecks(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeHealthChecks(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("HealthCheckSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadLoadHealthCheck(d *schema.ResourceData, healthCheckId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if healthCheckId == "" {
		healthCheckId = d.Id()
	}
	req := map[string]interface{}{
		"HealthCheckId.1": healthCheckId,
	}
	results, err = s.ReadLoadHealthChecks(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("HealthCheck %s not exist ", healthCheckId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetHealCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLoadHealthCheck(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	listener, err := s.ReadListener(nil, data["ListenerId"].(string))
	if err != nil {
		return err
	}
	if listener["ListenerProtocol"] == "HTTP" || listener["ListenerProtocol"] == "HTTPS" {
		if data["HostName"] == nil {
			err = d.Set("is_default_host_name", true)
		} else {
			err = d.Set("is_default_host_name", false)
		}
		if err != nil {
			return err
		}
	}
	return d.Set("listener_protocol", listener["ListenerProtocol"])
}

func (s *SlbService) ReadAndSetHealthChecks(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "HealthCheckId",
			Type:    TransformWithN,
		},
		"listener_id": {
			mapping: "listener-id",
			Type:    TransformWithFilter,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadLoadHealthChecks(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "HealthCheckId",
		targetField: "health_checks",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *SlbService) CreateHealthCheckCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"is_default_host_name": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return data.Get("is_default_host_name"), true
			},
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	callback = ApiCall{
		param:  &req,
		action: "ConfigureHealthCheck",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ConfigureHealthCheck(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("HealthCheckId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateHealthCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateHealthCheckCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) ModifyHealthCheckCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"health_check_state": {
			forceUpdateParam: true,
		},
		"url_path": {
			forceUpdateParam: true,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["HealthCheckId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyHealthCheck",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyHealthCheck(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (s *SlbService) ModifyHealthCheck(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyHealthCheckCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveHealthCheckCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"HealthCheckId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteHealthCheck",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteHealthCheck(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLoadHealthCheck(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading health check when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *SlbService) RemoveHealthCheck(d *schema.ResourceData) (err error) {
	call, err := s.RemoveHealthCheckCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

// LbRule start

func (s *SlbService) ReadLbRules(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.slbconn
	action := "DescribeRules"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeRules(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeRules(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("RuleSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *SlbService) ReadLbRule(d *schema.ResourceData, lbRuleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if lbRuleId == "" {
		lbRuleId = d.Id()
	}
	req := map[string]interface{}{
		"RuleId.1": lbRuleId,
	}
	results, err = s.ReadLbRules(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Lb Rule %s not exist ", lbRuleId)
	}
	return data, err
}

func (s *SlbService) ReadAndSetLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadLbRule(d, "")
	if err != nil {
		return err
	}
	hostname, _ := getSdkValue("HealthCheck.HostName", data)
	if hostname == nil {
		data["HealthCheck"].(map[string]interface{})["IsDefaultHostName"] = true
	} else {
		data["HealthCheck"].(map[string]interface{})["IsDefaultHostName"] = false
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	if err != nil {
		return err
	}
	return err
}

func (s *SlbService) CreateLbRuleCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"Session." 和 "HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	//非同步模式下 需要检查一下
	if req["ListenerSync"] == "off" && req["CookieType"] == "RewriteCookie" {
		if _, ok := req["CookieName"]; !ok {
			return callback, fmt.Errorf("Session.CookieType  is RewriteCookie,must set CookieName")
		}
	}
	callback = ApiCall{
		param:  &req,
		action: "CreateSlbRule",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateSlbRule(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Rule.RuleId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *SlbService) CreateLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateLbRuleCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) ModifyLbRuleCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"session": {
			Type: TransformListUnique,
		},
		"health_check": {
			Type: TransformListUnique,
		},
	}
	req, err := SdkRequestAutoMapping(d, r, true, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})
	if err != nil {
		return callback, err
	}
	//特殊处理下"Session." 和 "HealthCheck."
	for k, v := range req {
		if strings.HasPrefix(k, "Session.") {
			req[strings.Replace(k, "Session.", "", -1)] = v
			delete(req, k)
		}
		if strings.HasPrefix(k, "HealthCheck.") {
			req[strings.Replace(k, "HealthCheck.", "", -1)] = v
			delete(req, k)
		}
	}
	if len(req) > 0 {
		req["RuleId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifySlbRule",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.slbconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifySlbRule(call.param)
				return resp, err
			},
			afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
				logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
				return err
			},
		}
	}
	return callback, err
}

func (s *SlbService) ModifyLbRule(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyLbRuleCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *SlbService) RemoveLbRuleCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"RuleId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteRule",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.slbconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteRule(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadLbRule(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading lb rule when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *SlbService) RemoveLbRule(d *schema.ResourceData) (err error) {
	call, err := s.RemoveLbRuleCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
