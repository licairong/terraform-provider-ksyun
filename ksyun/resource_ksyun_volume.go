package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKsyunVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunVolumeCreate,
		Update: resourceKsyunVolumeUpdate,
		Read:   resourceKsyunVolumeRead,
		Delete: resourceKsyunVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"volume_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					//"SSD2.0",
					//"SSD3.0",
					//"EHDD",
					//"SATA2.0",
					"SSD3.0",
					"EHDD",
					"ESSD_PL1",
					"ESSD_PL2",
					"ESSD_PL3",
				}, false),
				Default: "SSD3.0",
			},
			"volume_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(10, 32000),
				Default:      10,
			},

			"online_resize": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: volumeDiffSuppressFunc,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"charge_type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HourlyInstantSettlement",
					"Daily",
				}, false),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"volume_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// 快照建盘（API不返回这个值，所以diff时忽略这个值）
			"snapshot_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: kecDiskSnapshotIdDiffSuppress,
			},
		},
	}
}

func resourceKsyunVolumeCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.CreateVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on creating volume %q, %s", d.Id(), err)
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func resourceKsyunVolumeRead(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ReadAndSetVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on reading volume %q, %s", d.Id(), err)
	}
	return err
}

func resourceKsyunVolumeUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.ModifyVolume(d, resourceKsyunVolume())
	if err != nil {
		return fmt.Errorf("error on updating volume %q, %s", d.Id(), err)
	}
	err = resourceKsyunVolumeRead(d, meta)
	return err
}

func resourceKsyunVolumeDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ebsService := EbsService{meta.(*KsyunClient)}
	err = ebsService.RemoveVolume(d)
	if err != nil {
		return fmt.Errorf("error on deleting volume %q, %s", d.Id(), err)
	}
	return err
}
