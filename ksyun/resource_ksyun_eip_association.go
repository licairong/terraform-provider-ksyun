package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func resourceKsyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipAssociationCreate,
		Read:   resourceKsyunEipAssociationRead,
		Delete: resourceKsyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: importAddressAssociate,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Ipfwd",
					"Slb",
				}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_band_width_share": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.CreateAddressAssociate(d, resourceKsyunEipAssociation())
	if err != nil {
		return fmt.Errorf("error on creating address association %q, %s", d.Id(), err)
	}
	return resourceKsyunEipAssociationRead(d, meta)
}

func resourceKsyunEipAssociationRead(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.ReadAndSetAddressAssociate(d, resourceKsyunEipAssociation())
	if err != nil {
		if strings.Contains(err.Error(), "not associate in") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading address association %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) (err error) {
	eipService := EipService{meta.(*KsyunClient)}
	err = eipService.RemoveAddressAssociate(d)
	if err != nil {
		return fmt.Errorf("error on deleting address association %q, %s", d.Id(), err)
	}
	return err
}
