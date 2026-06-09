package group

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Group Data Source",
		ReadContext: dataSourceTaikunGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Group's ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "Account ID the group belongs to.",
				Type:        schema.TypeInt,
				Required:    true,
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
	}
}

func dataSourceTaikunGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := int32(d.Get("account_id").(int))

	response, res, err := apiClient.Client.GroupsAPI.GroupsList(ctx).AccountId(accountId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	for _, item := range response.GetData() {
		if item.GetId() == id {
			err = utils.SetResourceDataFromMap(d, flattenTaikunGroup(&item, accountId))
			if err != nil {
				return diag.FromErr(err)
			}
			d.SetId(utils.I32toa(id))
			return nil
		}
	}

	return diag.Errorf("group with ID %d not found in account %d", id, accountId)
}
