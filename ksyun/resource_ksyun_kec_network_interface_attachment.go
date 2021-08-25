package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKecNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKecNetworkInterfaceAttachmentCreate,
		Read:   resourceKsyunKecNetworkInterfaceAttachmentRead,
		Delete: resourceKsyunKecNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importKecNetworkInterfaceAttachment,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunKecNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on reading network interface attachement %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunKecNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on creating network interface attachement %q, %s", d.Id(), err)
	}
	return resourceKsyunKecNetworkInterfaceAttachmentRead(d, meta)
}

func resourceKsyunKecNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyNetworkInterfaceAttachment(d, resourceKsyunKecNetworkInterfaceAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting network interface attachement %q, %s", d.Id(), err)
	}
	return err
}
