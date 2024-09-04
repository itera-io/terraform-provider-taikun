package alerting_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunAlertingProfiles() *schema.Resource {
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
				ValidateDiagFunc: utils.StringIsInt,
			},
		},
	}
}

func dataSourceTaikunAlertingProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.AlertingProfilesAPI.AlertingprofilesList(context.TODO())

	if organizationIDData, organizationIDProvided := d.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var alertingProfileDTOs []tkcore.AlertingProfilesListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		alertingProfileDTOs = append(alertingProfileDTOs, response.Data...)
		if len(alertingProfileDTOs) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(alertingProfileDTOs))
	}

	alertingProfiles := make([]map[string]interface{}, len(alertingProfileDTOs))
	for i, alertingProfileDTO := range alertingProfileDTOs {

		alertingIntegrationsResponse, res, err := apiClient.Client.AlertingIntegrationsAPI.AlertingintegrationsList(context.TODO(), alertingProfileDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		alertingProfiles[i] = flattenTaikunAlertingProfile(&alertingProfileDTO, alertingIntegrationsResponse)
	}

	if err := d.Set("alerting_profiles", alertingProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
