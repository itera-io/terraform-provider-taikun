package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/slack"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunSlackConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"channel": {
			Description: "Slack channel for notifications.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"id": {
			Description: "ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 40),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"organization_id": {
			Description:      "Organization ID.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "Organization Name.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"type": {
			Description:  "Alert (receive only alert-type of notification) or General (receive all notifications).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Alert", "General"}, false),
		},
		"url": {
			Description:  "Webhook URL from Slack app.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

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
		Schema: resourceTaikunSlackConfigurationSchema(),
	}
}

func resourceTaikunSlackConfigurationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.UpsertSlackConfigurationCommand{
		Name:      data.Get("name").(string),
		URL:       data.Get("url").(string),
		Channel:   data.Get("channel").(string),
		SlackType: getSlackConfigurationType(data.Get("type").(string)),
	}

	if organizationIDData, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.OrganizationID = organizationID
	}

	params := slack.NewSlackCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Slack.SlackCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(response.Payload))

	return resourceTaikunSlackConfigurationRead(ctx, data, meta)
}

func resourceTaikunSlackConfigurationRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if len(response.Payload.Data) != 1 {
		return nil
	}

	rawSlackConfiguration := response.Payload.Data[0]

	err = setResourceDataFromMap(data, flattenTaikunSlackConfiguration(rawSlackConfiguration))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(id))

	return nil
}

func resourceTaikunSlackConfigurationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.UpsertSlackConfigurationCommand{
		ID:        id,
		Name:      data.Get("name").(string),
		URL:       data.Get("url").(string),
		Channel:   data.Get("channel").(string),
		SlackType: getSlackConfigurationType(data.Get("type").(string)),
	}

	if organizationIDData, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.OrganizationID = organizationID
	}

	params := slack.NewSlackCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Slack.SlackCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(response.Payload))

	return resourceTaikunSlackConfigurationRead(ctx, data, meta)
}

func resourceTaikunSlackConfigurationDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteSlackConfigurationCommand{ID: id}
	params := slack.NewSlackDeleteParams().WithV(ApiVersion).WithBody(&body)
	_, _, err = apiClient.client.Slack.SlackDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return nil
}

func flattenTaikunSlackConfiguration(rawSlackConfiguration *models.SlackConfigurationDto) map[string]interface{} {
	return map[string]interface{}{
		"channel":           rawSlackConfiguration.Channel,
		"id":                i32toa(rawSlackConfiguration.ID),
		"name":              rawSlackConfiguration.Name,
		"organization_id":   i32toa(rawSlackConfiguration.OrganizationID),
		"organization_name": rawSlackConfiguration.OrganizationName,
		"type":              rawSlackConfiguration.SlackType,
		"url":               rawSlackConfiguration.URL,
	}
}
