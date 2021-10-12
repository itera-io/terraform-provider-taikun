package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunSlackConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"channel": {
			Description: "Slack channel for notifications",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description:  "ID",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: stringIsInt,
		},
		"name": {
			Description: "Name",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"organization_id": {
			Description: "Organization ID",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"organization_name": {
			Description: "Organization Name",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"type": {
			Description: "Alert (receive only alert-type of notification) or General (receive all notifications)",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"url": {
			Description: "Webhook URL from Slack app",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func dataSourceTaikunSlackConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Get a slack configuration by its ID.",
		ReadContext: dataSourceTaikunSlackConfigurationRead,
		Schema:      dataSourceTaikunSlackConfigurationSchema(),
	}
}

func dataSourceTaikunSlackConfigurationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))
	return resourceTaikunSlackConfigurationRead(ctx, data, meta)
}
