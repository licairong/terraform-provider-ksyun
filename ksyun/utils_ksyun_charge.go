package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func durationSchemaDiffSuppressFunc(field string, value interface{}) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		if v, ok := d.GetOk(field); ok && v == value {
			return false
		}
		return true
	}
}

func durationTime(s string, e string) int {
	end, _ := time.Parse("2006-01-02 15:04:05", e)
	start, _ := time.Parse("2006-01-02 15:04:05", s)
	year := end.Year() - start.Year()
	month := int(end.Month()) - int(start.Month())
	if month < 0 {
		month = month + 12
	}
	return year*12 + month
}

func chargeExtraForVpc(data map[string]interface{}) map[string]SdkResponseMapping {
	extra := map[string]SdkResponseMapping{
		"ServiceEndTime": {
			Field: "purchase_time",
			FieldRespFunc: func(i interface{}) interface{} {
				return durationTime(data["CreateTime"].(string), i.(string))
			},
		},
		"ChargeType": {
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
				case "PrePaidByMonth":
					return "Monthly"
				default:
					return charge
				}
			},
		},
	}
	return extra
}
