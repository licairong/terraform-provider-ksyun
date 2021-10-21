package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunLoadBalancerAcl() *schema.Resource {
	entry := resourceKsyunLoadBalancerAclEntry().Schema
	for k, v := range entry {
		if k == "load_balancer_acl_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclCreate,
		Read:   resourceKsyunLoadBalancerAclRead,
		Update: resourceKsyunLoadBalancerAclUpdate,
		Delete: resourceKsyunLoadBalancerAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_acl_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ipv4",
					"ipv6",
				}, false),
				Default: "ipv4",
			},
			"load_balancer_acl_entry_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Computed: true,
				Set:      loadBalancerAclEntryHash,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunLoadBalancerAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancerAcl(d, resourceKsyunLoadBalancerAcl())
	if err != nil {
		return fmt.Errorf("error on creating lb acl %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclRead(d, meta)
}

func resourceKsyunLoadBalancerAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancerAcl(d, resourceKsyunLoadBalancerAcl())
	if err != nil {
		return fmt.Errorf("error on reading lb acl %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLoadBalancerAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLoadBalancerAcl(d, resourceKsyunLoadBalancerAcl())
	if err != nil {
		return fmt.Errorf("error on updating lb acl %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclRead(d, meta)
}

func resourceKsyunLoadBalancerAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancerAcl(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb acl  %q, %s", d.Id(), err)
	}
	return err
}
