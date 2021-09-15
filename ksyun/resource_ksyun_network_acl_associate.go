package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunNetworkAclAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNetworkAclAssociateCreate,
		Read:   resourceKsyunNetworkAclAssociateRead,
		Delete: resourceKsyunNetworkAclAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importNetworkAclAssociate,
		},
		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKsyunNetworkAclAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAclAssociate(d, resourceKsyunNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on creating network acl associate  %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclAssociateRead(d, meta)
}

func resourceKsyunNetworkAclAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAclAssociate(d, resourceKsyunNetworkAclAssociate())
	if err != nil {
		return fmt.Errorf("error on reading network acl associate  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAclAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl associate  %q, %s", d.Id(), err)
	}
	return err
}
