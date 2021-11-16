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
			Description: "The creator of the alerting profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"emails": {
			Description: "The list of emails to notify.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"id": {
			Description: "The alerting profile's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"integration": {
			Description: "List of alerting integrations.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "The alerting integration's ID.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"token": {
						Description: "The token (required by Opsgenie, Pagerduty and Splunk).",
						Type:        schema.TypeString,
						Optional:    true,
						Default:     "",
					},
					"type": {
						Description: "The type of integration: `Opsgenie`, `Pagerduty`, `Splunk` or `MicrosoftTeams`.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Opsgenie",
							"Pagerduty",
							"Splunk",
							"MicrosoftTeams",
						}, false),
					},
					"url": {
						Description:  "The alerting integration's URL.",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
				},
			},
		},
		"lock": {
			Description: "Indicates whether to lock the profile.",
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
			Description: "The frequency of notifications: `HalfHour`, `Hourly`, `Daily` or `None`.",
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
		ReadContext:   generateResourceTaikunAlertingProfileRead(false),
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
	id, err := atoi32(response.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.Payload.ID)

	if err := resourceTaikunAlertingProfileSetIntegrations(data, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if data.Get("lock").(bool) {
		if err := resourceTaikunAlertingProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunAlertingProfileRead(true), ctx, data, meta)
}

func generateResourceTaikunAlertingProfileRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			if isAfterUpdateOrCreate {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}
		alertingProfileDTO := response.Payload.Data[0]

		alertingIntegrationsParams := alerting_integrations.NewAlertingIntegrationsListParams().WithV(ApiVersion).WithAlertingProfileID(alertingProfileDTO.ID)
		alertingIntegrationsResponse, err := apiClient.client.AlertingIntegrations.AlertingIntegrationsList(alertingIntegrationsParams, apiClient)
		if err != nil {
			if _, ok := err.(*alerting_integrations.AlertingIntegrationsListNotFound); ok && isAfterUpdateOrCreate {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return diag.FromErr(err)
		}

		err = setResourceDataFromMap(data, flattenTaikunAlertingProfile(alertingProfileDTO, alertingIntegrationsResponse.Payload))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(alertingProfileDTO.ID))

		return nil
	}
}

func resourceTaikunAlertingProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := data.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunAlertingProfileLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
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
		_, err := apiClient.client.AlertingProfiles.AlertingProfilesEdit(params, apiClient)
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

	if err := resourceTaikunAlertingProfileUnsetIntegrations(data, apiClient); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceTaikunAlertingProfileSetIntegrations(data, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if data.Get("lock").(bool) {
		if err := resourceTaikunAlertingProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunAlertingProfileRead(true), ctx, data, meta)
}

func resourceTaikunAlertingProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceTaikunAlertingProfileUnsetIntegrations(data *schema.ResourceData, apiClient *apiClient) error {
	oldIntegrationsData, _ := data.GetChange("integration")
	oldIntegrations := oldIntegrationsData.([]interface{})
	for _, oldIntegrationData := range oldIntegrations {
		oldIntegration := oldIntegrationData.(map[string]interface{})
		oldIntegrationID, _ := atoi32(oldIntegration["id"].(string))
		params := alerting_integrations.NewAlertingIntegrationsDeleteParams().WithV(ApiVersion).WithID(oldIntegrationID)
		_, _, err := apiClient.client.AlertingIntegrations.AlertingIntegrationsDelete(params, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunAlertingProfileSetIntegrations(data *schema.ResourceData, id int32, apiClient *apiClient) error {
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
			_, err := apiClient.client.AlertingIntegrations.AlertingIntegrationsCreate(alertingIntegrationParams, apiClient)
			if err != nil {
				return err
			}
		}
	}
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
			"id":    i32toa(integrationDTO.ID),
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

func flattenTaikunAlertingProfile(alertingProfileDTO *models.AlertingProfilesListDto, alertingIntegrationDto []*models.AlertingIntegrationsListDto) map[string]interface{} {
	return map[string]interface{}{
		"created_by":               alertingProfileDTO.CreatedBy,
		"emails":                   getAlertingProfileEmailsResourceFromEmailDTOs(alertingProfileDTO.Emails),
		"id":                       i32toa(alertingProfileDTO.ID),
		"integration":              getAlertingProfileIntegrationsResourceFromIntegrationDTOs(alertingIntegrationDto),
		"lock":                     alertingProfileDTO.IsLocked,
		"last_modified":            alertingProfileDTO.LastModified,
		"last_modified_by":         alertingProfileDTO.LastModifiedBy,
		"name":                     alertingProfileDTO.Name,
		"organization_id":          i32toa(alertingProfileDTO.OrganizationID),
		"organization_name":        alertingProfileDTO.OrganizationName,
		"reminder":                 alertingProfileDTO.Reminder,
		"slack_configuration_id":   i32toa(alertingProfileDTO.SlackConfigurationID),
		"slack_configuration_name": alertingProfileDTO.SlackConfigurationName,
		"webhook":                  getAlertingProfileWebhookResourceFromWebhookDTOs(alertingProfileDTO.Webhooks),
	}
}

func resourceTaikunAlertingProfileLock(id int32, lock bool, apiClient *apiClient) error {
	body := models.AlertingProfilesLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := alerting_profiles.NewAlertingProfilesLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.AlertingProfiles.AlertingProfilesLockManager(params, apiClient)
	return err
}
