package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunRegisterBackendServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunRegisterBackendServerCreate,
		Read:   resourceKsyunRegisterBackendServerRead,
		Update: resourceKsyunRegisterBackendServerUpdate,
		Delete: resourceKsyunRegisterBackendServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backend_server_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backend_server_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				Default:      1,
			},
			"register_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_ip": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"real_server_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_server_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunRegisterBackendServerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on creating backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterBackendServerRead(d, meta)
}

func resourceKsyunRegisterBackendServerRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on reading backend server %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunRegisterBackendServerUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyBackendServer(d, resourceKsyunRegisterBackendServer())
	if err != nil {
		return fmt.Errorf("error on updating backend server %q, %s", d.Id(), err)
	}
	return resourceKsyunRegisterBackendServerRead(d, meta)
}

func resourceKsyunRegisterBackendServerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveBackendServer(d)
	if err != nil {
		return fmt.Errorf("error on deleting backend server %q, %s", d.Id(), err)
	}
	return err
}
