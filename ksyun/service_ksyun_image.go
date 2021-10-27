package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

type ImageService struct {
	client *KsyunClient
}

func (s *ImageService) readKecImages(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	//return pageQuery(condition, "Limit", "Offset", 20, 0, func(condition map[string]interface{}) ([]interface{}, error) {
	//
	//})
	conn := s.client.kecconn
	action := "DescribeImages"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeImages(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeImages(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("ImagesSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *ImageService) ReadAndSetKecImages(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "ImageId",
			Type:    TransformWithN,
		},
		"platform": {
			Ignore: true,
		},
		"is_public": {
			Ignore: true,
		},
		"image_source": {
			Ignore: true,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.readKecImages(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		idFiled:     "ImageId",
		nameField:   "ImageName",
		targetField: "images",
		extra:       map[string]SdkResponseMapping{},
	}, func(data *schema.ResourceData, m map[string]interface{}) (result map[string]interface{}, flag bool, err error) {
		if name, ok := d.GetOk("platform"); ok {
			flag = true
			if v, ok1 := m["Platform"]; ok1 && v == name {
				result = m
			}
		}
		return result, flag, err
	}, func(data *schema.ResourceData, m map[string]interface{}) (result map[string]interface{}, flag bool, err error) {
		if name, ok := d.GetOk("image_source"); ok {
			flag = true
			if v, ok1 := m["ImageSource"]; ok1 && v == name {
				result = m
			}
		}
		return result, flag, err
	}, func(data *schema.ResourceData, m map[string]interface{}) (result map[string]interface{}, flag bool, err error) {
		if standard, ok := d.GetOk("is_public"); ok {
			flag = true
			if v, ok1 := m["IsPublic"]; ok1 && v == standard {
				result = m
			}
		}
		return result, flag, err
	})
}
