package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func checkConflictOnUpdate(key string, conflictKey string, d *schema.ResourceData) (err error) {
	if d.HasChange(conflictKey) && (d.HasChanges(key) || d.Get(key) != nil) {
		return fmt.Errorf("%s and %s conflict! ", conflictKey, key)
	}
	if d.HasChange(key) && (d.HasChanges(conflictKey) || d.Get(conflictKey) != nil) {
		return fmt.Errorf("%s and %s conflict! ", key, conflictKey)
	}
	return err
}

func conflictResourceImport(parent string, single string, multiple string, d *schema.ResourceData) ([]*schema.ResourceData, error) {
	var err error
	if !strings.Contains(d.Id(), ":") {
		err = d.Set(parent, d.Id())
		if err != nil {
			return nil, err
		}
		d.SetId(d.Id() + ":" + multiple)
	} else {
		items := strings.Split(d.Id(), ":")
		if len(items) != 2 {
			return nil, fmt.Errorf("id must split with %s and size %v", ":", 2)
		}
		err = d.Set(parent, items[0])
		if err != nil {
			return nil, err
		}
		err = d.Set(single, items[1])
		if err != nil {
			return nil, err
		}
		d.SetId(items[0] + ":" + items[1])
	}

	return []*schema.ResourceData{d}, err
}

func conflictResourceSetId(use string, parent string, single string, multiple string, d *schema.ResourceData) {
	if use == single {
		d.SetId(d.Get(parent).(string) + ":" + d.Get(single).(string))
	} else {
		d.SetId(d.Get(parent).(string) + ":" + multiple)
	}
}

func conflictResourceDiffSuppressForSingle(multiple string, old, new string, d *schema.ResourceData) bool {
	if d.Get(multiple) != nil {
		if v, ok := d.Get(multiple).(*schema.Set); ok && len(v.List()) > 0 {
			return true
		}
		if v, ok := d.Get(multiple).(string); ok && v != "" {
			return true
		}
	}
	if strings.Contains(d.Id(), ":"+multiple) {
		return true
	}
	if old != "" && new == "" {
		return true
	}
	return false
}

func conflictResourceDiffSuppressForMultiple(single string, multiple string, d *schema.ResourceData) bool {
	if d.Get(single) != nil && d.Get(single).(string) != "" {
		return true
	}

	if d.Id() != "" && !strings.Contains(d.Id(), ":"+multiple) {
		return true
	}

	return false
}

func checkConflictOnCreate(key string, conflictKey string, d *schema.ResourceData) (use string, err error) {
	_, existKey := d.GetOk(key)
	_, existConflictKey := d.GetOk(conflictKey)
	if !existKey && !existConflictKey {
		return use, fmt.Errorf("must set %s or set %s ", key, conflictKey)
	}
	if existKey {
		use = key
	} else {
		use = conflictKey
	}
	return use, err
}

func checkMultipleExist(multiple string, d *schema.ResourceData) bool {
	return strings.Contains(d.Id(), ":"+multiple)
}
