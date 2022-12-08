package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func resourceKsyunSlbRule() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}
	}
	return &schema.Resource{
		Create: resourceKsyunSlbRuleCreate,
		Read:   resourceKsyunSlbRuleRead,
		Update: resourceKsyunSlbRuleUpdate,
		Delete: resourceKsyunSlbRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_header_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_sync": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on",
					"off",
				}, false),
				Default: "on",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RoundRobin",
					"LeastConnections",
				}, false),
				Default:          "RoundRobin",
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
			},

			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
			},
			"session": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_state": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"start",
								"stop",
							}, false),
							Default: "start",
						},
						"session_persistence_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 86400),
							Default:      7200,
						},
						"cookie_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ImplantCookie",
								"RewriteCookie",
							}, false),
							Default: "ImplantCookie",
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				DiffSuppressFunc: lbRuleDiffSuppressFunc,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunSlbRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		return fmt.Errorf("error on creating lb rule %q, %s", d.Id(), err)
	}
	return resourceKsyunSlbRuleRead(d, meta)
}

func resourceKsyunSlbRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading lb rule %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSlbRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLbRule(d, resourceKsyunSlbRule())
	if err != nil {
		return fmt.Errorf("error on updating lb rule %q, %s", d.Id(), err)
	}
	return resourceKsyunSlbRuleRead(d, meta)
}

func resourceKsyunSlbRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLbRule(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb rule %q, %s", d.Id(), err)
	}
	return err
}
