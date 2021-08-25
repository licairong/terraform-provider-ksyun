package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunNatAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatAssociationCreate,
		Read:   resourceKsyunNatAssociationRead,
		Delete: resourceKsyunNatAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nat_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nat_ip_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"nat_ip_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"associate_nat_set": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunNatAssociationCreate(d *schema.ResourceData, m interface{}) error {
	vpcConn := m.(*KsyunClient).vpcconn

	req := make(map[string]interface{})
	creates := []string{
		"nat_id",
		"subnet_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "AssociateNat"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := vpcConn.AssociateNat(&req)
	logger.Debug(logger.AllFormat, action, req, resp, err)
	if err != nil {
		return fmt.Errorf("Error AssociateNat : %s ", err)
	}
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("Error AssociateNat ")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("Error AssociateNat:fail ")
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("nat_id"), d.Get("subnet_id")))
	return resourceKsyunNatAssociationRead(d, m)
}

func resourceKsyunNatAssociationRead(d *schema.ResourceData, m interface{}) error {
	vpcConn := m.(*KsyunClient).vpcconn
	p := strings.Split(d.Id(), ":")
	req := make(map[string]interface{})
	req["NatId.1"] = p[0]
	projectErr := addProjectInfo(d, &req, m.(*KsyunClient))
	if projectErr != nil {
		return projectErr
	}
	action := "DescribeNats"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := vpcConn.DescribeNats(&req)
	if err != nil {
		return fmt.Errorf("Error describeNats : %s ", err)
	}
	itemSet, ok := (*resp)["NatSet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items := itemSet.([]interface{})
	if len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetResourceDataByResp(d, items[0], natKeys)
	return nil
}

func resourceKsyunNatAssociationDelete(d *schema.ResourceData, m interface{}) error {
	vpcConn := m.(*KsyunClient).vpcconn
	deleteReq := make(map[string]interface{})
	p := strings.Split(d.Id(), ":")
	deleteReq["NatId"] = p[0]
	deleteReq["SubnetId"] = p[1]
	action := "DisassociateNat"
	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err1 := vpcConn.DisassociateNat(&deleteReq)
		logger.Debug(logger.AllFormat, action, deleteReq, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  disassociate nat %q, %s", d.Id(), err1))
		}
	})
}
