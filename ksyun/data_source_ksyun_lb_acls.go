package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceKsyunSlbAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSlbAclsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lb_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_entry_set": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_acl_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"load_balancer_acl_entry_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func dataSourceKsyunSlbAclsRead(d *schema.ResourceData, meta interface{}) error {
	slbService := SlbService{meta.(*KsyunClient)}
	return slbService.ReadAndSetLoadBalancerAcls(d, dataSourceKsyunSlbAcls())
}
