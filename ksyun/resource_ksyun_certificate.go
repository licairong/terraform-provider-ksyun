package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunCertificateCreate,
		Read:   resourceKsyunCertificateRead,
		Update: resourceKsyunCertificateUpdate,
		Delete: resourceKsyunCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunCertificateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	kcmService := KcmService{meta.(*KsyunClient)}
	err = kcmService.CreateCertificate(d, resourceKsyunCertificate())
	if err != nil {
		return fmt.Errorf("error on creating certificate %q, %s", d.Id(), err)
	}
	return resourceKsyunCertificateRead(d, meta)
}

func resourceKsyunCertificateRead(d *schema.ResourceData, meta interface{}) (err error) {
	kcmService := KcmService{meta.(*KsyunClient)}
	err = kcmService.ReadAndSetCertificate(d, resourceKsyunCertificate())
	if err != nil {
		return fmt.Errorf("error on reading certificate %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunCertificateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	kcmService := KcmService{meta.(*KsyunClient)}
	err = kcmService.ModifyCertificate(d, resourceKsyunCertificate())
	if err != nil {
		return fmt.Errorf("error on updating certificate %q, %s", d.Id(), err)
	}
	return resourceKsyunCertificateRead(d, meta)
}

func resourceKsyunCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	kcmService := KcmService{meta.(*KsyunClient)}
	err = kcmService.RemoveCertificate(d)
	if err != nil {
		return fmt.Errorf("error on deleting certificate %q, %s", d.Id(), err)
	}
	return err
}
