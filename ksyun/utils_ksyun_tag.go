package ksyun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func mergeTagsData(d *schema.ResourceData, data *map[string]interface{}, client *KsyunClient, resourceType string) (err error) {
	var tags []interface{}
	tagService := TagService{client}
	tags, err = tagService.ReadTagByResourceId(d, d.Id(), resourceType)
	if err != nil {
		return err
	}
	tagMap := make(map[string]interface{})
	for _, tag := range tags {
		_m := tag.(map[string]interface{})
		tagMap[_m["TagKey"].(string)] = _m["TagValue"].(string)
	}
	if len(tagMap) > 0 {
		(*data)["Tags"] = tagMap
	}
	return err
}
