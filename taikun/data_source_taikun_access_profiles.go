package taikun

import (
	"context"
	"fmt"

	"github.com/itera-io/taikungoclient/client/ssh_users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of access profiles, optionally filtered by organization.",
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"access_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAccessProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunAccessProfilesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
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

	accessProfiles := make([]map[string]interface{}, len(accessProfilesList), len(accessProfilesList))
	for i, rawAccessProfile := range accessProfilesList {

		sshParams := ssh_users.NewSSHUsersListParams().WithV(ApiVersion).WithAccessProfileID(rawAccessProfile.ID)
		sshResponse, err := apiClient.client.SSHUsers.SSHUsersList(sshParams, apiClient)
		if err != nil {
			fmt.Println(rawAccessProfile.ID)
			return diag.FromErr(err)
		}

		accessProfiles[i] = flattenDataSourceTaikunAccessProfilesItem(rawAccessProfile, sshResponse)
	}
	if err := data.Set("access_profiles", accessProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunAccessProfilesItem(rawAccessProfile *models.AccessProfilesListDto, sshResponse *ssh_users.SSHUsersListOK) map[string]interface{} {

	DNSServers := make([]map[string]interface{}, len(rawAccessProfile.DNSServers), len(rawAccessProfile.DNSServers))
	for i, rawDNSServer := range rawAccessProfile.DNSServers {
		DNSServers[i] = map[string]interface{}{
			"address": rawDNSServer.Address,
			"id":      i32toa(rawDNSServer.ID),
		}
	}

	NTPServers := make([]map[string]interface{}, len(rawAccessProfile.NtpServers), len(rawAccessProfile.NtpServers))
	for i, rawNTPServer := range rawAccessProfile.NtpServers {
		NTPServers[i] = map[string]interface{}{
			"address": rawNTPServer.Address,
			"id":      i32toa(rawNTPServer.ID),
		}
	}

	projects := make([]map[string]interface{}, len(rawAccessProfile.Projects), len(rawAccessProfile.Projects))
	for i, rawProject := range rawAccessProfile.Projects {
		projects[i] = map[string]interface{}{
			"id":   i32toa(rawProject.ID),
			"name": rawProject.Name,
		}
	}

	SSHUsers := make([]map[string]interface{}, len(sshResponse.Payload), len(sshResponse.Payload))
	for i, rawSSHUser := range sshResponse.Payload {
		SSHUsers[i] = map[string]interface{}{
			"id":         i32toa(rawSSHUser.ID),
			"name":       rawSSHUser.Name,
			"public_key": rawSSHUser.SSHPublicKey,
		}
	}

	return map[string]interface{}{
		"created_by":        rawAccessProfile.CreatedBy,
		"dns_server":        DNSServers,
		"http_proxy":        rawAccessProfile.HTTPProxy,
		"id":                i32toa(rawAccessProfile.ID),
		"is_locked":         rawAccessProfile.IsLocked,
		"last_modified":     rawAccessProfile.LastModified,
		"last_modified_by":  rawAccessProfile.LastModifiedBy,
		"name":              rawAccessProfile.Name,
		"ntp_server":        NTPServers,
		"organization_id":   i32toa(rawAccessProfile.OrganizationID),
		"organization_name": rawAccessProfile.OrganizationName,
		"project":           projects,
		"ssh_user":          SSHUsers,
	}
}
