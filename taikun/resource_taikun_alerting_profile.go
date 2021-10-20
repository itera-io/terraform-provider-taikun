package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/alerting_integrations"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunAlertingProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The profile creator.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"emails": {
			Description: "The list of e-mails to notify.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"id": {
			Description: "The alerting profile's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"integration": {
			Description: "list of alerting integrations",
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"token": {
						Description: "token (required from Opsgenie, Pagerduty and Splunk)",
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
						ForceNew:    true,
					},
					"type": {
						Description: "type of integration (Opsgenie, Pagerduty, Splunk or MicrosoftTeams)",
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Opsgenie",
							"Pagerduty",
							"Splunk",
							"MicrosoftTeams",
						}, false),
					},
					"url": {
						Description:  "URL",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},
		"is_locked": {
			Description: "Whether the profile is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"last_modified": {
			Description: "The time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The alerting profile's name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"reminder": {
			Description: "The frequency of notifications (HalfHour, Hourly, Daily or None).",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"HalfHour",
				"Hourly",
				"Daily",
				"None",
			}, false),
		},
		"slack_configuration_id": {
			Description:      "The ID of the Slack configuration to notify.",
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "0",
			ValidateDiagFunc: stringIsInt,
		},
		"slack_configuration_name": {
			Description: "The name of the Slack configuration to notify.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"webhook": {
			Description: "The list of webhooks to notify.",
			Type:        schema.TypeSet,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"header": {
						Description: "The list of headers.",
						Type:        schema.TypeSet,
						Optional:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Description: "The header key.",
									Type:        schema.TypeString,
									Required:    true,
								},
								"value": {
									Description: "The header value.",
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
					"url": {
						Description: "The webhook URL.",
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

func resourceTaikunAlertingProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateAlertingProfileCommand{
		Name: data.Get("name").(string),
	}

	if _, emailsIsSet := data.GetOk("emails"); emailsIsSet {
		body.Emails = getEmailDTOsFromAlertingProfileResourceData(data)
	}

	if organizationIDData, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.OrganizationID = organizationID
	}

	if reminderData, reminderIsSet := data.GetOk("reminder"); reminderIsSet {
		body.Reminder = getAlertingProfileReminder(reminderData.(string))
	}

	if slackConfigIDData, slackConfigIDIsSet := data.GetOk("slack_configuration_id"); slackConfigIDIsSet {
		slackConfigID, err := atoi32(slackConfigIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SlackConfigurationID = slackConfigID
	}

	if _, webhookIsSet := data.GetOk("webhook"); webhookIsSet {
		body.Webhooks = getWebhookDTOsFromAlertingProfileResourceData(data)
	}

	params := alerting_profiles.NewAlertingProfilesCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.AlertingProfiles.AlertingProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(response.Payload.ID)
	id, _ := atoi32(response.Payload.ID)

	if _, integrationIsSet := data.GetOk("integration"); integrationIsSet {
		alertingIntegrationDTOs := getIntegrationDTOsFromAlertingProfileResourceData(data)
		for _, alertingIntegration := range alertingIntegrationDTOs {
			alertingIntegrationCreateBody := models.CreateAlertingIntegrationCommand{
				AlertingIntegration: &models.AlertingIntegrationDto{
					AlertingIntegrationType: alertingIntegration.AlertingIntegrationType,
					Token:                   alertingIntegration.Token,
					URL:                     alertingIntegration.URL,
				},
				AlertingProfileID: id,
			}
			alertingIntegrationParams := alerting_integrations.NewAlertingIntegrationsCreateParams().WithV(ApiVersion).WithBody(&alertingIntegrationCreateBody)
			_, err = apiClient.client.AlertingIntegrations.AlertingIntegrationsCreate(alertingIntegrationParams, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if isLocked, isLockedIsSet := data.GetOk("is_locked"); isLockedIsSet && isLocked.(bool) {
		body := models.AlertingProfilesLockManagerCommand{
			ID:   id,
			Mode: "lock",
		}
		params := alerting_profiles.NewAlertingProfilesLockManagerParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.AlertingProfiles.AlertingProfilesLockManager(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunAlertingProfileRead(ctx, data, meta)
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
	alertingProfileDTO := response.Payload.Data[0]

	alertingIntegrationsParams := alerting_integrations.NewAlertingIntegrationsListParams().WithV(ApiVersion).WithAlertingProfileID(alertingProfileDTO.ID)
	alertingIntegrationsResponse, err := apiClient.client.AlertingIntegrations.AlertingIntegrationsList(alertingIntegrationsParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	alertingIntegrationDTOs := alertingIntegrationsResponse.Payload

	if err := data.Set("created_by", alertingProfileDTO.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("emails", getAlertingProfileEmailsResourceFromEmailDTOs(alertingProfileDTO.Emails)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", i32toa(alertingProfileDTO.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("integration", getAlertingProfileIntegrationsResourceFromIntegrationDTOs(alertingIntegrationDTOs)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_locked", alertingProfileDTO.IsLocked); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("last_modified", alertingProfileDTO.LastModified); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("last_modified_by", alertingProfileDTO.LastModifiedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", alertingProfileDTO.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_id", i32toa(alertingProfileDTO.OrganizationID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_name", alertingProfileDTO.OrganizationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("reminder", alertingProfileDTO.Reminder); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slack_configuration_id", i32toa(alertingProfileDTO.SlackConfigurationID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("slack_configuration_name", alertingProfileDTO.SlackConfigurationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("webhook", getAlertingProfileWebhookResourceFromWebhookDTOs(alertingProfileDTO.Webhooks)); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(alertingProfileDTO.ID))

	return nil
}

func resourceTaikunAlertingProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChanges("name", "organization_id", "reminder", "slack_configuration_id") {
		body := models.UpdateAlertingProfileCommand{
			ID:   id,
			Name: data.Get("name").(string),
		}
		if organizationIDData, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
			organizationID, err := atoi32(organizationIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			body.OrganizationID = organizationID
		}
		if reminderData, reminderIsSet := data.GetOk("reminder"); reminderIsSet {
			body.Reminder = getAlertingProfileReminder(reminderData.(string))
		}
		if slackConfigIDData, slackConfigIDIsSet := data.GetOk("slack_configuration_id"); slackConfigIDIsSet {
			slackConfigID, err := atoi32(slackConfigIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			body.SlackConfigurationID = slackConfigID
		}
		params := alerting_profiles.NewAlertingProfilesEditParams().WithV(ApiVersion).WithBody(&body)
		response, err := apiClient.client.AlertingProfiles.AlertingProfilesEdit(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		data.SetId(response.Payload.ID)
	}

	if data.HasChange("is_locked") {
		body := models.AlertingProfilesLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		params := alerting_profiles.NewAlertingProfilesLockManagerParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.AlertingProfiles.AlertingProfilesLockManager(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("emails") {
		body := getEmailDTOsFromAlertingProfileResourceData(data)
		params := alerting_profiles.NewAlertingProfilesAssignEmailsParams().WithV(ApiVersion).WithID(id).WithBody(body)
		_, err := apiClient.client.AlertingProfiles.AlertingProfilesAssignEmails(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("webhook") {
		body := getWebhookDTOsFromAlertingProfileResourceData(data)
		params := alerting_profiles.NewAlertingProfilesAssignWebhooksParams().WithV(ApiVersion).WithID(id).WithBody(body)
		_, err := apiClient.client.AlertingProfiles.AlertingProfilesAssignWebhooks(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunAlertingProfileRead(ctx, data, meta)
}

func resourceTaikunAlertingProfileDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteAlertingProfilesCommand{ID: id}
	params := alerting_profiles.NewAlertingProfilesDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.client.AlertingProfiles.AlertingProfilesDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func getAlertingProfileEmailsResourceFromEmailDTOs(emailDTOs []*models.AlertingEmailDto) []string {
	emails := make([]string, len(emailDTOs))
	for i, emailDTO := range emailDTOs {
		emails[i] = emailDTO.Email
	}
	return emails
}

func getAlertingProfileWebhookResourceFromWebhookDTOs(webhookDTOs []*models.AlertingWebhookDto) []map[string]interface{} {
	webhooks := make([]map[string]interface{}, len(webhookDTOs))
	for i, webhookDTO := range webhookDTOs {
		headers := make([]map[string]interface{}, len(webhookDTO.Headers))
		for i, rawHeader := range webhookDTO.Headers {
			headers[i] = map[string]interface{}{
				"key":   rawHeader.Key,
				"value": rawHeader.Value,
			}
		}
		webhooks[i] = map[string]interface{}{
			"header": headers,
			"url":    webhookDTO.URL,
		}
	}
	return webhooks
}

func getAlertingProfileIntegrationsResourceFromIntegrationDTOs(integrationDTOs []*models.AlertingIntegrationsListDto) []map[string]interface{} {
	integrations := make([]map[string]interface{}, len(integrationDTOs))
	for i, integrationDTO := range integrationDTOs {
		integrations[i] = map[string]interface{}{
			"token": integrationDTO.Token,
			"type":  integrationDTO.AlertingIntegrationType,
			"url":   integrationDTO.URL,
		}
	}
	return integrations
}

func getEmailDTOsFromAlertingProfileResourceData(data *schema.ResourceData) []*models.AlertingEmailDto {
	emails := data.Get("emails").([]interface{})
	emailDTOs := make([]*models.AlertingEmailDto, len(emails))
	for i, email := range emails {
		emailDTOs[i] = &models.AlertingEmailDto{
			Email: email.(string),
		}
	}
	return emailDTOs
}

func getWebhookDTOsFromAlertingProfileResourceData(data *schema.ResourceData) []*models.AlertingWebhookDto {
	webhooks := data.Get("webhook").(*schema.Set).List()
	alertingWebhookDTOs := make([]*models.AlertingWebhookDto, len(webhooks))
	for i, webhookData := range webhooks {
		webhook := webhookData.(map[string]interface{})
		headers := webhook["header"].(*schema.Set).List()
		headerDTOs := make([]*models.WebhookHeaderDto, len(headers))
		for i, headerData := range headers {
			header := headerData.(map[string]interface{})
			headerDTOs[i] = &models.WebhookHeaderDto{
				Key:   header["key"].(string),
				Value: header["value"].(string),
			}
		}
		alertingWebhookDTOs[i] = &models.AlertingWebhookDto{
			Headers: headerDTOs,
			URL:     webhook["url"].(string),
		}
	}
	return alertingWebhookDTOs
}

func getIntegrationDTOsFromAlertingProfileResourceData(data *schema.ResourceData) []*models.AlertingIntegrationDto {
	integrations := data.Get("integration").([]interface{})
	alertingIntegrationDTOs := make([]*models.AlertingIntegrationDto, len(integrations))
	for i, integrationData := range integrations {
		integration := integrationData.(map[string]interface{})
		alertingIntegrationDTOs[i] = &models.AlertingIntegrationDto{
			AlertingIntegrationType: getAlertingIntegrationType(integration["type"].(string)),
			Token:                   integration["token"].(string),
			URL:                     integration["url"].(string),
		}
	}
	return alertingIntegrationDTOs
}
