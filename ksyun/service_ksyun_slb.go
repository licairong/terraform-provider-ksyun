package ksyun

import "github.com/terraform-providers/terraform-provider-ksyun/logger"

type SlbService struct {
	client *KsyunClient
}

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
