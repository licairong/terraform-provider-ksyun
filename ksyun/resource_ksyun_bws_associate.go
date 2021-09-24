package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunBandWidthShareAssociate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBandWidthShareAssociateCreate,
		Read:   resourceKsyunBandWidthShareAssociateRead,
		Delete: resourceKsyunBandWidthShareAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importBandWidthShareAssociate,
		},

		Schema: map[string]*schema.Schema{
			"band_width_share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceKsyunBandWidthShareAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.AssociateBandWidthShare(d, resourceKsyunBandWidthShareAssociate())
	if err != nil {
		return fmt.Errorf("error on associate bandWidthShare %q, %s", d.Id(), err)
	}
	return resourceKsyunBandWidthShareAssociateRead(d, meta)
}

func resourceKsyunBandWidthShareAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.ReadAndSetAssociateBandWidthShare(d, resourceKsyunBandWidthShareAssociate())
	if err != nil {
		return fmt.Errorf("error on reading bandWidthShare associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBandWidthShareAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bwsService := BwsService{meta.(*KsyunClient)}
	err = bwsService.DisassociateBandWidthShare(d)
	if err != nil {
		return fmt.Errorf("error on disAssociate bandWidthShare %q, %s", d.Id(), err)
	}
	return err

}
