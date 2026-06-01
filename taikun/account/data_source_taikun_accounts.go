package account

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunAccounts() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Accounts Data Source",
		ReadContext: dataSourceTaikunAccountsRead,
		Schema: map[string]*schema.Schema{
			"accounts": {
				Description: "List of accounts.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Account's ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Account's name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organizations_count": {
							Description: "Number of organizations in the account.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"users_count": {
							Description: "Number of users in the account.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"groups_count": {
							Description: "Number of groups in the account.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"projects_count": {
							Description: "Number of projects in the account.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunAccountsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	accounts := make([]tkcore.AccountList, 0)
	hasMore := true
	var cursorId int32 = 0

	for hasMore {
		request := apiClient.Client.AccountsAPI.AccountsListAccounts(ctx)
		if cursorId != 0 {
			request = request.CursorId(cursorId)
		}

		response, res, err := request.Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		accounts = append(accounts, response.Data...)

		hasMore = response.GetHasMore()
		if hasMore && response.NextCursor.IsSet() {
			cursorId = *response.NextCursor.Get()
		} else {
			hasMore = false
		}
	}

	err := d.Set("accounts", flattenAccounts(accounts))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("accounts")

	return nil
}

func flattenAccounts(accounts []tkcore.AccountList) []interface{} {
	result := make([]interface{}, len(accounts))
	for i, account := range accounts {
		result[i] = map[string]interface{}{
			"id":                  utils.I32toa(account.GetId()),
			"name":                account.GetName(),
			"organizations_count": account.GetOrganizationsCount(),
			"users_count":         account.GetUsersCount(),
			"groups_count":        account.GetGroupsCount(),
			"projects_count":      account.GetProjectsCount(),
		}
	}
	return result
}
