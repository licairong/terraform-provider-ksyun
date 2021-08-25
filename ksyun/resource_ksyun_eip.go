package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunEipCreate,
		Read:   resourceKsyunEipRead,
		Update: resourceKsyunEipUpdate,
		Delete: resourceKsyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"line_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PrePaidByMonth",
					"Monthly",
					"PostPaidByPeak",
					"Peak",
					"PostPaidByDay",
					"Daily",
					"PostPaidByTransfer",
					"TrafficMonthly",
					"DailyPaidByTransfer",
					"HourlySettlement",
					"PostPaidByHour",
					"HourlyInstantSettlement",
				}, false),
			},
			"purchase_time": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: purchaseTimeDiffSuppressFunc,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  0,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allocation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_gateway_id": {
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
func resourceKsyunEipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn
	r := resourceKsyunEip()

	var resp *map[string]interface{}
	var err error

	req, err := SdkRequestAutoMapping(d, r, false, nil, nil)
	if err != nil {
		return fmt.Errorf("error on creating Address, %s", err)
	}
	err = validatePurchaseTime(&req, "PurchaseTime", "ChargeType", []string{"PrePaidByMonth", "Monthly"})
	if err != nil {
		return fmt.Errorf("error on creating Address, %s", err)
	}

	action := "AllocateAddress"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = conn.AllocateAddress(&req)
	if err != nil {
		return fmt.Errorf("error on creating Address, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["AllocationId"].(string))
	}
	return resourceKsyunEipRead(d, meta)
}

func resourceKsyunEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn

	req := make(map[string]interface{})
	req["AllocationId.1"] = d.Id()
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
		SdkResponseAutoResourceData(d, resourceKsyunEip(), items[0], nil)
	}
	return nil
}

func resourceKsyunEipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn
	r := resourceKsyunEip()

	var err error

	var only map[string]SdkReqTransform

	only = map[string]SdkReqTransform{
		"band_width": {},
		"project_id": {},
	}

	req, err := SdkRequestAutoMapping(d, r, true, only, nil)
	logger.Debug(logger.ReqFormat, "%%%%%%%%%%%%%%%%", &req)
	if err != nil {
		return fmt.Errorf("error on modifying Address, %s", err)
	}
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return fmt.Errorf("error on updating Eip, %s", err)
	}
	if len(req) > 0 {
		req["AllocationId"] = d.Id()
		action := "ModifyAddress"
		logger.Debug(logger.ReqFormat, action, req)
		_, err = conn.ModifyAddress(&req)
		if err != nil {
			return fmt.Errorf("error on modifying Address, %s", err)
		}
	}
	return resourceKsyunEipRead(d, meta)
}

func resourceKsyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.eipconn
	req := make(map[string]interface{})
	req["AllocationId"] = d.Id()
	action := "ReleaseAddress"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, req)
		resp, err1 := conn.ReleaseAddress(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting Address %q, %s", d.Id(), err1))
		}
	})

}
