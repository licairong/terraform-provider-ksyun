package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/rabbitmq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func checkRabbitmqAvailabilityZone(d *schema.ResourceData, meta interface{}, req map[string]interface{}) error {
	var (
		resp        *map[string]interface{}
		err         error
		data        interface{}
		regionCheck bool
		azCheck     bool
	)
	if v, ok := req["AvailabilityZone"]; ok {
		conn := meta.(*KsyunClient).rabbitmqconn
		currentRegion := *conn.Config.Region
		resp, err = conn.DescribeValidRegion(nil)
		if err != nil {
			return err
		}
		data, err = getSdkValue("Data.Regions", *resp)
		if err != nil {
			return err
		}
		for _, region := range data.([]interface{}) {
			if region.(map[string]interface{})["Code"].(string) != currentRegion {
				continue
			} else {
				regionCheck = true
				for _, az := range region.(map[string]interface{})["AvailabilityZones"].([]interface{}) {
					if az.(map[string]interface{})["Code"].(string) == v.(string) {
						azCheck = true
					}
				}
				break
			}
		}
		if !regionCheck {
			err = fmt.Errorf("region %s is not support", currentRegion)
		}

		if !azCheck {
			err = fmt.Errorf("availabilityZone %s is not support", v)
		}
	}

	return err
}

func checkRabbitmqPlugins(d *schema.ResourceData, meta interface{}, req map[string]interface{}) (error, bool) {
	var (
		resp        *map[string]interface{}
		err         error
		data        interface{}
		needRestart bool
	)
	if p, ok := req["EnablePlugins"]; ok {
		conn := meta.(*KsyunClient).rabbitmqconn
		resp, err = conn.SupportPlugins(nil)
		if err != nil {
			return err, needRestart
		}
		data, err = getSdkValue("Data", *resp)
		if err != nil {
			return err, needRestart
		}

		plugins := strings.Split(p.(string), ",")
		for _, plugin := range plugins {
			check := false
			for _, support := range data.([]interface{}) {
				if plugin == support.(map[string]interface{})["PluginName"].(string) {
					check = true
					//set ture only one time
					if !needRestart {
						needRestart = support.(map[string]interface{})["NeedToRestart"].(bool)
					}
					break
				}
			}
			if !check {
				return fmt.Errorf("plugin %s is not support", plugin), needRestart
			}
		}
	}
	return err, needRestart
}

func checkRabbitmqState(d *schema.ResourceData, meta interface{}, timeout time.Duration) error {
	var (
		err error
	)
	conn := meta.(*KsyunClient).rabbitmqconn
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"running"},
		Refresh:    rabbitmqStateRefreshFunc(conn, d.Id(), []string{"error"}),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
	}
	_, err = stateConf.WaitForState()
	return err
}

func rabbitmqStateRefreshFunc(conn *rabbitmq.Rabbitmq, instanceId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		queryReq := map[string]interface{}{"InstanceId": instanceId}
		logger.Debug(logger.ReqFormat, "DescribeRabbitmqInstance", queryReq)

		resp, err := conn.DescribeInstance(&queryReq)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, "DescribeRabbitmqInstance", queryReq, *resp)

		item, ok := (*resp)["Data"].(map[string]interface{})

		if !ok {
			return nil, "", fmt.Errorf("no instance information was queried. InstanceId:%s", instanceId)
		}
		status := item["Status"].(string)
		if status == "error" {
			return nil, "", fmt.Errorf("instance create error, status:%v", status)
		}

		for _, v := range failStates {
			if v == status {
				return nil, "", fmt.Errorf("instance create error, status:%v", status)
			}
		}
		return resp, status, nil
	}
}

func restartRabbitmqInstance(d *schema.ResourceData, meta interface{}) (err error) {
	if d.HasChange("force_restart") && d.Get("force_restart").(bool) {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		req := map[string]interface{}{
			"InstanceIds": d.Id(),
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		logger.Debug(logger.ReqFormat, "RestartInstance", req)
		_, err = conn.RestartInstance(&req)
	}
	return err
}

func modifyRabbitmqInstanceNameAndProject(d *schema.ResourceData, meta interface{}) (err error) {
	transform := map[string]SdkReqTransform{
		"instance_name": {},
		"project_id":    {},
	}

	req, err := SdkRequestAutoMapping(d, resourceKsyunRabbitmq(), true, transform, nil)
	if err != nil {
		return err
	}
	err = ModifyProjectInstance(d.Id(), &req, meta)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req["InstanceId"] = d.Id()
		logger.Debug(logger.ReqFormat, "RenameRabbitmqName", req)
		_, err = conn.Rename(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func modifyRabbitmqInstancePassword(d *schema.ResourceData, meta interface{}) (err error) {
	transform := map[string]SdkReqTransform{
		"instance_password": {},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunRabbitmq(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req["InstanceId"] = d.Id()
		logger.Debug(logger.ReqFormat, "ResetPassword", req)
		_, err = conn.ResetPassword(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func addRabbitmqRules(d *schema.ResourceData, meta interface{}, instanceId string, cidrs string) (err error) {
	if cidrs != "" {
		if instanceId == "" {
			instanceId = d.Id()
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req := map[string]interface{}{
			"InstanceId": instanceId,
			"Cidrs":      cidrs,
		}
		logger.Debug(logger.ReqFormat, "AddSecurityGroupRule", req)
		_, err = conn.AddSecurityGroupRule(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func validModifyRabbitmqInstanceRules(d *schema.ResourceData, r *schema.Resource, meta interface{}, instanceId string, isUpdate bool) (err error, add string, del string) {
	transform := map[string]SdkReqTransform{
		"cidrs": {mapping: "Cidrs"},
		"cidr":  {mapping: "Cidrs"},
	}
	req, err := SdkRequestAutoMapping(d, r, isUpdate, transform, nil)
	if err != nil {
		return err, add, add
	}
	if len(req) > 0 {
		var (
			rules []interface{}
			old   interface{}
		)
		if instanceId == "" && d.Id() == "" {
			add = req["Cidrs"].(string)
			return err, add, del
		}
		if d.HasChange("cidrs") {
			old, _ = d.GetChange("cidrs")
		}
		if d.HasChange("cidr") {
			old, _ = d.GetChange("cidr")
		}
		rules, err = readRabbitmqInstanceRules(d, meta, instanceId)
		if err != nil {
			return err, add, add
		}
		for _, rule := range rules {
			existCidr := rule.(map[string]interface{})["Cidr"].(string)
			exist := false
			for _, change := range strings.Split(req["Cidrs"].(string), ",") {
				if change == existCidr {
					exist = true
					break
				}
			}
			if !exist && strings.Contains(old.(string), existCidr) {
				del = del + existCidr + ","
			}
		}
		for _, change := range strings.Split(req["Cidrs"].(string), ",") {
			exist := false
			for _, rule := range rules {
				existCidr := rule.(map[string]interface{})["Cidr"].(string)
				if existCidr == change {
					exist = true
					break
				}
			}
			if !exist {
				add = add + change + ","
			}
		}
	}
	if add != "" {
		add = add[0 : len(add)-1]
	}
	if del != "" {
		del = del[0 : len(del)-1]
	}
	return err, add, del
}

func deleteRabbitmqRules(d *schema.ResourceData, meta interface{}, instanceId string, cidrs string) (resp *map[string]interface{}, err error) {
	if cidrs != "" {
		if instanceId == "" {
			instanceId = d.Id()
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req := map[string]interface{}{
			"InstanceId": instanceId,
			"Cidrs":      cidrs,
		}
		logger.Debug(logger.ReqFormat, "DeleteSecurityGroupRules", req)
		resp, err = conn.DeleteSecurityGroupRules(&req)
		if err != nil {
			return resp, err
		}
	}
	return resp, err
}

func readRabbitmqInstance(d *schema.ResourceData, meta interface{}, id string) (data map[string]interface{}, err error) {
	var (
		resp *map[string]interface{}
		ok   bool
	)
	conn := meta.(*KsyunClient).rabbitmqconn
	queryReq := make(map[string]interface{})
	if id == "" {
		id = d.Id()
	}
	queryReq["InstanceId"] = id
	action := "DescribeInstance"
	logger.Debug(logger.ReqFormat, action, queryReq)
	resp, err = conn.DescribeInstance(&queryReq)
	if err != nil {
		return data, fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, queryReq, *resp)
	if data, ok = (*resp)["Data"].(map[string]interface{}); !ok || data == nil || len(data) == 0 {
		return data, fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	return data, err
}

func readRabbitmqInstancePlugins(d *schema.ResourceData, meta interface{}, id string) (data []interface{}, err error) {
	var (
		resp  *map[string]interface{}
		value interface{}
		ok    bool
	)
	conn := meta.(*KsyunClient).rabbitmqconn
	queryReq := make(map[string]interface{})
	if id == "" {
		id = d.Id()
	}
	queryReq["InstanceId"] = id
	logger.Debug(logger.ReqFormat, "ListInstancePlugins", queryReq)
	resp, err = conn.ListInstancePlugins(&queryReq)
	if err != nil {
		return data, fmt.Errorf("error on reading instance %q, %s", d.Id(), err)
	}
	value, err = getSdkValue("Data.Plugins", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = value.([]interface{}); !ok {
		return data, fmt.Errorf("error on reading instance %q", d.Id())
	}
	return data, err
}

func readRabbitmqInstanceRules(d *schema.ResourceData, meta interface{}, id string) (data []interface{}, err error) {
	var (
		resp  *map[string]interface{}
		value interface{}
		ok    bool
	)
	conn := meta.(*KsyunClient).rabbitmqconn
	queryReq := make(map[string]interface{})
	if id == "" {
		id = d.Id()
	}
	queryReq["InstanceId"] = id
	logger.Debug(logger.ReqFormat, "DescribeSecurityGroupRules", queryReq)
	resp, err = conn.DescribeSecurityGroupRules(&queryReq)
	if err != nil {
		return data, fmt.Errorf("error on reading instance  %q security_roup_rules , %s", d.Id(), err)
	}
	value, err = getSdkValue("Data", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = value.([]interface{}); !ok {
		return data, fmt.Errorf("error on reading instance %q security_roup_rules ", d.Id())
	}
	return data, err
}

func allocateRabbitmqInstanceEip(d *schema.ResourceData, meta interface{}) (err error) {
	if d.HasChange("enable_eip") && d.Get("enable_eip").(bool) {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req := map[string]interface{}{
			"InstanceId": d.Id(),
		}
		logger.Debug(logger.ReqFormat, "AllocateEip", req)
		_, err = conn.AllocateEip(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func deallocateRabbitmqInstanceEip(d *schema.ResourceData, meta interface{}) (err error) {
	if d.HasChange("enable_eip") && !d.Get("enable_eip").(bool) {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req := map[string]interface{}{
			"InstanceId": d.Id(),
		}
		logger.Debug(logger.ReqFormat, "DeallocateEip", req)
		_, err = conn.DeallocateEip(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func modifyRabbitmqInstanceBandWidth(d *schema.ResourceData, meta interface{}) (err error) {
	transform := map[string]SdkReqTransform{
		"band_width": {},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunRabbitmq(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(req) > 0 {
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		conn := meta.(*KsyunClient).rabbitmqconn
		req["InstanceId"] = d.Id()
		logger.Debug(logger.ReqFormat, "ModifyBandWidth", req)
		_, err = conn.ModifyBandWidth(&req)
		if err != nil {
			return err
		}
	}
	return err
}

func validModifyRabbitmqInstancePlugins(d *schema.ResourceData, meta interface{}) (err error, enable string, disable string) {
	transform := map[string]SdkReqTransform{
		"enable_plugins": {},
	}
	req, err := SdkRequestAutoMapping(d, resourceKsyunRabbitmq(), true, transform, nil)
	if err != nil {
		return err, enable, disable
	}
	if len(req) > 0 {
		var (
			plugins        []interface{}
			enablePlugins  []string
			disablePlugins []string
			enableReq      []string
			disableReq     []string
			enableFlag     bool
			disableFlag    bool
			old            interface{}
		)
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err, enable, disable
		}
		if d.HasChange("enable_plugins") {
			old, _ = d.GetChange("enable_plugins")
		}
		plugins, err = readRabbitmqInstancePlugins(d, meta, "")
		if err != nil {
			return err, enable, disable
		}
		for _, plugin := range plugins {
			if int64(plugin.(map[string]interface{})["PluginStatus"].(float64)) == 1 {
				enablePlugins = append(enablePlugins, plugin.(map[string]interface{})["PluginName"].(string))
			} else {
				disablePlugins = append(disablePlugins, plugin.(map[string]interface{})["PluginName"].(string))
			}
		}
		for _, e := range enablePlugins {
			exist := false
			for _, change := range strings.Split(req["EnablePlugins"].(string), ",") {
				if change == e {
					exist = true
					break
				}
			}
			if !exist && strings.Contains(old.(string), e) {
				disableReq = append(disableReq, e)
			}
		}
		for _, e := range disablePlugins {
			for _, change := range strings.Split(req["EnablePlugins"].(string), ",") {
				if change == e {
					enableReq = append(enableReq, e)
					break
				}
			}
		}

		err, enableFlag, enable = checkEnableOrDisablePlugins(d, meta, enableReq)
		if err != nil {
			return err, enable, disable
		}
		err, disableFlag, disable = checkEnableOrDisablePlugins(d, meta, disableReq)
		if err != nil {
			return err, enable, disable
		}
		if ((enableFlag && enable != "") || (disableFlag && disable != "")) && !d.Get("force_restart").(bool) {
			err = fmt.Errorf("some plugin change need restart instance ,must set force_restart true")
			return err, enable, disable
		}

	}

	return err, enable, disable
}

func modifyRabbitmqInstancePlugins(d *schema.ResourceData, meta interface{}, enable string, disable string) (err error) {
	if enable != "" || disable != "" {
		conn := meta.(*KsyunClient).rabbitmqconn
		err = checkRabbitmqState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		if enable != "" {
			r := make(map[string]interface{})
			r["InstanceId"] = d.Id()
			r["EnablePlugins"] = enable
			logger.Debug(logger.ReqFormat, "EnableInstancePlugins", r)
			_, err = conn.EnableInstancePlugins(&r)
		}
		if disable != "" {
			r := make(map[string]interface{})
			r["InstanceId"] = d.Id()
			r["DisablePlugins"] = disable
			logger.Debug(logger.ReqFormat, "DisableInstancePlugins", r)
			_, err = conn.DisableInstancePlugins(&r)
		}
	}
	return err
}

func checkEnableOrDisablePlugins(d *schema.ResourceData, meta interface{}, plugins []string) (err error, restart bool, str string) {
	if len(plugins) > 0 {
		r := make(map[string]interface{})
		for i, p := range plugins {
			if i < len(plugins)-1 {
				str = str + p + ","
			} else {
				str = str + p
			}
		}
		r["EnablePlugins"] = str
		err, restart = checkRabbitmqPlugins(d, meta, r)
	}
	return err, restart, str
}
