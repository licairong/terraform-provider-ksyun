package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

type KcmService struct {
	client *KsyunClient
}

func (s *KcmService) ReadCertificates(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	conn := s.client.kcmconn
	action := "DescribeCertificates"
	logger.Debug(logger.ReqFormat, action, condition)
	if condition == nil {
		resp, err = conn.DescribeCertificates(nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = conn.DescribeCertificates(&condition)
		if err != nil {
			return data, err
		}
	}

	results, err = getSdkValue("CertificateSet", *resp)
	if err != nil {
		return data, err
	}
	data = results.([]interface{})
	return data, err
}

func (s *KcmService) ReadCertificate(d *schema.ResourceData, certificateId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
	)
	if certificateId == "" {
		certificateId = d.Id()
	}
	req := map[string]interface{}{
		"CertificateId.1": certificateId,
	}
	results, err = s.ReadCertificates(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		data = v.(map[string]interface{})
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Certificate %s not exist ", certificateId)
	}
	return data, err
}

func (s *KcmService) ReadAndSetCertificate(d *schema.ResourceData, r *schema.Resource) (err error) {
	data, err := s.ReadCertificate(d, "")
	if err != nil {
		return err
	}
	SdkResponseAutoResourceData(d, r, data, nil)
	return err
}

func (s *KcmService) ReadAndSetCertificates(d *schema.ResourceData, r *schema.Resource) (err error) {
	transform := map[string]SdkReqTransform{
		"ids": {
			mapping: "CertificateId",
			Type:    TransformWithN,
		},
	}
	req, err := mergeDataSourcesReq(d, r, transform)
	if err != nil {
		return err
	}
	data, err := s.ReadCertificates(req)
	if err != nil {
		return err
	}

	return mergeDataSourcesResp(d, r, ksyunDataSource{
		collection:  data,
		nameField:   "CertificateName",
		idFiled:     "CertificateId",
		targetField: "certificates",
		extra:       map[string]SdkResponseMapping{},
	})
}

func (s *KcmService) CreateCertificateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	transform := map[string]SdkReqTransform{
		"private_key": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return strings.Replace(fmt.Sprintf("%s", data.Get("private_key").(string)), "\n", "\\n", -1), true
			},
		},
		"public_key": {
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				return strings.Replace(fmt.Sprintf("%s", data.Get("public_key").(string)), "\n", "\\n", -1), true
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
		action: "CreateCertificate",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcmconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.CreateCertificate(call.param)
			return resp, err
		},
		afterCall: func(d *schema.ResourceData, client *KsyunClient, resp *map[string]interface{}, call ApiCall) (err error) {
			logger.Debug(logger.RespFormat, call.action, *(call.param), *resp)
			id, err := getSdkValue("Certificate.CertificateId", *resp)
			if err != nil {
				return err
			}
			d.SetId(id.(string))
			return err
		},
	}
	return callback, err
}

func (s *KcmService) CreateCertificate(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.CreateCertificateCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KcmService) ModifyCertificateCall(d *schema.ResourceData, r *schema.Resource) (callback ApiCall, err error) {
	req, err := SdkRequestAutoMapping(d, r, true, nil, nil)
	if err != nil {
		return callback, err
	}
	if len(req) > 0 {
		req["CertificateId"] = d.Id()
		callback = ApiCall{
			param:  &req,
			action: "ModifyCertificate",
			executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
				conn := client.kcmconn
				logger.Debug(logger.RespFormat, call.action, *(call.param))
				resp, err = conn.ModifyCertificate(call.param)
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

func (s *KcmService) ModifyCertificate(d *schema.ResourceData, r *schema.Resource) (err error) {
	call, err := s.ModifyCertificateCall(d, r)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}

func (s *KcmService) RemoveCertificateCall(d *schema.ResourceData) (callback ApiCall, err error) {
	removeReq := map[string]interface{}{
		"CertificateId": d.Id(),
	}
	callback = ApiCall{
		param:  &removeReq,
		action: "DeleteCertificate",
		executeCall: func(d *schema.ResourceData, client *KsyunClient, call ApiCall) (resp *map[string]interface{}, err error) {
			conn := client.kcmconn
			logger.Debug(logger.RespFormat, call.action, *(call.param))
			resp, err = conn.DeleteCertificate(call.param)
			return resp, err
		},
		callError: func(d *schema.ResourceData, client *KsyunClient, call ApiCall, baseErr error) error {
			return resource.Retry(15*time.Minute, func() *resource.RetryError {
				_, callErr := s.ReadCertificate(d, "")
				if callErr != nil {
					if notFoundError(callErr) {
						return nil
					} else {
						return resource.NonRetryableError(fmt.Errorf("error on  reading certificate when delete %q, %s", d.Id(), callErr))
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

func (s *KcmService) RemoveCertificate(d *schema.ResourceData) (err error) {
	call, err := s.RemoveCertificateCall(d)
	if err != nil {
		return err
	}
	return ksyunApiCallNew([]ApiCall{call}, d, s.client, true)
}
