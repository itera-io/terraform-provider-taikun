package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/alerting_integrations"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunAlertingProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all alerting profiles.",
		ReadContext: dataSourceTaikunAlertingProfilesRead,
		Schema: map[string]*schema.Schema{
			"alerting_profiles": {
				Description: "List of retrieved alerting profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAlertingProfileSchema(),
				},
			},
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
		},
	}
}

func dataSourceTaikunAlertingProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	dataSourceID := "all"

	params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion)
	if organizationIDData, organizationIDProvided := d.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var alertingProfileDTOs []*models.AlertingProfilesListDto
	for {
		response, err := apiClient.Client.AlertingProfiles.AlertingProfilesList(params, apiClient)
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

		alertingIntegrationsParams := alerting_integrations.NewAlertingIntegrationsListParams().WithV(ApiVersion).WithAlertingProfileID(alertingProfileDTO.ID)
		alertingIntegrationsResponse, err := apiClient.Client.AlertingIntegrations.AlertingIntegrationsList(alertingIntegrationsParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		alertingProfiles[i] = flattenTaikunAlertingProfile(alertingProfileDTO, alertingIntegrationsResponse.Payload)
	}

	if err := d.Set("alerting_profiles", alertingProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
