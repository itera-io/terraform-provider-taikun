package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunSlackConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Slack Configuration",
		CreateContext: resourceTaikunSlackConfigurationCreate,
		ReadContext:   resourceTaikunSlackConfigurationRead,
		UpdateContext: resourceTaikunSlackConfigurationUpdate,
		DeleteContext: resourceTaikunSlackConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"channel": {
				Description: "Slack channel for notifications",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"organization_id": {
				Description: "Organization ID",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"organization_name": {
				Description: "Organization Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description:  "Alert (receive only alert-type of notification) or General (receive all notifications)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Alert", "General"}, false),
			},
			"url": {
				Description: "Webhook URL from Slack app",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceTaikunSlackConfigurationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTaikunSlackConfigurationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTaikunSlackConfigurationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTaikunSlackConfigurationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
