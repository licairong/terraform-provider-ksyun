package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunListenerAssociateAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerAssociateAclCreate,
		Read:   resourceKsyunListenerAssociateAclRead,
		Delete: resourceKsyunListenerAssociateAclDelete,
		Importer: &schema.ResourceImporter{
			State: importLoadBalancerAclAssociate,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}
func resourceKsyunListenerAssociateAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateLoadBalancerAclAssociate(d, resourceKsyunListenerAssociateAcl())
	if err != nil {
		return fmt.Errorf("error on creating listener acl associate %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerAssociateAclRead(d, meta)
}

func resourceKsyunListenerAssociateAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetLoadBalancerAclAssociate(d, resourceKsyunListenerAssociateAcl())
	if err != nil {
		return fmt.Errorf("error on reading  listener acl associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerAssociateAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveLoadBalancerAclAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting listener acl associate %q, %s", d.Id(), err)
	}
	return err
}
