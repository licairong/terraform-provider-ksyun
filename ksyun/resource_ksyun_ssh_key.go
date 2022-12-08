package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceKsyunSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSSHKeyCreate,
		Read:   resourceKsyunSSHKeyRead,
		Update: resourceKsyunSSHKeyUpdate,
		Delete: resourceKsyunSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_key": {
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
func resourceKsyunSSHKeyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	sksService := SksService{meta.(*KsyunClient)}
	err = sksService.CreateKey(d, resourceKsyunSSHKey())
	if err != nil {
		return fmt.Errorf("error on creating ssh key %q, %s", d.Id(), err)
	}
	return resourceKsyunSSHKeyRead(d, meta)
}

func resourceKsyunSSHKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	sksService := SksService{meta.(*KsyunClient)}
	err = sksService.ReadAndSetKey(d, resourceKsyunSSHKey())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading ssh key %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSSHKeyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	sksService := SksService{meta.(*KsyunClient)}
	err = sksService.ModifyKey(d, resourceKsyunSSHKey())
	if err != nil {
		return fmt.Errorf("error on updating ssh key %q, %s", d.Id(), err)
	}
	return resourceKsyunSSHKeyRead(d, meta)
}

func resourceKsyunSSHKeyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	sksService := SksService{meta.(*KsyunClient)}
	err = sksService.RemoveKey(d)
	if err != nil {
		return fmt.Errorf("error on deleting ssh key %q, %s", d.Id(), err)
	}
	return err
}
