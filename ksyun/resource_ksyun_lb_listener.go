package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerCreate,
		Read:   resourceKsyunListenerRead,
		Update: resourceKsyunListenerUpdate,
		Delete: resourceKsyunListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "start",
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TCP",
				ValidateFunc: validation.StringInSlice([]string{
					"TCP",
					"UDP",
					"HTTP",
					"HTTPS",
				}, false),
				ForceNew: true,
			},
			"certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "RoundRobin",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
					"MasterSlave",
				}, false),
			},

			"enable_http2": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"tls_cipher_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "TlsCipherPolicy1.0",
				ValidateFunc: validation.StringInSlice([]string{
					"TlsCipherPolicy1.0",
					"TlsCipherPolicy1.1",
					"TlsCipherPolicy1.2",
					"TlsCipherPolicy1.2-strict",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},
			"http_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "HTTP1.1",
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP1.0",
					"HTTP1.1",
				}, false),
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"redirect_listener_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbListenerDiffSuppressFunc,
			},

			"session": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "stop",
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
						},
						"session_persistence_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validation.IntBetween(1, 86400),
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ImplantCookie",
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_state": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "start",
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
						},
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      60,
							ValidateFunc: validation.IntBetween(1, 3600),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(1, 3600),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      5,
							ValidateFunc: validation.IntBetween(1, 10),
						},
						"url_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "DEFAULT",
						},
					},
				},
			},

			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunListenerExtra(r map[string]SdkReqTransform, d *schema.ResourceData, forceGet bool) map[string]SdkRequestMapping {
	var extra map[string]SdkRequestMapping
	extra = SdkRequestAutoExtra(r, d, forceGet)
	return extra
}

func resourceKsyunListenerReq(req *map[string]interface{}, isUpdate bool) (map[string]interface{}, error) {
	healthCheckReq := make(map[string]interface{})
	for k, v := range *req {
		if strings.HasPrefix(k, "Session.") {
			newK := strings.Replace(k, "Session.", "", -1)
			(*req)[newK] = v
			delete(*req, k)
		}
		if strings.HasPrefix(k, "HealthCheck.") {
			newK := strings.Replace(k, "HealthCheck.", "", -1)
			healthCheckReq[newK] = v
			delete(*req, k)
		}
	}
	for k, v := range *req {
		if k == "ListenerProtocol" && !isUpdate {
			if v.(string) == "HTTPS" || v.(string) == "HTTP" {
				if v.(string) == "HTTPS" {
					if _, ok := (*req)["CertificateId"]; !ok {
						return healthCheckReq, fmt.Errorf(" certificate_id must set On listener_protocol is HTTPS")
					}
				}
				if _, ok := (*req)["CookieType"]; !ok {
					return healthCheckReq, fmt.Errorf(" cookie_type must set On listener_protocol is HTTPS or HTTP")
				} else {
					cookieType := (*req)["CookieType"]
					if _, ok := (*req)["CookieName"]; !ok {
						if cookieType == "RewriteCookie" {
							return healthCheckReq, fmt.Errorf(" cookie_name must set On listener_protocol is HTTPS or HTTP and cookie_type is RewriteCookie")
						}
					}
				}
				if len(healthCheckReq) > 0 {
					if _, ok := healthCheckReq["UrlPath"]; !ok {
						return healthCheckReq, fmt.Errorf(" url_path must set On listener_protocol is HTTPS or HTTP")
					}
				}

			} else {
				if _, ok := (*req)["CertificateId"]; ok {
					return healthCheckReq, fmt.Errorf(" certificate_id must not set On listener_protocol is not HTTPS")
				}
				if _, ok := (*req)["CookieType"]; ok {
					return healthCheckReq, fmt.Errorf(" cookie_type must not set On listener_protocol is not HTTPS or HTTP")
				}
				if _, ok := (*req)["CookieName"]; ok {
					return healthCheckReq, fmt.Errorf(" cookie_name must not set On listener_protocol is not  HTTPS or HTTP")
				}
				if len(healthCheckReq) > 0 {
					if _, ok := healthCheckReq["UrlPath"]; ok {
						return healthCheckReq, fmt.Errorf(" url_path must not set On listener_protocol is not HTTPS or HTTP")
					}
					if v1, ok := healthCheckReq["HostName"]; ok {
						if v1.(string) != "DEFAULT" {
							return healthCheckReq, fmt.Errorf(" host_name must not set On listener_protocol is not HTTPS or HTTP")
						}

					}
				}

			}

		}
	}

	//default host_name
	if !isUpdate {
		for k, v := range healthCheckReq {
			if k == "HostName" && v.(string) == "DEFAULT" {
				delete(healthCheckReq, k)
			}
		}
	}
	return healthCheckReq, nil
}

func resourceKsyunListenerCreate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	r := resourceKsyunListener()
	req := make(map[string]interface{})

	var extra map[string]SdkReqTransform

	extra = map[string]SdkReqTransform{
		"health_check": {Type: TransformListUnique},
		"session":      {Type: TransformListUnique},
	}

	req, err := SdkRequestAutoMapping(d, r, false, nil, resourceKsyunListenerExtra(extra, d, false))
	if err != nil {
		return fmt.Errorf(" Error CreateListeners : %s", err)
	}
	healthCheckReq, err := resourceKsyunListenerReq(&req, false)
	if err != nil {
		return fmt.Errorf(" Error CreateListeners : %s", err)
	}

	action := "CreateListeners"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := slbconn.CreateListeners(&req)
	if err != nil {
		return fmt.Errorf("Error CreateListeners : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	id, ok := (*resp)["ListenerId"]
	if !ok {
		return fmt.Errorf(" Error CreateListeners : no ListenerId found")
	}
	if len(healthCheckReq) > 0 {
		healthCheckReq["ListenerId"] = id.(string)
		// create healthCheck
		_, err = slbconn.ConfigureHealthCheck(&healthCheckReq)
	}
	if err != nil {
		return fmt.Errorf(" Error CreateListeners : %s", err)
	}

	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf(" Error CreateListeners : no ListenerId found")
	}
	d.SetId(idres)
	if err := d.Set("listener_id", idres); err != nil {
		return err
	}
	return resourceKsyunListenerRead(d, m)
}

func resourceKsyunListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.slbconn

	req := make(map[string]interface{})
	req["ListenerId.1"] = d.Id()
	action := "DescribeListeners"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeListeners(&req)
	if err != nil {
		return fmt.Errorf("error on reading DescribeListeners %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["ListenerSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseDefault("HealthCheck.HostName", "DEFAULT", &(items[0]))
		SdkResponseAutoResourceData(d, resourceKsyunListener(), items[0], nil)
	}
	return nil
}

func resourceKsyunListenerUpdate(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	r := resourceKsyunListener()
	healthCheckId := d.Get("health_check.0.health_check_id").(string)
	var healthCheckReq map[string]interface{}
	var err error
	req := make(map[string]interface{})

	var onlyHealthCheck map[string]SdkReqTransform
	onlyHealthCheck = map[string]SdkReqTransform{
		"health_check": {Type: TransformListUnique},
	}

	var extra map[string]SdkReqTransform

	extra = map[string]SdkReqTransform{
		"health_check": {Type: TransformListUnique},
		"session":      {Type: TransformListUnique},
	}

	if healthCheckId == "" {
		healthCheckReq, err = SdkRequestAutoMapping(d, r, false, onlyHealthCheck, nil)
	} else {
		healthCheckReq, err = SdkRequestAutoMapping(d, r, true, onlyHealthCheck, nil)
	}

	if err != nil {
		return fmt.Errorf(" Error UpdateListeners : %s", err)
	}
	healthCheckReq, err = resourceKsyunListenerReq(&healthCheckReq, true)
	if err != nil {
		return fmt.Errorf(" Error UpdateListeners : %s", err)
	}

	req, err = SdkRequestAutoMapping(d, r, true, nil, resourceKsyunListenerExtra(extra, d, true))
	if err != nil {
		return fmt.Errorf(" Error UpdateListeners : %s", err)
	}
	_, err = resourceKsyunListenerReq(&req, true)
	if err != nil {
		return fmt.Errorf(" Error UpdateListeners : %s", err)
	}

	if len(healthCheckReq) > 0 {
		if healthCheckId == "" {
			//create
			healthCheckReq["ListenerId"] = d.Id()
			_, err = slbconn.ConfigureHealthCheck(&healthCheckReq)
		} else {
			//update
			logger.Debug(logger.ReqFormat, "ModifyHealthCheck", req)
			if v, ok := healthCheckReq["HostName"]; ok {
				if v.(string) == "DEFAULT" || v.(string) == "" {
					healthCheckReq["IsDefaultHostName"] = true
					delete(healthCheckReq, "HostName")
				}
			}
			//force set HealthCheckState if not set
			if _, ok := healthCheckReq["HealthCheckState"]; !ok {
				healthCheckReq["HealthCheckState"] = d.Get("health_check.0.health_check_state")
			}
			healthCheckReq["HealthCheckId"] = healthCheckId
			_, err = slbconn.ModifyHealthCheck(&healthCheckReq)
		}
		if err != nil {
			return fmt.Errorf(" Error UpdateListeners : %s", err)
		}
	}

	if len(req) > 0 {
		req["ListenerId"] = d.Id()
		action := "ModifyListeners"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = slbconn.ModifyListeners(&req)
		if err != nil {
			return fmt.Errorf(" Error UpdateListeners : %s", err)
		}
	}

	return resourceKsyunListenerRead(d, m)
}

func resourceKsyunListenerDelete(d *schema.ResourceData, m interface{}) error {
	slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["ListenerId"] = d.Id()
	/*
		req["LoadBalancerId"] = d.Id()
		_, err := slbconn.DeleteLoadBalancer(&req)
		if err != nil {
			return fmt.Errorf("release Listener error:%v", err)
		}
		return nil
	*/
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		action := "DeleteListeners"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := slbconn.DeleteListeners(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["ListenerId.1"] = d.Id()
		action = "DescribeListeners"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := slbconn.DescribeListeners(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading Listener when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["ListenerSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified Listener %q has not been deleted due to unknown error", d.Id()))
	})
}
