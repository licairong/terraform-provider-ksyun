package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceIamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamUserCreate,
		Read:   resourceIamUserRead,
		Update: resourceIamUserUpdate,
		Delete: resourceIamUserDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type: schema.TypeString,
				Optional: true,
				Default: "/",
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIamUserCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]  = d.Get("username").(string)
	user["RealName"]  = d.Get("realname").(string)
	user["Remark"]    = d.Get("remark").(string)
	user["Path"]      = d.Get("path").(string)
	user["Phone"]     = d.Get("phone").(string)
	email := d.Get("email").(string)
	if email != "" {
		user["Email"]     = d.Get("email").(string)
	}

	resp, err := c.CreateUser(&user)
	if err != nil {
		return err
	}

	userinfo := (*resp)["CreateUserResult"].(map[string]interface{})["User"].(map[string]interface{})
	d.SetId(userinfo["UserName"].(string))

	_ = resourceIamUserRead(d, m)

	return nil
}

func resourceIamUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]  = d.Id()

	resp, err := c.GetUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "UserNoSuchEntity") {
			d.SetId("")
			return nil
		}
		return err
	}

	userinfo := (*resp)["GetUserResult"].(map[string]interface{})["User"].(map[string]interface{})
	d.Set("username", userinfo["UserName"].(string))
	d.Set("path", userinfo["Path"].(string))
	d.Set("create_date", userinfo["CreateDate"].(string))
	v, ok := userinfo["RealName"].(string); if ok {
		d.Set("realname", v)
	}
	v, ok = userinfo["Remark"].(string); if ok {
		d.Set("remark", v)
	}
	v, ok = userinfo["Phone"].(string); if ok {
		d.Set("phone", v)
	}
	v, ok = userinfo["Email"].(string); if ok {
		d.Set("email", v)
	}

	return nil
}

func resourceIamUserUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]    = d.Id()
	user["NewUserName"]    = d.Get("username").(string)
	if d.Get("realname").(string) != "" {
		user["NewRealName"] = d.Get("realname").(string)
	}
	if d.Get("remark").(string) != "" {
		user["NewRemark"]   = d.Get("remark").(string)
	}
	if d.Get("path").(string) != "" {
		user["NewPath"]     = d.Get("path").(string)
	}
	if d.Get("phone").(string) != "" {
		user["NewPhone"]    = d.Get("phone").(string)
	}
	if d.Get("email").(string) != "" {
		user["NewEmail"] = d.Get("email").(string)
	}

	_, err := c.UpdateUser(&user)
	if err != nil {
		return err
	}
	d.SetId(d.Get("username").(string))

	_ = resourceIamUserRead(d, m)

	return nil
}

func resourceIamUserDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*KsyunClient).iamconn

	r, err := c.GetAccountAllProjectList(nil)
	if err != nil {
		return err
	}

	projectList := (*r)["ListProjectResult"].(map[string]interface{})["ProjectList"].([]interface{})
	for _, p := range projectList {
		member := make(map[string]interface{})
		member["ProjectId"]  = p.(map[string]interface{})["ProjectId"].(float64)

		r, err = c.ListProjectMember(&member)
		if err != nil {
			return err
		}

		memberList := (*r)["ListProjectMember"].([]interface{})

		for _, v := range memberList {
			if v.(map[string]interface{})["IdentityName"].(string) == d.Get("username").(string) {
				member["MemberIds"] = v.(map[string]interface{})["MemberId"].(float64)
				_, err = c.DeleteProjectMember(&member)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	user := make(map[string]interface{})
	user["UserName"]  = d.Id()

	_, err = c.DeleteUser(&user)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
