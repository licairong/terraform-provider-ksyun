package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceIamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamAccessKeyCreate,
		Read:   resourceIamAccessKeyRead,
		Update: resourceIamAccessKeyUpdate,
		Delete: resourceIamAccessKeyDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_access_key": {
				Type:     schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIamAccessKeyCreate(d *schema.ResourceData, m interface{}) (err error) {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]  = d.Get("username").(string)

	resp, err := c.CreateAccessKey(&user)
	if err != nil {
		return err
	}

	userinfo := (*resp)["CreateAccessKeyResult"].(map[string]interface{})["AccessKey"].(map[string]interface{})
	d.SetId(userinfo["AccessKeyId"].(string))
	d.Set("secret_access_key", userinfo["SecretAccessKey"].(string))

	_ = resourceIamAccessKeyRead(d, m)

	return nil
}

func resourceIamAccessKeyRead(d *schema.ResourceData, m interface{}) (err error) {
	c := m.(*KsyunClient).iamconn

	user := make(map[string]interface{})
	user["UserName"]  = d.Get("username").(string)

	resp, err := c.ListAccessKeys(&user)
	if err != nil {
		if strings.Contains(err.Error(), "UserNoSuchEntity") {
			d.SetId("")
			return nil
		}
		return err
	}

	aks := (*resp)["ListAccessKeyResult"].(map[string]interface{})["AccessKeyMetadata"].(map[string]interface{})["member"].([]interface{})
	for _, ak := range aks {
		if ak.(map[string]interface{})["AccessKeyId"].(string) == d.Id() {
			d.Set("username", ak.(map[string]interface{})["UserName"].(string))
			d.Set("access_key_id", ak.(map[string]interface{})["AccessKeyId"].(string))
			d.Set("status", ak.(map[string]interface{})["Status"].(string))
			d.Set("create_date", ak.(map[string]interface{})["CreateDate"].(string))
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceIamAccessKeyUpdate(d *schema.ResourceData, m interface{}) (err error) {
	return nil
}

func resourceIamAccessKeyDelete(d *schema.ResourceData, m interface{}) (err error) {
	c := m.(*KsyunClient).iamconn

	ak := make(map[string]interface{})
	ak["UserName"]  = d.Get("username").(string)
	ak["AccessKeyId"]  = d.Id()

	_, err = c.DeleteAccessKey(&ak)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
