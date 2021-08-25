package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKecNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceKecNetworkInterfaceCreate,
		Update: resourceKecNetworkInterfaceUpdate,
		Read:   resourceKecNetworkInterfaceRead,
		Delete: resourceKecNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: kecNetworkInterfaceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"network_interface_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKecNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.createNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on creating network interface %q, %s", d.Id(), err)
	}
	return resourceKecNetworkInterfaceRead(d, meta)
}

func resourceKecNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.modifyNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on updating network interface %q, %s", d.Id(), err)
	}
	return resourceKecNetworkInterfaceRead(d, meta)
}

func resourceKecNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) (err error) {
	kecService := KecService{meta.(*KsyunClient)}
	err = kecService.readAndSetNetworkInterface(d, resourceKsyunKecNetworkInterface())
	if err != nil {
		return fmt.Errorf("error on reading network interface %q, %s", d.Id(), err)
	}
	return err
}

func resourceKecNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	return vpcService.removeNetworkInterface(d)
}
