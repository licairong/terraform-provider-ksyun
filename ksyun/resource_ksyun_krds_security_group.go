package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceKsyunKrdsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKrdsSecurityGroupCreate,
		Update: resourceKsyunKrdsSecurityGroupUpdate,
		Read:   resourceKsyunKrdsSecurityGroupRead,
		Delete: resourceKsyunKrdsSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf_krds_security_group",
			},
			"security_group_description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tf_krds_security_group_desc",
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_rule": {
				Type:     schema.TypeSet,
				Set:      secParameterToHash,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_group_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_rule_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group_rule_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func secParameterToHash(ruleMap interface{}) int {
	rule := ruleMap.(map[string]interface{})
	return hashcode.String(rule["security_group_rule_protocol"].(string) + "|" + rule["security_group_rule_name"].(string))
}

func resourceKsyunKrdsSecurityGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	err = createKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on creating sg , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	err = modifyKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on updating sg , error is %e", err)
	}
	return resourceKsyunKrdsSecurityGroupRead(d, meta)
}

func resourceKsyunKrdsSecurityGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	err = readKrdsAndSetSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on reading sg , error is %e", err)
	}
	return err
}

func resourceKsyunKrdsSecurityGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	err = removeKrdsSecurityGroup(d, meta)
	if err != nil {
		return fmt.Errorf("error on deleting sg , error is %e", err)
	}
	return err
}
