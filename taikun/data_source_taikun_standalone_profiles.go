package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/security_group"
	"github.com/itera-io/taikungoclient/client/stand_alone_profile"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunStandaloneProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Standalone profiles.",
		ReadContext: dataSourceTaikunStandaloneProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"standalone_profiles": {
				Description: "List of retrieved Standalone profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunStandaloneProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunStandaloneProfilesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var standaloneProfilesListDtos []*models.StandAloneProfilesListDto
	for {
		response, err := apiClient.client.StandAloneProfile.StandAloneProfileList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		standaloneProfilesListDtos = append(standaloneProfilesListDtos, response.GetPayload().Data...)
		if len(standaloneProfilesListDtos) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(standaloneProfilesListDtos))
		params = params.WithOffset(&offset)
	}

	standaloneProfiles := make([]map[string]interface{}, len(standaloneProfilesListDtos))
	for i, rawStandaloneProfile := range standaloneProfilesListDtos {

		params := security_group.NewSecurityGroupListParams().WithV(ApiVersion).WithStandAloneProfileID(rawStandaloneProfile.ID)
		securityGroupResponse, err := apiClient.client.SecurityGroup.SecurityGroupList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		standaloneProfiles[i] = flattenTaikunStandaloneProfile(rawStandaloneProfile, securityGroupResponse.GetPayload())
	}
	if err := data.Set("standalone_profiles", standaloneProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
