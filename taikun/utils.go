package taikun

import (
	"fmt"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/users"
)

func int32Address(value interface{}) *int32 {
	var intValue int
	switch v := value.(type) {
	case int:
		intValue = v
	case int32:
		intValue = int(v)
	default:
		panic(fmt.Sprintf("expected value to be int or int32, but got %T", value))
	}
	int32Value := int32(intValue)
	return &int32Value
}

func stringAddress(value interface{}) *string {
	strValue, ok := value.(string)
	if !ok {
		return nil
	}
	return &strValue
}

func strfmtEmailAddress(value interface{}) *strfmt.Email {
	strValue, ok := value.(string)
	if !ok {
		// Handle the case when value is not of type string
		// You may return a nil pointer or a default value depending on your use case
		return nil
	}
	email := strfmt.Email(strValue)
	return &email
}

func getDefaultOrganization(defaultOrganizationID *int32, apiClient *taikungoclient.Client) error {
	params := users.NewUsersDetailsParams().WithV(ApiVersion)
	response, err := apiClient.Client.Users.UsersDetails(params, apiClient)
	if err != nil {
		return err
	}
	*defaultOrganizationID = response.Payload.Data.OrganizationID
	return nil
}

func getOrganizationFromDataOrElseDefault(d *schema.ResourceData, apiClient *taikungoclient.Client) (organizationID int32, err error) {
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
