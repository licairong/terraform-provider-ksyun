package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunListenerHostHeader() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerHostHeaderCreate,
		Read:   resourceKsyunListenerHostHeaderRead,
		Update: resourceKsyunListenerHostHeaderUpdate,
		Delete: resourceKsyunListenerHostHeaderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_header": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: hostHeaderDiffSuppressFunc,
			},
			"host_header_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunListenerHostHeaderCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on creating host header %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerHostHeaderRead(d, meta)
}

func resourceKsyunListenerHostHeaderRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on reading host header %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunListenerHostHeaderUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyHostHeader(d, resourceKsyunListenerHostHeader())
	if err != nil {
		return fmt.Errorf("error on updating host header %q, %s", d.Id(), err)
	}
	return resourceKsyunListenerHostHeaderRead(d, meta)
}

func resourceKsyunListenerHostHeaderDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveHostHeader(d)
	if err != nil {
		return fmt.Errorf("error on deleting host header %q, %s", d.Id(), err)
	}
	return err
}
