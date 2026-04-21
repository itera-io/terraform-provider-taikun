package user

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all users.",
		ReadContext: dataSourceTaikunUsersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
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

func dataSourceTaikunUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	var rawUserList []tkcore.UsersSearchResponseData

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}

		dropdownRes, _, err := apiClient.Client.UsersAPI.UsersDropdown(ctx).OrganizationId(organizationID).Execute()
		if err != nil {
			return diag.FromErr(err)
		}

		for _, u := range dropdownRes.GetData() {
			searchBody := tkcore.UsersSearchCommand{}
			searchBody.SetSearchTerm(u.GetId())
			searchRes, _, err := apiClient.Client.SearchAPI.SearchUsers(ctx).UsersSearchCommand(searchBody).Execute()
			if err == nil {
				for _, su := range searchRes.GetData() {
					if su.GetId() == u.GetId() {
						rawUserList = append(rawUserList, su)
						break
					}
				}
			}
		}
	} else {
		searchBody := tkcore.UsersSearchCommand{}
		searchRes, _, err := apiClient.Client.SearchAPI.SearchUsers(ctx).UsersSearchCommand(searchBody).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		rawUserList = searchRes.GetData()
	}

	userList := make([]map[string]interface{}, len(rawUserList))
	for i, rawUser := range rawUserList {
		userList[i] = map[string]interface{}{
			"id":           rawUser.GetId(),
			"user_name":    rawUser.GetName(),
			"email":        rawUser.GetEmail(),
			"account_id":   utils.I32toa(rawUser.GetAccountId()),
			"account_name": rawUser.GetAccountName(),
		}
	}
	if err := d.Set("users", userList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)
	return nil
}
