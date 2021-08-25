package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceKsyunLb() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLbCreate,
		Read:   resourceKsyunLbRead,
		Update: resourceKsyunLbUpdate,
		Delete: resourceKsyunLbDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"load_balancer_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateLbType,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateLbState,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_waf": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceKsyunLbCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		req["VpcId"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("load_balancer_name"); ok {
		req["LoadBalancerName"] = fmt.Sprintf("%v", v)
	} else {
		req["LoadBalancerName"] = resource.PrefixedUniqueId("tf-lb-")
	}
	if v, ok := d.GetOk("load_balancer_state"); ok {
		if v == "start" {
			req["AdminStateUp"] = true
		} else {
			req["AdminStateUp"] = false
		}

	}
	internalFlag := false
	if v, ok := d.GetOk("type"); ok {
		if v == "internal" {
			internalFlag = true
		}
		req["Type"] = fmt.Sprintf("%v", v)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		if !internalFlag {
			return fmt.Errorf(" Error CreateLoadBalancer : public lb can not set subnet id ")
		}
		req["SubnetId"] = fmt.Sprintf("%v", v)
	} else if internalFlag {
		return fmt.Errorf(" Error CreateLoadBalancer : internal lb must set subnet id ")
	}
	if v, ok := d.GetOk("private_ip_address"); ok {
		if !internalFlag {
			return fmt.Errorf(" Error CreateLoadBalancer : public lb can not set private ip address  ")
		}
		req["PrivateIpAddress"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("project_id"); ok {
		req["ProjectId"] = fmt.Sprintf("%v", v)
	}
	action := "CreateLoadBalancer"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := slbconn.CreateLoadBalancer(&req)
	if err != nil {
		return fmt.Errorf("Error CreateLoadBalancer : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	id, ok := (*resp)["LoadBalancerId"]
	if !ok {
		return fmt.Errorf(" Error CreateLoadBalancer : no LoadBalancerId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf(" Error CreateLoadBalancer : no LoadBalancerId found")
	}
	if err := d.Set("load_balancer_id", idres); err != nil {
		return err
	}
	d.SetId(idres)
	return resourceKsyunLbRead(d, m)
}

func resourceKsyunLbRead(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId.1"] = d.Id()
	err := addProjectInfo(d, &req, m.(*KsyunClient))
	if err != nil {
		return fmt.Errorf("error on DescribeLoadBalancers  %q, %s", d.Id(), err)
	}
	action := "DescribeLoadBalancers"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.DescribeLoadBalancers(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeLoadBalancers : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	if resp != nil {
		items, ok := (*resp)["LoadBalancerDescriptions"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SetDByResp(d, items[0], slbKeys, map[string]bool{})
		//the api return is string,but the resource set is int, so transfer it
		projectId, _ := getSdkValue("ProjectId", items[0])
		p, _ := strconv.Atoi(projectId.(string))
		_ = d.Set("project_id", p)
		if t, _ := getSdkValue("Type", items[0]); t == "internal" {
			ip, _ := getSdkValue("PublicIp", items[0])
			err1 := d.Set("private_ip_address", ip)
			if err1 != nil {
				log.Println(err1.Error())
				panic("ERROR: " + err1.Error())
			}
			err2 := d.Set("public_ip", nil)
			if err2 != nil {
				log.Println(err2.Error())
				panic("ERROR: " + err2.Error())
			}
		}
	}
	return nil
}

func resourceKsyunLbUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId"] = d.Id()
	if v, ok := d.GetOk("load_balancer_name"); ok {
		req["LoadBalancerName"] = fmt.Sprintf("%v", v)
	}
	if v, ok := d.GetOk("load_balancer_state"); ok {
		req["LoadBalancerState"] = fmt.Sprintf("%v", v)
	} else {
		return fmt.Errorf("cann't change load_balancer_state to empty string")
	}
	// Enable partial attribute modification
	d.Partial(true)
	if d.HasChange("project_id") {
		req["ProjectId"] = fmt.Sprintf("%v", d.Get("project_id"))
		err := ModifyProjectInstance(d.Id(), &req, m)
		if err != nil {
			return fmt.Errorf("error on updating lb, %s", err)
		}
	}
	// Whether the representative has any modifications
	attributeUpdate := false
	if d.HasChange("load_balancer_name") {
		attributeUpdate = true
	}
	if d.HasChange("load_balancer_state") {
		attributeUpdate = true
	}
	if attributeUpdate {
		action := "ModifyLoadBalancer"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.ModifyLoadBalancer(&req)
		if err != nil {
			logger.Debug(logger.AllFormat, action+" first", req, *resp, err)
			if strings.Contains(err.Error(), "400") {
				time.Sleep(time.Second * 2)
				resp, err = slbconn.ModifyLoadBalancer(&req)
				if err != nil {
					return fmt.Errorf("update Slb (%v)error twice:%v", req, err)
				}
			}
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		d.SetPartial("load_balancer_name")
		d.SetPartial("load_balancer_state")
	}
	d.Partial(false)
	return resourceKsyunLbRead(d, m)
}

func resourceKsyunLbDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerId"] = d.Id()
	/*
		_, err := slbconn.DeleteLoadBalancer(&req)
		if err != nil {
			return fmt.Errorf("release Slb error:%v", err)
		}
		return nil
	*/
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancer"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteLoadBalancer(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["LoadBalancerId.1"] = d.Id()
		if pd, ok := d.GetOk("project_id"); ok {
			req["ProjectId.1"] = fmt.Sprintf("%v", pd)
		}
		action = "DescribeLoadBalancers"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeLoadBalancers(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading lb when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)

		itemSet, ok := (*resp)["LoadBalancerDescriptions"]
		if !ok {
			return nil
		}
		items, ok := itemSet.([]interface{})
		if !ok || len(items) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified lb %q has not been deleted due to unknown error", d.Id()))
	})
}
