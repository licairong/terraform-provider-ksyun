package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunSecurityGroup() *schema.Resource {
	entry := resourceKsyunSecurityGroupEntry().Schema
	for k, v := range entry {
		if k == "security_group_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupCreate,
		Update: resourceKsyunSecurityGroupUpdate,
		Read:   resourceKsyunSecurityGroupRead,
		Delete: resourceKsyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_group_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      securityGroupEntryHash,
				Elem: &schema.Resource{
					Schema: entry,
				},
			},
			"security_group_id": {
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

func resourceKsyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on creating security group  %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on reading security group  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySecurityGroup(d, resourceKsyunSecurityGroup())
	if err != nil {
		return fmt.Errorf("error on updating security group  %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupRead(d, meta)
}

func resourceKsyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSecurityGroup(d)
	if err != nil {
		return fmt.Errorf("error on deleting security group  %q, %s", d.Id(), err)
	}
	return err
}
