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
		Description: "Retrieve all users.",
		ReadContext: dataSourceTaikunUsersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"users": {
				Description: "List of retrieved users.",
				Type:        schema.TypeList,
				Computed:    true,
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
		userList[i] = flattenTaikunUser(rawUser)
	}
	if err := data.Set("users", userList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
