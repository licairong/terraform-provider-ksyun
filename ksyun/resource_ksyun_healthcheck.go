package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunHealthCheckCreate,
		Read:   resourceKsyunHealthCheckRead,
		Update: resourceKsyunHealthCheckUpdate,
		Delete: resourceKsyunHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start",
					"stop",
				}, false),
				Default: "start",
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      5,
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 1000),
				Default:      5,
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 3600),
				Default:      4,
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10),
				Default:      4,
			},
			"url_path": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "/",
				DiffSuppressFunc: lbHealthCheckDiffSuppressFunc,
				ValidateFunc:     validation.StringIsNotEmpty,
			},
			"is_default_host_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: lbHealthCheckDiffSuppressFunc,
			},
			"health_check_id": {
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
func resourceKsyunHealthCheckCreate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.CreateHealthCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on creating health check %q, %s", d.Id(), err)
	}
	return resourceKsyunHealthCheckRead(d, meta)
}

func resourceKsyunHealthCheckRead(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ReadAndSetHealCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on reading health check %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunHealthCheckUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.ModifyHealthCheck(d, resourceKsyunHealthCheck())
	if err != nil {
		return fmt.Errorf("error on updating health check %q, %s", d.Id(), err)
	}
	return resourceKsyunHealthCheckRead(d, meta)
}

func resourceKsyunHealthCheckDelete(d *schema.ResourceData, meta interface{}) (err error) {
	slbService := SlbService{meta.(*KsyunClient)}
	err = slbService.RemoveHealthCheck(d)
	if err != nil {
		return fmt.Errorf("error on deleting health check %q, %s", d.Id(), err)
	}
	return err
}
