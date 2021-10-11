package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/slack"
)

func resourceTaikunSlackConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Slack Configuration",
		CreateContext: resourceTaikunSlackConfigurationCreate,
		ReadContext:   resourceTaikunSlackConfigurationRead,
		DeleteContext: resourceTaikunSlackConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"channel": {
				Description: "Slack channel for notifications",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				ForceNew:    true,
			},
			"organization_id": {
				Description: "Organization ID",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
				ForceNew:     true,
			},
			"url": {
				Description:  "Webhook URL from Slack app",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
		},
	}
}

func resourceTaikunSlackConfigurationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	params := slack.NewSlackListParams().WithV(ApiVersion).WithID(&id)
	response, err := apiClient.client.Slack.SlackList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if response.Payload.TotalCount != 1 {
		return diag.Errorf("Slack configuration with ID %d not found", id)
	}

	rawSlackConfiguration := response.Payload.Data[0]
	if err := data.Set("channel", rawSlackConfiguration.Channel); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", rawSlackConfiguration.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", rawSlackConfiguration.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_id", rawSlackConfiguration.OrganizationID); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_name", rawSlackConfiguration.OrganizationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slack_type", rawSlackConfiguration.SlackType); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("url", rawSlackConfiguration.URL); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(id))

	return nil
}

func resourceTaikunSlackConfigurationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil

}

func resourceTaikunSlackConfigurationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
