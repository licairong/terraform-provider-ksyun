package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunLoadBalancerAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclEntryCreate,
		Delete: resourceKsyunLoadBalancerAclEntryDelete,
		Update: resourceKsyunLoadBalancerAclEntryUpdate,
		Read:   resourceKsyunLoadBalancerAclEntryRead,
		Importer: &schema.ResourceImporter{
			State: importLoadBalancerAclEntry,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_number": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 32766),
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow",
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ip",
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
				}, false),
			},
			"load_balancer_acl_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunLoadBalancerAclEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on reading lb acl entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunLoadBalancerAclEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on creating lb acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyLoadBalancerAclEntry(d, resourceKsyunLoadBalancerAclEntry())
	if err != nil {
		return fmt.Errorf("error on updating lb acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunLoadBalancerAclEntryRead(d, meta)
}

func resourceKsyunLoadBalancerAclEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancerAclEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting lb acl entry  %q, %s", d.Id(), err)
	}
	return err
}
