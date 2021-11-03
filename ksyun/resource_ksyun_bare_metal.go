package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func resourceKsyunBareMetal() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunBareMetalCreate,
		Read:   resourceKsyunBareMetalRead,
		Update: resourceKsyunBareMetalUpdate,
		Delete: resourceKsyunBareMetalDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
		CustomizeDiff: bareMetalCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hyper_threading": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Open",
					"Close",
					"NoChange",
				}, false),
				Default: "NoChange",
			},
			"raid": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Raid0",
					"Raid1",
					"Raid5",
					"Raid10",
					"Raid50",
					"SRaid0",
				}, false),
				ConflictsWith: []string{"raid_id"},
			},
			"raid_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"raid"},
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_interface_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bond4",
					"single",
					"dual",
				}, false),
				Default: "bond4",
			},
			"bond_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bond0",
					"bond1",
				}, false),
				Default:          "bond1",
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 3,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"dns1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ksc_epc",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"security_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"classic",
					"no",
				}, false),
				Default: "no",
			},
			"cloud_monitor_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"classic",
					"no",
				}, false),
				Default: "no",
			},
			"extension_subnet_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"extension_private_ip_address": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"extension_security_group_ids": {
				Type:     schema.TypeSet,
				MaxItems: 3,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:              schema.HashString,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"extension_dns1": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"extension_dns2": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"system_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
				Default: "XFS",
			},
			"data_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
				Default: "XFS",
			},
			"data_disk_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"/DATA/disk",
					"/data",
				}, false),
				Default: "/DATA/disk",
			},
			"data_disk_catalogue_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSuffix",
					"NaturalNumber",
					"NaturalNumberFromZero",
				}, false),
				Default: "NaturalNumber",
			},
			"nvme_data_file_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"EXT4",
					"XFS",
				}, false),
			},
			"nvme_data_disk_catalogue": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"/DATA/disk",
					"/data",
				}, false),
			},
			"nvme_data_disk_catalogue_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NoSuffix",
					"NaturalNumber",
					"NaturalNumberFromZero",
				}, false),
			},
			"container_agent": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"supported",
					"unsupported",
				}, false),
				Default: "unsupported",
			},
			"computer_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"path": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"force_re_install": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: bareMetalDiffSuppressFunc,
			},
			"extension_network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKsyunBareMetalCreate(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.CreateBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on creating bare metal %q, %s", d.Id(), err)
	}
	return resourceKsyunBareMetalRead(d, meta)
}

func resourceKsyunBareMetalRead(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.ReadAndSetBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on reading bare metal %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunBareMetalUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.ModifyBareMetal(d, resourceKsyunBareMetal())
	if err != nil {
		return fmt.Errorf("error on updating bare metal %q, %s", d.Id(), err)
	}
	return resourceKsyunBareMetalRead(d, meta)
}

func resourceKsyunBareMetalDelete(d *schema.ResourceData, meta interface{}) (err error) {
	bareMetalService := BareMetalService{meta.(*KsyunClient)}
	err = bareMetalService.RemoveBareMetal(d)
	if err != nil {
		return fmt.Errorf("error on deleting bare metal %q, %s", d.Id(), err)
	}
	return err

}
