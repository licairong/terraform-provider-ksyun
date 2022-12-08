package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
)

func resourceKsyunSecurityGroupEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunSecurityGroupEntryCreate,
		Read:   resourceKsyunSecurityGroupEntryRead,
		Update: resourceKsyunSecurityGroupEntryUpdate,
		Delete: resourceKsyunSecurityGroupEntryDelete,
		Importer: &schema.ResourceImporter{
			State: importSecurityGroupEntry,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
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
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in",
					"out",
				}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip",
					"tcp",
					"udp",
					"icmp",
				}, false),
			},
			"icmp_type": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
			},
			"icmp_code": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
			},
			"port_range_from": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
			},
			"port_range_to": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.IntBetween(1, 65535),
				DiffSuppressFunc: securityGroupEntryDiffSuppressFunc,
			},
			"security_group_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunSecurityGroupEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.CreateSecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		return fmt.Errorf("error on creating security group entry %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupEntryRead(d, meta)
}
func resourceKsyunSecurityGroupEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ReadAndSetSecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading security group entry %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunSecurityGroupEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.ModifySecurityGroupEntry(d, resourceKsyunSecurityGroupEntry())
	if err != nil {
		return fmt.Errorf("error on updating security group entry %q, %s", d.Id(), err)
	}
	return resourceKsyunSecurityGroupEntryRead(d, meta)
}

func resourceKsyunSecurityGroupEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcService := VpcService{meta.(*KsyunClient)}
	err = vpcService.RemoveSecurityGroupEntry(d)
	if err != nil {
		return fmt.Errorf("error on deleting security group entry %q, %s", d.Id(), err)
	}
	return err
}
