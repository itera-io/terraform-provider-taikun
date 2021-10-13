package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunAlertingProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "profile creator",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"emails": {
			Description: "list of e-mails to notify",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"id": {
			Description: "ID",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "whether the profile is locked or not",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"last_modified": {
			Description: "time and date of last modification",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "last user to have modified the profile",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "name",
			Type:        schema.TypeString,
			Required:    true,
		},
		"organization_id": {
			Description:  "ID of the organization which owns the profile",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: stringIsInt,
		},
		"organization_name": {
			Description: "name of the organization which owns the profile",
			Type:        schema.TypeString,
			Computed:    true,
		},
		// TODO add "projects" ?
		"reminder": {
			Description: "frequency of notifications (None, HalfHour, Hourly or Daily)",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"None",
				"HalfHour",
				"Hourly",
				"Daily",
			}, false),
		},
		"slack_configuration_id": {
			Description:  "ID of Slack configuration to notify",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: stringIsInt,
		},
		"slack_configuration_name": {
			Description: "name of Slack configuration to notify",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"webhooks": {
			Description: "list of webhooks to notify",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"headers": {
						Description: "list of headers",
						Type:        schema.TypeList,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Description: "key",
									Type:        schema.TypeString,
									Required:    true,
								},
								"value": {
									Description: "value",
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
					"url": {
						Description: "URL",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
			},
		},
	}
}
