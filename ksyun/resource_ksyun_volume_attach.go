package ksyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVolumeAttach() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeAttachCreate,
		Read:   resourceKsyunVolumeAttachRead,
		Update: resourceKsyunVolumeAttachUpdate,
		Delete: resourceKsyunVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			State: importVolumeAttach,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"volume_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceKsyunVolumeAttachCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.CreateVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on creating volume attach %q, %s", d.Id(), err)
	}
	return resourceKsyunVolumeAttachRead(d, meta)
}

func resourceKsyunVolumeAttachRead(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ReadAndSetVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on reading volume attach %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVolumeAttachUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ModifyVolumeAttach(d, resourceKsyunVolumeAttach())
	if err != nil {
		return fmt.Errorf("error on updating volume attach %q, %s", d.Id(), err)
	}
	return resourceKsyunVolumeAttachRead(d, meta)
}

func resourceKsyunVolumeAttachDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.RemoveVolumeAttach(d)
	if err != nil {
		return fmt.Errorf("error on deleting volume attach %q, %s", d.Id(), err)
	}
	return err
}
