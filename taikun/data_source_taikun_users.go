package taikun

import (
	"context"
	"github.com/itera-io/taikungoclient/client/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of users, optionally filtered by organization.",
		ReadContext: dataSourceTaikunUsersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunUserSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunUsersRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := users.NewUsersListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var rawUserList []*models.UserForListDto
	for {
		response, err := apiClient.client.Users.UsersList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		rawUserList = append(rawUserList, response.GetPayload().Data...)
		if len(rawUserList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(rawUserList))
		params = params.WithOffset(&offset)
	}

	userList := make([]map[string]interface{}, len(rawUserList), len(rawUserList))
	for i, rawUser := range rawUserList {
		userList[i] = flattenDataSourceTaikunUserItem(rawUser)
	}
	if err := data.Set("users", userList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunUserItem(rawUser *models.UserForListDto) map[string]interface{} {

	return map[string]interface{}{
		"id":                         rawUser.ID,
		"user_name":                  rawUser.Username,
		"organization_id":            i32toa(rawUser.OrganizationID),
		"organization_name":          rawUser.OrganizationName,
		"role":                       rawUser.Role,
		"email":                      rawUser.Email,
		"display_name":               rawUser.DisplayName,
		"email_confirmed":            rawUser.IsEmailConfirmed,
		"email_notification_enabled": rawUser.IsEmailNotificationEnabled,
		"is_csm":                     rawUser.IsCsm,
		"is_owner":                   rawUser.Owner,
		"user_disabled":              rawUser.IsLocked,
		"approved_by_partner":        rawUser.IsApprovedByPartner,
	}
}
