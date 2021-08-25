package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

type VpcService struct {
	client *KsyunClient
}

func (s *VpcService) readNetworkInterfaces(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp                    *map[string]interface{}
		networkInterfaceResults interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeNetworkInterfaces"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeNetworkInterfaces(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeNetworkInterfaces(&condition)
		if err != nil {
			return data, err
		}
	}

	networkInterfaceResults, err = getSdkValue("NetworkInterfaceSet", *resp)
	if err != nil {
		return data, err
	}
	data = networkInterfaceResults.([]interface{})
	return data, err
}

func (s *VpcService) readNetworkInterface(d *schema.ResourceData, instanceId string) (data map[string]interface{}, err error) {
	var (
		networkInterfaceResults []interface{}
	)
	if instanceId == "" {
		instanceId = d.Id()
	}
	req := map[string]interface{}{
		"NetworkInterfaceId.1": instanceId,
	}
	networkInterfaceResults, err = s.readNetworkInterfaces(req)
	if err != nil {
		return data, err
	}
	for _, v := range networkInterfaceResults {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("NetworkInterface %s not exist ", instanceId)
	}
	return data, err
}

func (s *VpcService) removeNetworkInterface(d *schema.ResourceData) (err error) {
	call, err := s.removeNetworkInterfaceCall(d)
	if err != nil {
		return err
	}
	err = ksyunApiCallNew([]apiCall{call}, d, s.client, true)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]apiCall{call}, d, s.client, false)
}

func (s *VpcService) removeNetworkInterfaceCall(d *schema.ResourceData) (callback apiCall, err error) {
	removeReq := map[string]interface{}{
		"NetworkInterfaceId": d.Id(),
	}
	callback = apiCall{
		param:  &removeReq,
		action: "DeleteNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call apiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteNetworkInterface(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call apiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.readNetworkInterface(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading network interface when delete %q, %s", d.Id(), callErr))
					}
				}
				_, callErr = call.executeCall(d, client, call)
				if callErr == nil {
					return nil
				}
				return resource.RetryableError(callErr)
			})
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call apiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *VpcService) createNetworkInterfaceCall(createReq *map[string]interface{}) (callback apiCall, err error) {
	callback = apiCall{
		param:  createReq,
		action: "CreateNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call apiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call apiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			var (
				instanceId interface{}
			)
			if resp != nil {
				instanceId, err = getSdkValue("NetworkInterfaceId", *resp)
				if err != nil {
					return err
				}
				d.SetId(instanceId.(string))
			}
			return err
		},
	}
	return callback, err
}

func (s *VpcService) modifyNetworkInterfaceCall(updateReq *map[string]interface{}) (callback apiCall, err error) {
	callback = apiCall{
		param:  updateReq,
		action: "ModifyNetworkInterface",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call apiCall) (resp *map[string]interface{}, err error) {
			conn := client.vpcconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.ModifyNetworkInterface(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call apiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			return err
		},
	}
	return callback, err
}

func (s *VpcService) readSubnets(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.vpcconn
	action := "DescribeSubnets"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeSubnets(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeSubnets(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("SubnetSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *VpcService) readSubnet(d *schema.ResourceData, subnetId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if subnetId == "" {
		subnetId = d.Id()
	}
	req := map[string]interface{}{
		"SubnetId.1": subnetId,
	}
	results, err = s.readSubnets(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Subnet %s not exist ", subnetId)
	}
	return data, err
}
