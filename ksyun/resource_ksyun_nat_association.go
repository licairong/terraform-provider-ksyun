package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunNatAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatAssociationCreate,
		Read:   resourceKsyunNatAssociationRead,
		Delete: resourceKsyunNatAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importNatAssociate,
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
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
func resourceKsyunNatAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNatAssociate(d, resourceKsyunNatAssociation())
	if err != nil {
		return fmt.Errorf("error on creating nat associate %q, %s", d.Id(), err)
	}
	return resourceKsyunNatAssociationRead(d, meta)
}

func resourceKsyunNatAssociationRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNatAssociate(d, resourceKsyunNatAssociation())
	if err != nil {
		return fmt.Errorf("error on reading nat associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNatAssociationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNatAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting nat associate %q, %s", d.Id(), err)
	}
	return err
}
