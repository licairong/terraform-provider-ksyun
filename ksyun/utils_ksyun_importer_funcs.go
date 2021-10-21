package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"strings"
)

func importKecNetworkInterfaceAttachment(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("network_interface_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("instance_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importNatAssociate(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("nat_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("subnet_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importNetworkAclEntry(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("network_acl_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	num, err := strconv.Atoi(items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("rule_number", num)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	if items[2] != "in" && items[2] != "out" {
		return []*schema.ResourceData{d}, fmt.Errorf("direction must in or out")
	}
	err = d.Set("direction", items[2])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importLoadBalancerAclEntry(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("load_balancer_acl_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	num, err := strconv.Atoi(items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("rule_number", num)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("cidr_block", items[2])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importLoadBalancerAclAssociate(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("listener_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("load_balancer_acl_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importNetworkAclAssociate(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	err = d.Set("network_acl_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("subnet_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importSecurityGroupEntry(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 4 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}

	protocol := items[1]
	direction := items[2]
	cidrBlock := items[3]

	if protocol != "ip" {
		if len(items) != 6 {
			return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':' and size must  5")
		}
		if protocol == "icmp" {
			var (
				t int
				c int
			)
			t, err = strconv.Atoi(items[4])
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			c, err = strconv.Atoi(items[5])
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			err = d.Set("icmp_type", t)
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			err = d.Set("icmp_code", c)
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
		} else {
			var (
				from int
				to   int
			)
			from, err = strconv.Atoi(items[4])
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			to, err = strconv.Atoi(items[5])
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			err = d.Set("port_range_from", from)
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
			err = d.Set("port_range_to", to)
			if err != nil {
				return []*schema.ResourceData{d}, err
			}
		}
	}

	err = d.Set("security_group_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("protocol", protocol)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("direction", direction)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("cidr_block", cidrBlock)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func importAddressAssociate(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}
	err = d.Set("allocation_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("instance_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	if items[2] != "" {
		err = d.Set("network_interface_id", items[2])
		if err != nil {
			return []*schema.ResourceData{d}, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

func importBandWidthShareAssociate(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(d.Id(), ":")
	if len(items) < 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("import id must split with ':'")
	}
	err = d.Set("band_width_share_id", items[0])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	err = d.Set("allocation_id", items[1])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	return []*schema.ResourceData{d}, nil
}
