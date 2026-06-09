package account

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Account Data Source",
		ReadContext: dataSourceTaikunAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Account's ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Account's name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "Account's email.",
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
			"created_at": {
				Description: "Time and date of creation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_2fa_enabled": {
				Description: "Indicates whether 2FA is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceTaikunAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	response, res, err := apiClient.Client.AccountsAPI.AccountsDetails(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	err = utils.SetResourceDataFromMap(d, flattenTaikunAccount(response))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.I32toa(id))

	return nil
}
