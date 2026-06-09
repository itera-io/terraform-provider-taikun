package group

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Groups Data Source",
		ReadContext: dataSourceTaikunGroupsRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account ID to list groups for.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"groups": {
				Description: "List of groups.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Group's ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Group's name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"claim_value": {
							Description: "Claim value for the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	accountId := int32(d.Get("account_id").(int))

	response, res, err := apiClient.Client.GroupsAPI.GroupsList(ctx).AccountId(accountId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	err = d.Set("groups", flattenGroups(response.GetData()))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("groups")

	return nil
}

func flattenGroups(groups []tkcore.GroupListItem) []interface{} {
	result := make([]interface{}, len(groups))
	for i, group := range groups {
		result[i] = map[string]interface{}{
			"id":          utils.I32toa(group.GetId()),
			"name":        group.GetName(),
			"claim_value": group.GetClaimValue(),
		}
	}
	return result
}
