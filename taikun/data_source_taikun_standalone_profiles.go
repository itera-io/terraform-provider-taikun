package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/security_group"
	"github.com/itera-io/taikungoclient/client/stand_alone_profile"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunStandaloneProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all standalone profiles.",
		ReadContext: dataSourceTaikunStandaloneProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"standalone_profiles": {
				Description: "List of retrieved standalone profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunStandaloneProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunStandaloneProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	dataSourceID := "all"

	params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
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
		response, err := apiClient.Client.StandAloneProfile.StandAloneProfileList(params, apiClient)
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
		securityGroupResponse, err := apiClient.Client.SecurityGroup.SecurityGroupList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		standaloneProfiles[i] = flattenTaikunStandaloneProfile(rawStandaloneProfile, securityGroupResponse.GetPayload())
	}
	if err := d.Set("standalone_profiles", standaloneProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
