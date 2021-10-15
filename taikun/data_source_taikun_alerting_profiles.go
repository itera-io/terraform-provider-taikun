package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunAlertingProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of alerting profiles for your organizations, or filter by organization if Partner or Admin",
		ReadContext: dataSourceTaikunAlertingProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "organization ID filter (for Partner and Admin roles)",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"alerting_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAlertingProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunAlertingProfilesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion)
	if organizationIDData, organizationIDProvided := data.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var alertingProfileDTOs []*models.AlertingProfilesListDto
	for {
		response, err := apiClient.client.AlertingProfiles.AlertingProfilesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		alertingProfileDTOs = append(alertingProfileDTOs, response.Payload.Data...)
		if len(alertingProfileDTOs) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(alertingProfileDTOs))
		params = params.WithOffset(&offset)
	}

	alertingProfiles := make([]map[string]interface{}, len(alertingProfileDTOs))
	for i, alertingProfileDTO := range alertingProfileDTOs {
		alertingProfiles[i] = flattenDataSourceTaikunAlertingProfilesItem(alertingProfileDTO)
	}

	if err := data.Set("alerting_profiles", alertingProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunAlertingProfilesItem(alertingProfileDTO *models.AlertingProfilesListDto) map[string]interface{} {
	return map[string]interface{}{
		"created_by":               alertingProfileDTO.CreatedBy,
		"emails":                   getAlertingProfileEmailsResourceFromEmailDTOs(alertingProfileDTO.Emails),
		"id":                       i32toa(alertingProfileDTO.ID),
		"is_locked":                alertingProfileDTO.IsLocked,
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
