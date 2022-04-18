package taikun

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/users"
)

func getDefaultOrganization(defaultOrganizationID *int32, apiClient *apiClient) error {
	params := users.NewUsersDetailsParams().WithV(ApiVersion)
	response, err := apiClient.client.Users.UsersDetails(params, apiClient)
	if err != nil {
		return err
	}
	*defaultOrganizationID = response.Payload.Data.OrganizationID
	return nil
}

func getOrganizationFromDataOrElseDefault(d *schema.ResourceData, apiClient *apiClient) (organizationID int32, err error) {
	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationID, err = atoi32(organizationIDData.(string))
		if err != nil {
			err = fmt.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
	} else {
		err = getDefaultOrganization(&organizationID, apiClient)
	}

	return
}
