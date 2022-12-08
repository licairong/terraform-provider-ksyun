package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceKsyunNetworkAcl() *schema.Resource {
	entry := resourceKsyunNetworkAclEntry().Schema
	for k, v := range entry {
		if k == "network_acl_id" {
			delete(entry, k)
		} else {
			v.ForceNew = false
		}
	}
	return &schema.Resource{
		Create: resourceKsyunNetworkAclCreate,
		Read:   resourceKsyunNetworkAclRead,
		Delete: resourceKsyunNetworkAclDelete,
		Update: resourceKsyunNetworkAclUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: networkAclEntryCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_acl_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: entry,
				},
				Set: networkAclEntryHash,
			},
			"network_acl_id": {
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

func resourceKsyunNetworkAclCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		return fmt.Errorf("error on creating network acl %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclRead(d, meta)
}

func resourceKsyunNetworkAclRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading network acl  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyNetworkAcl(d, resourceKsyunNetworkAcl())
	if err != nil {
		return fmt.Errorf("error on updating network acl %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclRead(d, meta)
}

func resourceKsyunNetworkAclDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAcl(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl %q, %s", d.Id(), err)
	}
	return err
}
