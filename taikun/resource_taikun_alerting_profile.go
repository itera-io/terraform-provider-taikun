package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		"lock": {
			Description: "Indicates whether to lock the profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
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
		ReadContext:   generateResourceTaikunAlertingProfileReadWithoutRetries(),
		UpdateContext: resourceTaikunAlertingProfileUpdate,
		DeleteContext: resourceTaikunAlertingProfileDelete,
		Schema:        resourceTaikunAlertingProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunAlertingProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateAlertingProfileCommand{
		Name: *tkcore.NewNullableString(stringPtr(d.Get("name").(string))),
	}

	apiClient.Client.AlertingProfilesApi.AlertingprofilesCreate(ctx).CreateAlertingProfileCommand(body)

	if _, emailsIsSet := d.GetOk("emails"); emailsIsSet {
		body.SetEmails(getEmailDTOsFromAlertingProfileResourceData(d))
	}

	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetOrganizationId(organizationID)
	}

	if reminderData, reminderIsSet := d.GetOk("reminder"); reminderIsSet {
		body.SetReminder(tkcore.AlertingReminder(reminderData.(string)))
	}

	if slackConfigIDData, slackConfigIDIsSet := d.GetOk("slack_configuration_id"); slackConfigIDIsSet {
		slackConfigID, err := atoi32(slackConfigIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if slackConfigID != 0 {
			body.SetSlackConfigurationId(slackConfigID)
		}
	}

	if _, webhookIsSet := d.GetOk("webhook"); webhookIsSet {
		body.SetWebhooks(getWebhookDTOsFromAlertingProfileResourceData(d))
	}

	if _, integrationIsSet := d.GetOk("integration"); integrationIsSet {
		body.SetAlertingIntegrations(getIntegrationDTOsFromAlertingProfileResourceData(d))
	}

	response, bodyResponse, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesCreate(context.TODO()).CreateAlertingProfileCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(bodyResponse, err))
	}
	id, err := atoi32(response.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunAlertingProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunAlertingProfileReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunAlertingProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunAlertingProfileRead(true)
}
func generateResourceTaikunAlertingProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunAlertingProfileRead(false)
}
func generateResourceTaikunAlertingProfileRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, _, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}
		alertingProfileDTO := response.Data[0]

		alertingIntegrationsResponse, _, err := apiClient.Client.AlertingIntegrationsApi.AlertingintegrationsList(context.TODO(), alertingProfileDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(err)
		}

		err = setResourceDataFromMap(d, flattenTaikunAlertingProfile(&alertingProfileDTO, alertingIntegrationsResponse))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(alertingProfileDTO.GetId()))

		return nil
	}
}

func resourceTaikunAlertingProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunAlertingProfileLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "organization_id", "reminder", "slack_configuration_id") {
		body := tkcore.UpdateAlertingProfileCommand{}
		body.SetId(id)
		body.SetName(d.Get("name").(string))

		if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
			organizationID, err := atoi32(organizationIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			body.SetOrganizationId(organizationID)
		}
		if reminderData, reminderIsSet := d.GetOk("reminder"); reminderIsSet {
			body.SetReminder(tkcore.AlertingReminder(reminderData.(string)))
		}
		if slackConfigIDData, slackConfigIDIsSet := d.GetOk("slack_configuration_id"); slackConfigIDIsSet {
			slackConfigID, err := atoi32(slackConfigIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			body.SetSlackConfigurationId(slackConfigID)
		}

		_, res, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesEdit(ctx).UpdateAlertingProfileCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChange("emails") {
		alertEmails := getEmailDTOsFromAlertingProfileResourceData(d)
		res, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesAssignEmail(ctx, id).AlertingEmailDto(alertEmails).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChange("webhook") {
		webhooks := getWebhookDTOsFromAlertingProfileResourceData(d)
		res, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesAssignWebhooks(ctx, id).AlertingWebhookDto(webhooks).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if err := resourceTaikunAlertingProfileUpdateIntegrations(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunAlertingProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunAlertingProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunAlertingProfileUpdateIntegrations(d *schema.ResourceData, id int32, apiClient *tk.Client) (err error) {
	if !d.HasChange("integration") {
		return
	}

	// Remove old integrations
	oldIntegrationsData, _ := d.GetChange("integration")
	oldIntegrations := oldIntegrationsData.([]interface{})
	for _, oldIntegrationData := range oldIntegrations {
		oldIntegration := oldIntegrationData.(map[string]interface{})
		oldIntegrationID, _ := atoi32(oldIntegration["id"].(string))
		res, err := apiClient.Client.AlertingIntegrationsApi.AlertingintegrationsDelete(context.TODO(), oldIntegrationID).Execute()
		if err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	// Set new integrations
	if _, integrationIsSet := d.GetOk("integration"); integrationIsSet {
		alertingIntegrationDTOs := getIntegrationDTOsFromAlertingProfileResourceData(d)
		for _, alertingIntegration := range alertingIntegrationDTOs {
			alertingIntegrationCreateBody := tkcore.CreateAlertingIntegrationCommand{}
			alertingIntegrationCreateBody.SetAlertingIntegrationType(alertingIntegration.GetAlertingIntegrationType())
			alertingIntegrationCreateBody.SetToken(alertingIntegration.GetToken())
			alertingIntegrationCreateBody.SetUrl(alertingIntegration.GetUrl())
			alertingIntegrationCreateBody.SetAlertingProfileId(id)

			_, res, err := apiClient.Client.AlertingIntegrationsApi.AlertingintegrationsCreate(context.TODO()).CreateAlertingIntegrationCommand(alertingIntegrationCreateBody).Execute()
			if err != nil {
				err = tk.CreateError(res, err)
				return err
			}
		}
	}
	return
}

func resourceTaikunAlertingProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if res, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesDelete(context.TODO(), id).Execute(); err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func getAlertingProfileEmailsResourceFromEmailDTOs(emailDTOs []tkcore.AlertingEmailDto) []string {
	emails := make([]string, len(emailDTOs))
	for i, emailDTO := range emailDTOs {
		emails[i] = emailDTO.GetEmail()
	}
	return emails
}

func getAlertingProfileWebhookResourceFromWebhookDTOs(webhookDTOs []tkcore.AlertingWebhookDto) []map[string]interface{} {
	webhooks := make([]map[string]interface{}, len(webhookDTOs))
	for i, webhookDTO := range webhookDTOs {
		headers := make([]map[string]interface{}, len(webhookDTO.Headers))
		for i, rawHeader := range webhookDTO.Headers {
			headers[i] = map[string]interface{}{
				"key":   rawHeader.GetKey(),
				"value": rawHeader.GetValue(),
			}
		}
		webhooks[i] = map[string]interface{}{
			"header": headers,
			"url":    webhookDTO.GetUrl(),
		}
	}
	return webhooks
}

func getAlertingProfileIntegrationsResourceFromIntegrationDTOs(integrationDTOs []tkcore.AlertingIntegrationsListDto) []map[string]interface{} {
	integrations := make([]map[string]interface{}, len(integrationDTOs))
	for i, integrationDTO := range integrationDTOs {
		integrations[i] = map[string]interface{}{
			"id":    i32toa(integrationDTO.GetId()),
			"token": integrationDTO.GetToken(),
			"type":  integrationDTO.GetAlertingIntegrationType(),
			"url":   integrationDTO.GetUrl(),
		}
	}
	return integrations
}

func getEmailDTOsFromAlertingProfileResourceData(d *schema.ResourceData) []tkcore.AlertingEmailDto {
	emails := d.Get("emails").([]interface{})
	emailDTOs := make([]tkcore.AlertingEmailDto, len(emails))
	for i, email := range emails {
		emailDTOs[i] = tkcore.AlertingEmailDto{}
		emailDTOs[i].SetEmail(email.(string))
	}
	return emailDTOs
}

func getWebhookDTOsFromAlertingProfileResourceData(d *schema.ResourceData) []tkcore.AlertingWebhookDto {
	webhooks := d.Get("webhook").(*schema.Set).List()
	alertingWebhookDTOs := make([]tkcore.AlertingWebhookDto, len(webhooks))
	for i, webhookData := range webhooks {
		webhook := webhookData.(map[string]interface{})
		headers := webhook["header"].(*schema.Set).List()
		headerDTOs := make([]tkcore.WebhookHeaderDto, len(headers))
		for i, headerData := range headers {
			header := headerData.(map[string]interface{})
			headerDTOs[i] = tkcore.WebhookHeaderDto{}
			headerDTOs[i].SetKey(header["key"].(string))
			headerDTOs[i].SetValue(header["value"].(string))
		}
		alertingWebhookDTOs[i] = tkcore.AlertingWebhookDto{}
		alertingWebhookDTOs[i].SetHeaders(headerDTOs)
		alertingWebhookDTOs[i].SetUrl(webhook["url"].(string))
	}
	return alertingWebhookDTOs
}

func getIntegrationDTOsFromAlertingProfileResourceData(d *schema.ResourceData) []tkcore.AlertingIntegrationDto {
	integrations := d.Get("integration").([]interface{})
	alertingIntegrationDTOs := make([]tkcore.AlertingIntegrationDto, len(integrations))
	for i, integrationData := range integrations {
		integration := integrationData.(map[string]interface{})
		alertingIntegrationDTOs[i] = tkcore.AlertingIntegrationDto{}
		alertingIntegrationDTOs[i].SetAlertingIntegrationType(getAlertingIntegrationType(integration["type"].(string)))
		alertingIntegrationDTOs[i].SetToken(integration["token"].(string))
		alertingIntegrationDTOs[i].SetUrl(integration["url"].(string))
	}
	return alertingIntegrationDTOs
}

func flattenTaikunAlertingProfile(alertingProfileDTO *tkcore.AlertingProfilesListDto, alertingIntegrationDto []tkcore.AlertingIntegrationsListDto) map[string]interface{} {
	return map[string]interface{}{
		"created_by":               alertingProfileDTO.GetCreatedBy(),
		"emails":                   getAlertingProfileEmailsResourceFromEmailDTOs(alertingProfileDTO.GetEmails()),
		"id":                       i32toa(alertingProfileDTO.GetId()),
		"integration":              getAlertingProfileIntegrationsResourceFromIntegrationDTOs(alertingIntegrationDto),
		"lock":                     alertingProfileDTO.GetIsLocked(),
		"last_modified":            alertingProfileDTO.GetLastModified(),
		"last_modified_by":         alertingProfileDTO.GetLastModifiedBy(),
		"name":                     alertingProfileDTO.GetName(),
		"organization_id":          i32toa(alertingProfileDTO.GetOrganizationId()),
		"organization_name":        alertingProfileDTO.GetOrganizationName(),
		"reminder":                 alertingProfileDTO.GetReminder(),
		"slack_configuration_id":   i32toa(alertingProfileDTO.GetSlackConfigurationId()),
		"slack_configuration_name": alertingProfileDTO.GetSlackConfigurationName(),
		"webhook":                  getAlertingProfileWebhookResourceFromWebhookDTOs(alertingProfileDTO.GetWebhooks()),
	}
}

func resourceTaikunAlertingProfileLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.AlertingProfilesLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	res, err := apiClient.Client.AlertingProfilesApi.AlertingprofilesLockManager(context.TODO()).AlertingProfilesLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
