package ksyun

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type requestFunc func() (*map[string]interface{}, error)

type IntegrationAzConf struct {
	resourceData *schema.ResourceData
	client       *KsyunClient
	requestFunc  requestFunc
	req          *map[string]interface{}
	field        string
}

func (conf *IntegrationAzConf) integrationAz() (*map[string]interface{}, error) {
	var (
		resp *map[string]interface{}
		err  error
	)
	if v, ok := conf.resourceData.GetOk(conf.field); ok {
		(*conf.req)[Downline2Hump(conf.field)] = v
		return conf.requestFunc()
	}
	conn := conf.client.vpcconn
	req := make(map[string]interface{})
	resp, err = conn.DescribeAvailabilityZones(&req)
	if err != nil {
		return resp, err
	}
	obj, _ := getSdkValue("AvailabilityZoneInfo", *resp)
	for _, az := range obj.([]interface{}) {
		zone := az.(map[string]interface{})["AvailabilityZoneName"]
		(*conf.req)[Downline2Hump(conf.field)] = zone
		resp, err = conf.requestFunc()
		if err == nil {
			_ = conf.resourceData.Set(conf.field, zone)
			return resp, err
		} else {
			continue
		}
	}
	return resp, err

}
