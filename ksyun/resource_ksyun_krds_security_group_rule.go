package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunKrdsSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceKsyunKrdsSecurityGroupRuleCreate,
		Read:     resourceKsyunKrdsSecurityGroupRuleRead,
		Delete:   resourceKsyunKrdsSecurityGroupRuleDelete,
		Importer: importKrdsSecurityGroupRule(),
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_rule_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunKrdsSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating krds sg rule , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRuleRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return removeKrdsSecurityGroupRule(d, meta)
}

func resourceKsyunKrdsSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readAndSetKrdsSecurityGroupRule(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading krds sg rule , error is %e", err)
	}
	return err
}
