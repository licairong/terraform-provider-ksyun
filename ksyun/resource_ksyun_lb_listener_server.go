package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunInstancesWithListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstancesWithListenerCreate,
		Read:   resourceKsyunInstancesWithListenerRead,
		Update: resourceKsyunInstancesWithListenerUpdate,
		Delete: resourceKsyunInstancesWithListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"real_server_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"real_server_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"real_server_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"host",
					"DirectConnectGateway",
				}, false),
				Default: "host",
			},
			"instance_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				DiffSuppressFunc: lbRealServerDiffSuppressFunc,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				Default:      1,
			},
			"master_slave_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Master",
					"Slave",
				}, false),
				Default:          "Master",
				DiffSuppressFunc: lbRealServerDiffSuppressFunc,
			},
			"real_server_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"register_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"listener_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunInstancesWithListenerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on creating real server %q, %s", d.Id(), err)
	}
	return resourceKsyunInstancesWithListenerRead(d, meta)
}

func resourceKsyunInstancesWithListenerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on reading real server %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunInstancesWithListenerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyRealServer(d, resourceKsyunInstancesWithListener())
	if err != nil {
		return fmt.Errorf("error on updating real server %q, %s", d.Id(), err)
	}
	return resourceKsyunInstancesWithListenerRead(d, meta)
}

func resourceKsyunInstancesWithListenerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveRealServer(d)
	if err != nil {
		return fmt.Errorf("error on deleting real server %q, %s", d.Id(), err)
	}
	return err
}
