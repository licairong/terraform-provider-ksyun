package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceKsyunNetworkAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNetworkAclEntryCreate,
		Read:   resourceKsyunNetworkAclEntryRead,
		Delete: resourceKsyunNetworkAclEntryDelete,
		Update: resourceKsyunNetworkAclEntryUpdate,
		Importer: &schema.ResourceImporter{
			State: importNetworkAclEntry,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validation.StringIsEmpty,
					validation.IsCIDR,
				),
			},
			"rule_number": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in",
					"out",
				}, false),
				ForceNew: true,
			},
			"rule_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow",
					"deny",
				}, false),
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"tcp",
					"udp",
					"icmp",
				}, false),
				ForceNew: true,
			},
			"icmp_type": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
			},
			"icmp_code": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
			},
			"port_range_from": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
			},
			"port_range_to": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				ForceNew:         true,
				DiffSuppressFunc: networkAclEntryDiffSuppressFunc,
			},
			"network_acl_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunNetworkAclEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on creating network acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclEntryRead(d, meta)
}

func resourceKsyunNetworkAclEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on reading network acl entry  %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunNetworkAclEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifyNetworkAclEntry(d, resourceKsyunNetworkAclEntry())
	if err != nil {
		return fmt.Errorf("error on updating network acl entry %q, %s", d.Id(), err)
	}
	return resourceKsyunNetworkAclEntryRead(d, meta)
}

func resourceKsyunNetworkAclEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveNetworkAclEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting network acl entry %q, %s", d.Id(), err)
	}
	return err
}
