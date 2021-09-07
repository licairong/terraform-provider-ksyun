package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRouteCreate,
		Read:   resourceKsyunRouteRead,
		Delete: resourceKsyunRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"destination_cidr_block": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"route_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"InternetGateway",
					"Tunnel",
					"Host",
					"Peering",
					"DirectConnect",
					"Vpn",
				}, false),
			},
			"tunnel_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"direct_connect_gateway_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vpn_tunnel_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"next_hop_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunRouteCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateRoute(d, resourceKsyunRoute())
	if err != nil {
		return fmt.Errorf("error on creating route %q, %s", d.Id(), err)
	}
	return resourceKsyunRouteRead(d, meta)
}

func resourceKsyunRouteRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetRoute(d, resourceKsyunRoute())
	if err != nil {
		return fmt.Errorf("error on reading route %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunRouteDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveRoute(d)
	if err != nil {
		return fmt.Errorf("error on deleting route %q, %s", d.Id(), err)
	}
	return err
}
