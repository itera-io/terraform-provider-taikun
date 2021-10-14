package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunAlertingProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// TODO "alerting_integrations"
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

func resourceTaikunAlertingProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Alerting Profile",
		CreateContext: resourceTaikunAlertingProfileCreate,
		ReadContext:   resourceTaikunAlertingProfileRead,
		UpdateContext: resourceTaikunAlertingProfileUpdate,
		DeleteContext: resourceTaikunAlertingProfileDelete,
		Schema:        resourceTaikunAlertingProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunAlertingProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion).WithID(&id)
	response, err := apiClient.client.AlertingProfiles.AlertingProfilesList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(response.Payload.Data) != 1 {
		return diag.Errorf("Alerting profile with ID %d not found", id)
	}

	rawAlertingProfile := response.Payload.Data[0]

	emails := make([]string, len(rawAlertingProfile.Emails))
	for i, rawEmail := range rawAlertingProfile.Emails {
		emails[i] = rawEmail.Email
	}

	webhooks := make([]map[string]interface{}, len(rawAlertingProfile.Webhooks))
	for i, rawWebhook := range rawAlertingProfile.Webhooks {
		headers := make([]map[string]interface{}, len(rawWebhook.Headers))
		for i, rawHeader := range rawWebhook.Headers {
			headers[i] = map[string]interface{}{
				"key":   rawHeader.Key,
				"value": rawHeader.Value,
			}
		}
		webhooks[i] = map[string]interface{}{
			"headers": headers,
			"url":     rawWebhook.URL,
		}
	}

	if err := data.Set("created_by", rawAlertingProfile.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("emails", emails); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", i32toa(rawAlertingProfile.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_locked", rawAlertingProfile.IsLocked); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("last_modified", rawAlertingProfile.LastModified); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("last_modified_by", rawAlertingProfile.LastModifiedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", rawAlertingProfile.Name); err != nil {
		return diag.FromErr(err)

	}
	if err := data.Set("organization_id", i32toa(rawAlertingProfile.OrganizationID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_name", rawAlertingProfile.OrganizationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("reminder", rawAlertingProfile.Reminder); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slack_configuration_id", rawAlertingProfile.SlackConfigurationID); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slack_configuration_name", rawAlertingProfile.SlackConfigurationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("webhooks", webhooks); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTaikunAlertingProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateAlertingProfileCommand{
		Name:     data.Get("name").(string),
		Reminder: getAlertingProfileReminder(data.Get("reminder").(string)),
	}

	if emailsData, emailsIsSet := data.GetOk("emails"); emailsIsSet {
		emails := emailsData.([]string)
		emailDTOs := make([]*models.AlertingEmailDto, len(emails))
		for i, email := range emails {
			emailDTOs[i].Email = email
		}
		body.Emails = emailDTOs
	}

	if organizationIDData, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.OrganizationID = organizationID
	}

	if slackConfigIDData, slackConfigIDIsSet := data.GetOk("slack_configuration_id"); slackConfigIDIsSet {
		slackConfigID, err := atoi32(slackConfigIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SlackConfigurationID = slackConfigID
	}

	if webhooksData, webhooksIsSet := data.GetOk("webhooks"); webhooksIsSet {
		webhooks := webhooksData.([]interface{})
		webhookDTOs := make([]*models.AlertingWebhookDto, len(webhooks))
		for i, webhookData := range webhooks {
			webhook := webhookData.(map[string]interface{})
			webhookDTOs[i].URL = webhook["url"].(string)
			headers := webhook["headers"].([]map[string]string)
			headerDTOs := make([]*models.WebhookHeaderDto, len(headers))
			for i, header := range headers {
				headerDTOs[i].Key = header["key"]
				headerDTOs[i].Value = header["value"]
			}
		}
		body.Webhooks = webhookDTOs
	}

	// TODO handle alerting integrations

	params := alerting_profiles.NewAlertingProfilesCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.AlertingProfiles.AlertingProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO data.SetId()

	return resourceTaikunAlertingProfileRead(ctx, data, meta)
}
func resourceTaikunAlertingProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}
func resourceTaikunAlertingProfileDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}