package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func resourceKsyunBackendServerGroup() *schema.Resource {
	entry := resourceKsyunHealthCheck().Schema
	for k, v := range entry {
		if k == "listener_id" || k == "listener_protocol" || k == "is_default_host_name" || k == "host_name" {
			delete(entry, k)
		} else {
			v.ForceNew = false
			v.DiffSuppressFunc = nil
		}
	}
	entry["host_name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "default",
		ValidateFunc: validation.StringIsNotEmpty,
	}

	return &schema.Resource{
		Create: resourceKsyunBackendServerGroupCreate,
		Read:   resourceKsyunBackendServerGroupRead,
		Update: resourceKsyunBackendServerGroupUpdate,
		Delete: resourceKsyunBackendServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"backend_server_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "backend_server_group",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backend_server_group_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Server",
					"Mirror",
				}, false),
				Default:  "Server",
				ForceNew: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbBackendServerDiffSuppressFunc,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_server_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backend_server_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
func resourceKsyunBackendServerGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on creating backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunBackendServerGroupRead(d, meta)
}

func resourceKsyunBackendServerGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading backend server group %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBackendServerGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyBackendServerGroup(d, resourceKsyunBackendServerGroup())
	if err != nil {
		return fmt.Errorf("error on updating backend server group %q, %s", d.Id(), err)
	}
	return resourceKsyunBackendServerGroupRead(d, meta)
}

func resourceKsyunBackendServerGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveBackendServerGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting backend server group %q, %s", d.Id(), err)
	}
	return err
}
