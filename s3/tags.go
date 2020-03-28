package s3

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func MakeUserTags(i interface{}) map[string]string {
	switch value := i.(type) {
	case map[string]interface{}:
		tags := make(map[string]string)
		for k, v := range value {
			tags[k] = v.(string)
		}
		return tags

	case map[string]string:
		tags := make(map[string]string)
		for k, v := range value {
			tags[k] = v
		}
		return tags

	default:
		return make(map[string]string)
	}
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}
