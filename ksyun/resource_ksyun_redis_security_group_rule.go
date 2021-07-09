package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// redis security group rule
// Deprecated: Use ksyun.resourceRedisSecurityGroup instead.
func resourceRedisSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityGroupRuleCreate,
		Delete: resourceRedisSecurityGroupRuleDelete,
		Update: resourceRedisSecurityGroupRuleUpdate,
		Read:   resourceRedisSecurityGroupRuleRead,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				var err error
				err = d.Set("security_group_id",d.Id())
				if err !=nil {
					return  nil,err
				}
				d.SetId(d.Id()+"-rules")
				return []*schema.ResourceData{d},err
			},
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceRedisSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		err  error
	)
	transform := map[string]SdkReqTransform{
		"rules": {
			mapping: "Cidrs",
			Type:    TransformWithN,
		},
	}
	err = processRedisSecurityGroupRule(d, meta, transform, false,d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	d.SetId(d.Get("security_group_id").(string)+"-rules")
	return resourceRedisSecurityGroupRuleRead(d,meta)
}

func resourceRedisSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
		del []interface{}
	)
	resp, err = readRedisSecurityGroup(d, meta,  d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	data := (*resp)["Data"].(map[string]interface{})

	//get rule id for del
	if rules,ok := data["rules"]; ok{
		for _,r := range rules.([]interface{}){
			rule := r.(map[string]interface{})
			del = append(del,rule["id"])
		}
	}

	transformDel := map[string]SdkReqTransform{
		"rules": {
			mapping: "SecurityGroupRuleId",
			Type:    TransformWithN,
			ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
				if len(del) > 0{
					return del,true
				}
				return nil,true
			},
		},
	}
	err = processRedisSecurityGroupRule(d, meta, transformDel, true,d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	return nil
}

func resourceRedisSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("rules") {
		var (
			resp *map[string]interface{}
			err  error
			oldArray []string
			newArray []string
			add []interface{}
			del []interface{}
		)
		resp, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
		if err != nil {
			return err
		}
		data := (*resp)["Data"].(map[string]interface{})
		rulesMap :=make(map[string]interface{})
		//get rule id for del
		if rules,ok := data["rules"]; ok{
			for _,r := range rules.([]interface{}){
				rule := r.(map[string]interface{})
				rulesMap[rule["cidr"].(string)] = rule["id"]
			}
		}
		o, n := d.GetChange("rules")
		for _, v := range o.(*schema.Set).List() {
			oldArray = append(oldArray, v.(string))
		}
		for _, v := range n.(*schema.Set).List() {
			newArray = append(newArray, v.(string))
		}
		for _, a := range oldArray {
			if _, ok := rulesMap[a];!ok{
				continue
			}
			exist := false
			for _, b := range newArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				del = append(del, rulesMap[a])
			}
		}
		for _, a := range newArray {
			exist := false
			for _, b := range oldArray {
				if a == b {
					exist = true
					break
				}
			}
			if !exist {
				add = append(add, a)
			}
		}
		transformAdd := map[string]SdkReqTransform{
			"rules": {
				mapping: "Cidrs",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(add) > 0{
						return add,true
					}
					return nil,true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformAdd, false,d.Get("security_group_id").(string))
		if err != nil {
			return err
		}
		transformDel := map[string]SdkReqTransform{
			"rules": {
				mapping: "SecurityGroupRuleId",
				Type:    TransformWithN,
				ValueFunc: func(data *schema.ResourceData) (interface{}, bool) {
					if len(del) > 0{
						return del,true
					}
					return nil,true
				},
			},
		}
		err = processRedisSecurityGroupRule(d, meta, transformDel, true,d.Get("security_group_id").(string))
		if err != nil {
			return err
		}

	}

	return resourceRedisSecurityGroupRuleRead(d,meta)
}

func resourceRedisSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err  error
	)
	resp, err = readRedisSecurityGroup(d, meta, d.Get("security_group_id").(string))
	if err != nil {
		return err
	}
	data := (*resp)["Data"].(map[string]interface{})
	extra := map[string]SdkResponseMapping{
		"rules": {
			Field: "rules",
			FieldRespFunc: func(i interface{}) interface{} {
				var cidr []string
				for _, v := range i.([]interface{}) {
					cidr = append(cidr, v.(map[string]interface{})["cidr"].(string))
				}
				return cidr
			},
		},
	}
	SdkResponseAutoResourceData(d, resourceRedisSecurityGroupRule(), data, extra)
	return nil
}


