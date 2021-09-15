package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunVpnTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpnTunnelCreate,
		Update: resourceKsyunVpnTunnelUpdate,
		Read:   resourceKsyunVpnTunnelRead,
		Delete: resourceKsyunVpnTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_tunnel_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"GreOverIpsec",
					"Ipsec",
				}, false),
				Required: true,
				ForceNew: true,
			},

			"vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
			},

			"ha_vpn_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
			},

			"customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
			},

			"ha_customer_gre_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.IsIPAddress,
				),
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"pre_shared_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ike_authen_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"md5",
					"sha",
				}, false),
				Computed: true,
			},

			"ike_dh_group": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					1,
					2,
					5,
				}),
				Computed: true,
			},

			"ike_encry_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"3des",
					"aes",
					"des",
				}, false),
				Computed: true,
			},

			"ipsec_encry_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp-3des",
					"esp-aes",
					"esp-des",
					"esp-null",
					"esp-seal",
				}, false),
				Computed: true,
			},

			"ipsec_authen_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"esp-md5-hmac",
					"esp-sha-hmac",
				}, false),
				Computed: true,
			},

			"ipsec_lifetime_traffic": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(2560, 4608000),
				Computed:     true,
			},

			"ipsec_lifetime_second": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(120, 2592000),
				Computed:     true,
			},
		},
	}
}

func resourceKsyunVpnTunnelCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on creating vpn tunnel  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnTunnelRead(d, meta)
}

func resourceKsyunVpnTunnelRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on reading vpn tunnel  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpnTunnelUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpnTunnel(d, resourceKsyunVpnTunnel())
	if err != nil {
		return fmt.Errorf("error on updating vpn tunnel  %q, %s", d.Id(), err)
	}
	return resourceKsyunVpnTunnelRead(d, meta)
}

func resourceKsyunVpnTunnelDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpnTunnel(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpn tunnel  %q, %s", d.Id(), err)
	}
	return err
}
