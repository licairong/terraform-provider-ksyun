package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnCustomerGatewayCreate,
		Update: resourceKsyunVpnCustomerGatewayUpdate,
		Read:   resourceKsyunVpnCustomerGatewayRead,
		Delete: resourceKsyunVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"customer_gateway_address": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validation.IsIPAddress,
				),
			},

			"ha_customer_gateway_address": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validation.IsIPAddress,
				),
			},

			"customer_gateway_mame": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceKsyunVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpnCustomerGateway(d, resourceKsyunVpnCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on creating vpn customer gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnCustomerGatewayRead(d, meta)
}

func resourceKsyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpnCustomerGateway(d, resourceKsyunVpnCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on reading vpn customer gateway  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpnCustomerGateway(d, resourceKsyunVpnCustomerGateway())
	if err != nil {
		return fmt.Errorf("error on updating vpn customer gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnCustomerGatewayRead(d, meta)
}

func resourceKsyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpnCustomerGateway(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn customer gateway  %q, %s", d.Id(), err)
	}
	return err
}
