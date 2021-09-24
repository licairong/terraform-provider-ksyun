package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnGatewayCreate,
		Update: resourceKsyunVpnGatewayUpdate,
		Read:   resourceKsyunVpnGatewayRead,
		Delete: resourceKsyunVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"band_width": {
				Type: schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{
					5,
					10,
					20,
					50,
					100,
					200,
				}),
				Required: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Daily",
				}, false),
			},

			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(0, 36),
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunVpnGatewayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on creating vpn gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRead(d, meta)
}

func resourceKsyunVpnGatewayRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on reading vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpnGateway(d, resourceKsyunVpnGateway())
	if err != nil {
		return fmt.Errorf("error on updating vpn gateway  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnGatewayRead(d, meta)
}

func resourceKsyunVpnGatewayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpnGateway(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn gateway  %q, %s", d.Id(), err)
	}
	return err
}
