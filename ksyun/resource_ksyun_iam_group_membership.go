package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceIamGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamGroupMembershipCreate,
		Read:   resourceIamGroupMembershipRead,
		Update: resourceIamGroupMembershipUpdate,
		Delete: resourceIamGroupMembershipDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIamGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	member := make(map[string]interface{})
	member["GroupName"]  = d.Get("group_name").(string)
	member["UserName"]  = d.Get("username").(string)

	_, err := c.AddUserToGroup(&member)
	if err != nil {
		return err
	}

	id := d.Get("group_name").(string) + "+" + d.Get("username").(string)
	d.SetId(id)

	_ = resourceIamGroupMembershipRead(d, m)

	return nil
}

func resourceIamGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]  = d.Get("username").(string)

	resp, err := c.ListGroupsForUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "UserNoSuchEntity") {
			d.SetId("")
			return nil
		}
		return err
	}

	groupList := (*resp)["ListGroupsForUserResult"].(map[string]interface{})["Groups"].(map[string]interface{})["member"].([]interface{})
	for _, v := range groupList {
		if v.(map[string]interface{})["GroupName"].(string) == d.Get("group_name").(string) {
			return nil
		}
	}
	d.SetId("")

	return nil
}

func resourceIamGroupMembershipUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIamGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	member := make(map[string]interface{})
	member["GroupName"]  = d.Get("group_name").(string)
	member["UserName"]  = d.Get("username").(string)

	_, err := c.RemoveUserFromGroup(&member)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
