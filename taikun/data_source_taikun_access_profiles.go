package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/ssh_users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all access profiles.",
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema: map[string]*schema.Schema{
			"access_profiles": {
				Description: "List of retrieved access profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAccessProfileSchema(),
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

func dataSourceTaikunAccessProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var accessProfilesList []*models.AccessProfilesListDto
	for {
		response, err := apiClient.client.AccessProfiles.AccessProfilesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		accessProfilesList = append(accessProfilesList, response.GetPayload().Data...)
		if len(accessProfilesList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(accessProfilesList))
		params = params.WithOffset(&offset)
	}

	accessProfiles := make([]map[string]interface{}, len(accessProfilesList))
	for i, rawAccessProfile := range accessProfilesList {

		sshParams := ssh_users.NewSSHUsersListParams().WithV(ApiVersion).WithAccessProfileID(rawAccessProfile.ID)
		sshResponse, err := apiClient.client.SSHUsers.SSHUsersList(sshParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		accessProfiles[i] = flattenTaikunAccessProfile(rawAccessProfile, sshResponse)
	}
	if err := d.Set("access_profiles", accessProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
