package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunNatCreate,
		Update: resourceKsyunNatUpdate,
		Read:   resourceKsyunNatRead,
		Delete: resourceKsyunNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nat_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateName,
				Computed:     true,
			},

			"nat_mode": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNatMode,
			},

			"nat_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateNatType,
			},

			"nat_ip_number": {
				Type:         schema.TypeInt,
				ForceNew:     false,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateNatIpNumber,
			},

			"band_width": {
				Type:         schema.TypeInt,
				ForceNew:     false,
				Required:     true,
				ValidateFunc: validateNatBandWidth,
			},

			"charge_type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "DailyPaidByTransfer",
				ValidateFunc: validateNatChargeType,
			},

			"purchase_time": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
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

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunNatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	var resp *map[string]interface{}
	var err error

	createNat, _ := SdkRequestAutoMapping(d, resourceKsyunNat(), false, nil, nil)
	err = validatePurchaseTime(&createNat, "purchase_time", "charge_type", []string{"Monthly"})
	if err != nil {
		return fmt.Errorf("error on creating nat, %s", err)
	}
	action := "CreateNat"
	logger.Debug(logger.ReqFormat, action, createNat)
	resp, err = conn.CreateNat(&createNat)
	if err != nil {
		return fmt.Errorf("error on creating nat, %s", err)
	}
	if resp != nil {
		d.SetId((*resp)["NatId"].(string))
	}
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	readNat := make(map[string]interface{})
	readNat["NatId.1"] = d.Id()
	err := addProjectInfo(d, &readNat, meta.(*KsyunClient))
	if err != nil {
		return fmt.Errorf("error on reading nat %q, %s", d.Id(), err)
	}
	action := "DescribeNats"
	logger.Debug(logger.ReqFormat, action, readNat)
	resp, err := conn.DescribeNats(&readNat)
	if err != nil {
		return fmt.Errorf("error on reading nat %q, %s", d.Id(), err)
	}
	if resp != nil {
		items, ok := (*resp)["NatSet"].([]interface{})
		if !ok || len(items) == 0 {
			d.SetId("")
			return nil
		}
		SdkResponseAutoResourceData(d, resourceKsyunNat(), items[0], resourceKsyunNatExtra())
	}
	return nil
}

func resourceKsyunNatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn

	req, err := SdkRequestAutoMapping(d, resourceKsyunNat(), true, nil, nil)
	if err != nil {
		return fmt.Errorf("error on updating Nat, %s", err)
	}
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return fmt.Errorf("error on updating Nat, %s", err)
	}
	if len(req) > 0 {
		req["NatId"] = d.Id()
		action := "ModifyNat"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ModifyNat(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)
		if err != nil {
			return fmt.Errorf("error on updating Nat, %s", err)
		}
	}
	return resourceKsyunNatRead(d, meta)
}

func resourceKsyunNatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.vpcconn
	deleteNat := make(map[string]interface{})
	deleteNat["NatId"] = d.Id()
	action := "DeleteNat"

	return resource.Retry(25*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteNat)
		resp, err1 := conn.DeleteNat(&deleteNat)
		logger.Debug(logger.AllFormat, action, deleteNat, resp, err1)
		if err1 == nil {
			return nil
		} else if notFoundError(err1) {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("error on  deleting nat %q, %s", d.Id(), err1))
		}
	})

}

func resourceKsyunNatExtra() map[string]SdkResponseMapping {
	extra := make(map[string]SdkResponseMapping)
	extra["ChargeType"] = SdkResponseMapping{
		Field: "charge_type",
		FieldRespFunc: func(i interface{}) interface{} {
			charge := i.(string)
			switch charge {
			case "PostPaidByPeak":
				return "Peak"
			case "PostPaidByDay":
				return "Daily"
			case "PostPaidByTransfer":
				return "TrafficMonthly"
			default:
				return charge
			}
		},
	}
	return extra
}
