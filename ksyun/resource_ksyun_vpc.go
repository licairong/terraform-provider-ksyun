package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceKsyunVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVpcCreate,
		Update: resourceKsyunVpcUpdate,
		Read:   resourceKsyunVpcRead,
		Delete: resourceKsyunVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"is_default": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Default:  false,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunVpcCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateVpc(d, resourceKsyunVpc())
	if err != nil {
		return fmt.Errorf("error on creating vpc %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetVpc(d, resourceKsyunVpc())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading vpc %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVpcUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyVpc(d, resourceKsyunVpc())
	if err != nil {
		return fmt.Errorf("error on updating vpc %q, %s", d.Id(), err)
	}
	return resourceKsyunVpcRead(d, meta)
}

func resourceKsyunVpcDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveVpc(d)
	if err != nil {
		return fmt.Errorf("error on deleting vpc %q, %s", d.Id(), err)
	}
	return err
}
