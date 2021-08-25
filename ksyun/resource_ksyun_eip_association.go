package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunEipAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipAssociationCreate,
		Read:   resourceKsyunEipAssociationRead,
		Delete: resourceKsyunEipAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Ipfwd",
					"Slb",
				}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"line_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"band_width_share_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_band_width_share": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceKsyunEipAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn
	r := resourceKsyunEipAssociation()

	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating AssociateAddress, %s", err)
	}

	action := "CreateScalingPolicy"
	logger.Debug(logger.ReqFormat, action, req)
	_, err = conn.AssociateAddress(&req)
	if err != nil {
		return fmt.Errorf("error on creating AssociateAddress, %s", err)
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("allocation_id"), d.Get("instance_id")))
	return resourceKsyunEipAssociationRead(d, meta)
}

func resourceKsyunEipAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn

	req := make(map[string]interface{})
	req["AllocationId.1"] = strings.Split(d.Id(), ":")[0]
	err := addProjectInfo(d, &req, client)
	if err != nil {
		return fmt.Errorf("error on reading Address %q, %s", d.Id(), err)
	}
	action := "DescribeAddresses"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeAddresses(&req)
	if err != nil {
		return fmt.Errorf("error on reading Address %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["AddressesSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunEipAssociation(), items[0], nil)
	}
	return nil
}

func resourceKsyunEipAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn
	req := make(map[string]interface{})
	req["AllocationId"] = strings.Split(d.Id(), ":")[0]
	action := "DisassociateAddress"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.DisassociateAddress(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  DisassociateAddress %q, %s", d.Id(), err1))
		}
	})
}
