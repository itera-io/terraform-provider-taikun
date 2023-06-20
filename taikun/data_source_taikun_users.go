package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceTaikunUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0
	params := apiClient.Client.UsersApi.UsersList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var rawUserList []tkcore.UserForListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		rawUserList = append(rawUserList, response.Data...)
		if len(rawUserList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(rawUserList))
	}

	userList := make([]map[string]interface{}, len(rawUserList))
	for i, rawUser := range rawUserList {
		userList[i] = flattenTaikunUser(rawUser)
	}
	if err := d.Set("users", userList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
