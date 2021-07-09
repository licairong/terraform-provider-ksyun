package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strconv"
	"time"
)

func resourceKsyunKrds() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunKrdsCreate,
		Update: resourceKsyunMysqlUpdate,
		Read:   resourceKsyunMysqlRead,
		Delete: resourceKsyunMysqlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_identifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "source instance identifier",
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,3}|db.disk.d{1,5} , " +
					"db.ram is rds random access memory size, db.disk is disk size",
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HRDS",
				ValidateFunc: validation.StringInSlice([]string{
					"HRDS",
					"TRDS",
					"ERDS",
					"SINGLERDS",
				}, false),
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "engine is db type, only support mysql|percona",
				ForceNew:    true,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db engine version only support 5.5|5.6|5.7|8.0",
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"master_user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bill_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DAY",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DAY",
				}, false),
			},
			"duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "proprietary security group id for krds",
			},
			"db_parameter_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_backup_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone_2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"instance_create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_has_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"force_restart": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func parameterToHash(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	return hashcode.String(m["name"].(string) + "|" + m["value"].(string))
}

func resourceKsyunKrdsCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var resp *map[string]interface{}
	//create a temp parameterGroup if need
	err = createTempParameterGroup(d, meta)
	if err != nil {
		return err
	}
	r := resourceKsyunKrds()
	transform := map[string]SdkReqTransform{
		"db_instance_class":     {mapping: "DBInstanceClass"},
		"db_instance_name":      {mapping: "DBInstanceName"},
		"db_instance_type":      {mapping: "DBInstanceType"},
		"db_parameter_group_id": {mapping: "DBParameterGroupId"},
		"instance_has_eip":      {Ignore: true},
		"parameters":            {Ignore: true},
		"force_restart":         {Ignore: true},
		"availability_zone_1":   {mapping: "AvailabilityZone.1"},
		"availability_zone_2":   {mapping: "AvailabilityZone.2"},
	}

	createReq, err := SdkRequestAutoMapping(d, r, false, transform, nil, SdkReqParameter{
		onlyTransform: false,
	})

	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}
	action := "CreateDBInstance"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateDBInstance(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		krdsInstance := bodyData["DBInstance"].(map[string]interface{})
		instanceId := krdsInstance["DBInstanceIdentifier"].(string)
		d.SetId(instanceId)
	}
	pending := []string{tCreatingStatus}
	target := []string{tActiveStatus, tFailedStatus, tDeletedStatus, tStopedStatus}
	err = checkStatus(d, conn, pending, target, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	//clean temp parameterGroup if need
	err = resourceKsyunMysqlDeleteParameterGroup(d, meta)
	if err != nil {
		return err
	}
	if d.Get("instance_has_eip") == true {
		err = allocateOrReleaseInstanceEip(d, meta)
		if err != nil {
			return fmt.Errorf("error on allocate Instance(krds) eip: %s", err)
		}
	}
	return resourceKsyunMysqlRead(d, meta)
}

func processParameters(d *schema.ResourceData, meta interface{}, paramsReq *map[string]interface{}, isUpdate bool) (update bool, restart bool, err error) {
	// defer to check force restart
	defer func(d *schema.ResourceData, meta interface{}) {
		if isUpdate && !d.Get("force_restart").(bool) && restart && err == nil {
			err = fmt.Errorf("update parameters must restart,must set force_restart true ")
		} else if isUpdate && d.Get("force_restart").(bool) {
			restart = d.Get("force_restart").(bool)
		}
	}(d, meta)

	oldP, newP := d.GetChange("parameters")
	documentOld := oldP.(*schema.Set).List()
	documentNew := newP.(*schema.Set).List()
	var toDefault []string
	mapOld := make(map[string]interface{})
	mapNew := make(map[string]interface{})
	for _, i := range documentOld {
		name := i.(map[string]interface{})["name"].(string)
		value := i.(map[string]interface{})["value"].(string)
		mapOld[name] = value
	}
	for _, i := range documentNew {
		name := i.(map[string]interface{})["name"].(string)
		value := i.(map[string]interface{})["value"].(string)
		mapNew[name] = value
	}
	for k, _ := range mapOld {
		if _, ok := mapNew[k]; !ok {
			toDefault = append(toDefault, k)
		}
	}

	var resp *map[string]interface{}
	update = false
	restart = false

	if len(documentNew) > 0 || len(toDefault) > 0 {
		conn := meta.(*KsyunClient).krdsconn
		currentParameter := make(map[string]interface{})
		// get current db_parameter_group when exist
		if paramId, ok := d.GetOk("db_parameter_group_id"); ok && paramId != nil && paramId != "" {
			queryParam := map[string]interface{}{
				"DBParameterGroupId": paramId,
			}
			resp, err = conn.DescribeDBParameterGroup(&queryParam)
			if err != nil {
				err = fmt.Errorf("error on check parameters: error is %v", err)
				return update, restart, err
			}
			obj, err := getSdkValue("Data.DBParameterGroups.0.Parameters", *resp)
			if err != nil {
				err = fmt.Errorf("error on check parameters: error is %v", err)
				return update, restart, err
			}
			currentParameter = obj.(map[string]interface{})
		}
		defaultParam := map[string]interface{}{
			"Engine":        d.Get("engine"),
			"EngineVersion": d.Get("engine_version"),
		}
		resp, err = conn.DescribeEngineDefaultParameters(&defaultParam)
		if err == nil {
			obj, _ := getSdkValue("Data.Parameters", *resp)
			parameters := obj.(map[string]interface{})
			num := 0
			for _, i := range documentNew {
				name := i.(map[string]interface{})["name"].(string)
				if name == "" {
					continue
				}
				value := i.(map[string]interface{})["value"].(string)
				if _, ok := parameters[name]; !ok {
					err = fmt.Errorf("error on check parameters: parameter not support %s", name)
					return update, restart, err
				}
				if parameters[name].(map[string]interface{})["Type"] == "string" {
					inValid := true
					enum := parameters[name].(map[string]interface{})["Enums"].([]interface{})
					for _, e := range enum {
						if e.(string) == value {
							inValid = false
							break
						}
					}
					if inValid {
						err = fmt.Errorf("error on check parameters: parameter  %s value must in %v", name, enum)
						return update, restart, err
					}
				} else if parameters[name].(map[string]interface{})["Type"] == "integer" {
					valueNum, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						err = fmt.Errorf("error on check parameters: parameter  %s value must integer", name)
						return update, restart, err
					}
					max := int64(parameters[name].(map[string]interface{})["Max"].(float64))
					min := int64(parameters[name].(map[string]interface{})["Min"].(float64))
					if valueNum < min || valueNum > max {
						err = fmt.Errorf("error on check parameters: parameter  %s value must in (%v,%v)", name, min, max)
						return update, restart, err
					}
				}
				if v, ok := currentParameter[name]; ok {
					var currentValue string
					if f, ok := v.(float64); ok {
						currentValue = strconv.FormatInt(int64(f), 10)
					} else if s, ok := v.(string); ok {
						currentValue = s
					}
					if currentValue != value {
						num, restart = generateParameters(paramsReq, num, name, value, parameters)
					}
				} else if _, ok := currentParameter[name]; !ok {
					num, restart = generateParameters(paramsReq, num, name, value, parameters)
				}
			}
			for _, name := range toDefault {
				if _, ok := parameters[name]; ok {
					var defaultValue string
					if parameters[name].(map[string]interface{})["Type"] == "integer" {
						defaultValue = strconv.FormatInt(int64(parameters[name].(map[string]interface{})["Default"].(float64)), 10)
					} else {
						defaultValue = parameters[name].(map[string]interface{})["Default"].(string)
					}
					if v, ok := currentParameter[name]; ok {
						var currentValue string
						if f, ok := v.(float64); ok {
							currentValue = strconv.FormatInt(int64(f), 10)
						} else if s, ok := v.(string); ok {
							currentValue = s
						}
						if currentValue != defaultValue {
							num, restart = generateParameters(paramsReq, num, name, defaultValue, parameters)
						}
					} else if _, ok := currentParameter[name]; !ok {
						num, restart = generateParameters(paramsReq, num, name, defaultValue, parameters)
					}
				}
			}
		}
		if len(*paramsReq) > 0 {
			update = true
		}
		return update, restart, err
	}
	return update, restart, err
}

func generateParameters(paramsReq *map[string]interface{}, num int, name string, value string, parameters map[string]interface{}) (int, bool) {
	num = num + 1
	restart := false
	(*paramsReq)["Parameters.Name."+strconv.Itoa(num)] = name
	(*paramsReq)["Parameters.Value."+strconv.Itoa(num)] = value
	//need restart
	if parameters[name].(map[string]interface{})["RestartRequired"].(bool) {
		restart = true
	}
	return num, restart
}

func createTempParameterGroup(d *schema.ResourceData, meta interface{}) error {
	//create new parameter
	conn := meta.(*KsyunClient).krdsconn
	paramsReq := make(map[string]interface{})
	paramsReq["DBParameterGroupName"] = d.Get("db_instance_name").(string) + "_param"
	paramsReq["Description"] = d.Get("db_instance_name").(string) + "_desc"
	paramsReq["Engine"] = d.Get("engine")
	paramsReq["EngineVersion"] = d.Get("engine_version")
	create, _, err := processParameters(d, meta, &paramsReq, false)
	if err != nil {
		return err
	}
	if create {
		action := "CreateDBParameterGroup"
		logger.Debug(logger.RespFormat, action, paramsReq)
		paramResp, err := conn.CreateDBParameterGroup(&paramsReq)
		logger.Debug(logger.AllFormat, action, paramsReq, paramResp, err)
		if err != nil {
			return fmt.Errorf("error on create Instance(krds) DBParameterGroup : %s", err)
		}
		parameterId, err := getSdkValue("Data.DBParameterGroup.DBParameterGroupId", *paramResp)
		if err != nil {
			return fmt.Errorf("error on create Instance(krds) DBParameterGroup : %s", err)
		}
		return d.Set("db_parameter_group_id", parameterId)
	}
	return nil
}

func modifyParameters(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	paramsReq := make(map[string]interface{})
	modify, restart, err := processParameters(d, meta, &paramsReq, true)
	if err != nil {
		return err
	}
	if modify {
		paramsReq["DBParameterGroupId"] = d.Get("db_parameter_group_id").(string)
		mdAction := "ModifyDBParameterGroup"
		logger.Debug(logger.RespFormat, mdAction, paramsReq)
		paramResp, err := conn.ModifyDBParameterGroup(&paramsReq)
		logger.Debug(logger.AllFormat, mdAction, paramsReq, paramResp, err)
		if err != nil {
			return err
		}

	}
	if restart {
		err := checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		restartParam := make(map[string]interface{})
		restartParam["DBInstanceIdentifier"] = d.Id()
		mdAction := "RebootDBInstance"
		logger.Debug(logger.RespFormat, mdAction, restartParam)
		_, err = conn.RebootDBInstance(&restartParam)
		logger.Debug(logger.AllFormat, mdAction, restartParam, err)
		if err != nil {
			return err
		}
	}
	return nil
}

func allocateOrReleaseInstanceEip(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	err := checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}
	req := map[string]interface{}{
		"DBInstanceIdentifier": d.Id(),
	}

	if d.Get("instance_has_eip") == true {
		action := "AllocateDBInstanceEip"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.AllocateDBInstanceEip(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)

		if err != nil {
			return err
		}
		return nil
	} else {
		action := "ReleaseDBInstanceEip"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.ReleaseDBInstanceEip(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)

		if err != nil {
			return err
		}
		return nil
	}
}

func mysqlInstanceStateRefresh(client *krds.Krds, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"DBInstanceIdentifier": instanceId}
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeDBInstances(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)
		if err != nil {
			return nil, "", err
		}
		bodyData := (*resp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		krdsInstance := instances[0].(map[string]interface{})
		state := krdsInstance["DBInstanceStatus"].(string)

		return resp, state, nil
	}
}

func resourceKsyunMysqlReadParameter(dbParameterGroupId string, d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	req := map[string]interface{}{"DBParameterGroupId": dbParameterGroupId}
	action := "DescribeDBParameterGroup"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeDBParameterGroup(&req)
	logger.Debug(logger.AllFormat, action, req, resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(krds) %q, %s", d.Id(), err)
	}
	parameter, _ := getSdkValue("Data.DBParameterGroups.0.Parameters", *resp)
	var parameters []map[string]interface{}
	remote := make(map[string]map[string]interface{})
	if local, ok := d.GetOk("parameters"); ok {
		if parameter != nil {
			for k, v := range parameter.(map[string]interface{}) {
				m := make(map[string]interface{})
				m["name"] = k
				if vf, ok := v.(float64); ok {
					m["value"] = fmt.Sprintf("%v", strconv.FormatInt(int64(vf), 10))
				} else {
					m["value"] = fmt.Sprintf("%v", v)
				}

				remote[k] = m
			}
		}
		for _, value := range local.(*schema.Set).List() {
			name := value.(map[string]interface{})["name"]
			for k, v := range remote {
				if k == name {
					parameters = append(parameters, v)
					break
				}
			}
		}
	}
	return d.Set("parameters", parameters)

}

func resourceKsyunMysqlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.krdsconn
	req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := conn.DescribeDBInstances(&req)
	logger.Debug(logger.AllFormat, action, req, resp, err)

	if err != nil {
		return fmt.Errorf("error on reading Instance(krds) %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(krds) body %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})

	extra := make(map[string]SdkResponseMapping)
	for _, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		krdsMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
				extra["DBInstanceClass"] = SdkResponseMapping{
					Field: "db_instance_class",
					FieldRespFunc: func(i interface{}) interface{} {
						value := i.(map[string]interface{})
						return fmt.Sprintf("db.ram.%v|db.disk.%v", value["Ram"], value["Disk"])
					},
				}
				krdsMap["DBInstanceClass"] = v
			} else if k == "ReadReplicaDBInstanceIdentifiers" || k == "DBSource" {
				continue
			} else if k == "Eip" {
				krdsMap["instance_has_eip"] = true
			} else if k == "MasterAvailabilityZone" {
				krdsMap["availability_zone_1"] = v
			} else if k == "SlaveAvailabilityZone" {
				krdsMap["availability_zone_2"] = v
			} else {
				krdsMap[Camel2Hungarian(k)] = v
			}
		}
		if krdsMap["instance_has_eip"] == nil {
			krdsMap["instance_has_eip"] = false
		}
		SdkResponseAutoResourceData(d, resourceKsyunKrds(), krdsMap, extra)
		if d.Get("force_restart") != nil {
			_ = d.Set("force_restart", d.Get("force_restart"))
		} else {
			_ = d.Set("force_restart", false)
		}
		return resourceKsyunMysqlReadParameter(d.Get("db_parameter_group_id").(string), d, meta)
	}
	return nil
}

func modifyDBInstance(d *schema.ResourceData, meta interface{}, oldType interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var modifyInstanceParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"db_instance_name":      {mapping: "DBInstanceName"},
		"master_user_password":  {},
		"security_group_id":     {},
		"preferred_backup_time": {},
		"project_id":            {},
	}
	modifyInstanceParam, err = SdkRequestAutoMapping(d, resourceKsyunKrds(), true, transform, nil)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}
	//modify project
	err = ModifyProjectInstance(d.Id(), &modifyInstanceParam, meta)
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %s", err)
	}
	if len(modifyInstanceParam) > 0 {
		if _, ok := modifyInstanceParam["PreferredBackupTime"]; ok && oldType == "RR" {
			return fmt.Errorf("error on updating instance , krds rr is not support update %s", "preferred_backup_time")
		}
		modifyInstanceParam["DBInstanceIdentifier"] = d.Id()
		action := "ModifyDBInstance"
		logger.Debug(logger.ReqFormat, action, modifyInstanceParam)
		_, err = conn.ModifyDBInstance(&modifyInstanceParam)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	return err
}

func modifyDBInstanceType(d *schema.ResourceData, meta interface{}, oldType interface{}, newType interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var modifyDBInstanceTypeParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"db_instance_type": {mapping: "DBInstanceType"},
	}
	modifyDBInstanceTypeParam, err = SdkRequestAutoMapping(d, resourceKsyunKrds(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(modifyDBInstanceTypeParam) > 0 {
		if oldType != "TRDS" {
			return fmt.Errorf("error on updating instance , krds is not support %s to %s", oldType, newType)
		}
		if oldType == "RR" {
			return fmt.Errorf("error on updating instance , krds rr is not support update %s", "db_instance_type")
		}
		err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		modifyDBInstanceTypeParam["DBInstanceIdentifier"] = d.Id()
		action := "ModifyDBInstanceType"
		logger.Debug(logger.ReqFormat, action, modifyDBInstanceTypeParam)
		_, err = conn.ModifyDBInstanceType(&modifyDBInstanceTypeParam)
		logger.Debug(logger.AllFormat, action, modifyDBInstanceTypeParam, err)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	return err
}

func modifyDBInstanceSpec(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var modifyDBInstanceSpecParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"db_instance_class": {mapping: "DBInstanceClass"},
	}
	modifyDBInstanceSpecParam, err = SdkRequestAutoMapping(d, resourceKsyunKrds(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(modifyDBInstanceSpecParam) > 0 {
		err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		modifyDBInstanceSpecParam["DBInstanceIdentifier"] = d.Id()
		action := "ModifyDBInstanceSpec"
		logger.Debug(logger.ReqFormat, action, modifyDBInstanceSpecParam)
		_, err = conn.ModifyDBInstanceSpec(&modifyDBInstanceSpecParam)
		logger.Debug(logger.AllFormat, action, modifyDBInstanceSpecParam, err)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	return err
}

func upgradeDBInstanceEngineVersion(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var upgradeDBInstanceEngineVersionParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"engine":         {},
		"engine_version": {},
	}
	upgradeDBInstanceEngineVersionParam, err = SdkRequestAutoMapping(d, resourceKsyunKrds(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(upgradeDBInstanceEngineVersionParam) > 0 {
		err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		upgradeDBInstanceEngineVersionParam["DBInstanceIdentifier"] = d.Id()
		if _, ok := upgradeDBInstanceEngineVersionParam["Engine"]; !ok {
			upgradeDBInstanceEngineVersionParam["Engine"] = d.Get("engine")
		}
		if _, ok := upgradeDBInstanceEngineVersionParam["EngineVersion"]; !ok {
			upgradeDBInstanceEngineVersionParam["EngineVersion"] = d.Get("engine_version")
		}

		// check parameter valid on upgrade
		paramsReq := make(map[string]interface{})
		_, _, err := processParameters(d, meta, &paramsReq, true)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}

		action := "UpgradeDBInstanceEngineVersion"
		logger.Debug(logger.ReqFormat, action, upgradeDBInstanceEngineVersionParam)
		_, err = conn.UpgradeDBInstanceEngineVersion(&upgradeDBInstanceEngineVersionParam)
		logger.Debug(logger.AllFormat, action, upgradeDBInstanceEngineVersionParam, err)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}

		err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		//clean old db_parameter_group_id
		err = resourceKsyunMysqlDeleteParameterGroup(d, meta)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		//query db
		req := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
		action = "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := conn.DescribeDBInstances(&req)
		logger.Debug(logger.AllFormat, action, req, resp, err)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		value, err := getSdkValue("Data.Instances.0.DBParameterGroupId", *resp)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		return d.Set("db_parameter_group_id", value)
	}
	return err
}

func modifyDBInstanceAvailabilityZone(d *schema.ResourceData, meta interface{}, oldType interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var err error
	var modifyDBInstanceAvailabilityZoneParam map[string]interface{}
	transform := map[string]SdkReqTransform{
		"availability_zone_1": {mapping: "AvailabilityZone.1"},
		"availability_zone_2": {mapping: "AvailabilityZone.2"},
	}
	modifyDBInstanceAvailabilityZoneParam, err = SdkRequestAutoMapping(d, resourceKsyunKrds(), true, transform, nil)
	if err != nil {
		return err
	}
	if len(modifyDBInstanceAvailabilityZoneParam) > 0 {
		if oldType == "RR" {
			return fmt.Errorf("error on updating instance , krds rr is not support update %s", "availability_zone")
		}
		err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
		modifyDBInstanceAvailabilityZoneParam["DBInstanceIdentifier"] = d.Id()
		action := "ModifyDBInstanceAvailabilityZone"
		logger.Debug(logger.ReqFormat, action, modifyDBInstanceAvailabilityZoneParam)
		_, err = conn.ModifyDBInstanceAvailabilityZone(&modifyDBInstanceAvailabilityZoneParam)
		logger.Debug(logger.AllFormat, action, modifyDBInstanceAvailabilityZoneParam, err)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	return err
}

func resourceKsyunMysqlUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	conn := meta.(*KsyunClient).krdsconn
	// defer to read
	defer func(d *schema.ResourceData, meta interface{}) {
		_err := resourceKsyunMysqlRead(d, meta)
		if err == nil {
			err = _err
		} else {
			if _err != nil {
				err = fmt.Errorf(err.Error()+" %s", _err)
			}
		}

	}(d, meta)
	//instance_type
	oldType, newType := d.GetChange("db_instance_type")
	//rebuild ModifyDBInstance
	err = modifyDBInstance(d, meta, oldType)
	if err != nil {
		return err
	}
	//rebuild ModifyDBInstanceType
	err = modifyDBInstanceType(d, meta, oldType, newType)
	if err != nil {
		return err
	}
	//rebuild ModifyDBInstanceSpec
	err = modifyDBInstanceSpec(d, meta)
	if err != nil {
		return err
	}

	//rebuild UpgradeDBInstanceEngineVersion
	err = upgradeDBInstanceEngineVersion(d, meta)
	if err != nil {
		return err
	}

	//rebuild ModifyDBInstanceAvailabilityZone
	err = modifyDBInstanceAvailabilityZone(d, meta, oldType)
	if err != nil {
		return err
	}
	//rebuild ModifyDBInstanceEip
	if d.HasChange("instance_has_eip") {
		err = allocateOrReleaseInstanceEip(d, meta)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	//rebuild ModifyParameters
	if d.HasChange("parameters") {
		err = modifyParameters(d, meta)
		if err != nil {
			return fmt.Errorf("error on updating instance , error is %e", err)
		}
	}
	//wait
	err = checkStatus(d, conn, nil, nil, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error on updating instance , error is %e", err)
	}
	return err
}

func checkStatus(d *schema.ResourceData, conn *krds.Krds, pending []string, target []string, timeout time.Duration) error {
	if pending == nil {
		pending = waitStatus
	}
	if target == nil {
		target = []string{tActiveStatus}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    mysqlInstanceStateRefresh(conn, d.Id(), target),
	}
	_, err := stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("error on get krds DBInstanceStatus, err = %s", err)
	}
	return nil
}

func resourceKsyunMysqlDeleteDb(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	deleteReq := make(map[string]interface{})
	deleteReq["DBInstanceIdentifier"] = d.Id()
	var dbParameterGroupId string

	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		readReq := map[string]interface{}{"DBInstanceIdentifier": d.Id()}
		discribeAction := "DescribeInstances"
		logger.Debug(logger.ReqFormat, discribeAction, readReq)
		desResp, desErr := conn.DescribeDBInstances(&readReq)
		logger.Debug(logger.AllFormat, discribeAction, readReq, *desResp, desErr)

		if desErr != nil {
			if notFoundError(desErr) {
				return nil
			} else {
				return resource.NonRetryableError(desErr)
			}
		}

		bodyData := (*desResp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)
		dbParameterGroupId = sqlserverInstance["DBParameterGroupId"].(string)

		if state != tDeletedStatus {
			deleteAction := "DeleteDBInstance"
			logger.Debug(logger.ReqFormat, deleteAction, deleteReq)
			deleteResp, deleteErr := conn.DeleteDBInstance(&deleteReq)
			logger.Debug(logger.AllFormat, deleteAction, deleteReq, deleteResp, deleteErr)
			if deleteErr == nil || notFoundError(deleteErr) {
				return nil
			}
			if deleteErr != nil {
				return resource.RetryableError(deleteErr)
			}

			logger.Debug(logger.ReqFormat, discribeAction, readReq)
			postDesResp, postDesErr := conn.DescribeDBInstances(&readReq)
			logger.Debug(logger.AllFormat, discribeAction, readReq, *postDesResp, postDesErr)

			if desErr != nil {
				if notFoundError(desErr) {
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("error on  reading krds when delete %q, %s", d.Id(), desErr))
				}
			}
		}

		return resource.RetryableError(desErr)
	})
}

func resourceKsyunMysqlDeleteParameterGroup(d *schema.ResourceData, meta interface{}) error {
	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		conn := meta.(*KsyunClient).krdsconn
		if d.Get("db_parameter_group_id") != nil && d.Get("db_parameter_group_id").(string) != "" {
			delParam := make(map[string]interface{})
			delParam["DBParameterGroupId"] = d.Get("db_parameter_group_id").(string)
			_, deleteErr := conn.DeleteDBParameterGroup(&delParam)
			if deleteErr == nil || notFoundErrorNew(deleteErr) {
				return nil
			} else {
				return resource.RetryableError(deleteErr)
			}
		}
		return resource.RetryableError(nil)
	})
}

func resourceKsyunMysqlDelete(d *schema.ResourceData, meta interface{}) error {
	err := resourceKsyunMysqlDeleteDb(d, meta)
	if err != nil {
		return err
	}
	return resourceKsyunMysqlDeleteParameterGroup(d, meta)
}
