package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceKsyunListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunListenersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"certificate_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"method": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_http2": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"tls_cipher_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"load_balancer_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"health_check": {
							Type:     schema.TypeList,
							MaxItems: 1,
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
										Computed: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"url_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Computed: true,
						},
						"session": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"session_persistence_period": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"session_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"real_server": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_server_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"real_server_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"real_server_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"real_server_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"register_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"listener_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunListenersRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetListeners(d, dataSourceKsyunListeners())
}
