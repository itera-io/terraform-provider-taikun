package account

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTaikunAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Account's name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"email": {
			Description: "Account's email.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"create_organization": {
			Description: "Create an organization for this account.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
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
	}
}

func ResourceTaikunAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Account",
		CreateContext: resourceTaikunAccountCreate,
		ReadContext:   generateResourceTaikunAccountReadWithoutRetries(),
		UpdateContext: resourceTaikunAccountUpdate,
		DeleteContext: resourceTaikunAccountDelete,
		Schema:        resourceTaikunAccountSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateAccountCommand{}
	body.SetName(d.Get("name").(string))
	body.SetEmail(d.Get("email").(string))
	body.SetCreateOrganization(d.Get("create_organization").(bool))

	id, res, err := apiClient.Client.AccountsAPI.AccountsCreate(ctx).CreateAccountCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(utils.I32toa(id))

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunAccountReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunAccountReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunAccountRead(true)
}

func generateResourceTaikunAccountReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunAccountRead(false)
}

func generateResourceTaikunAccountRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.AccountsAPI.AccountsDetails(ctx, id).Execute()
		if err != nil {
			if res != nil && res.StatusCode == 404 {
				if withRetries {
					return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
				}
				d.SetId("")
				return nil
			}
			return diag.FromErr(tk.CreateError(res, err))
		}

		err = utils.SetResourceDataFromMap(d, flattenTaikunAccount(response))
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceTaikunAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.UpdateAccountCommand{}
	body.SetId(id)
	body.SetName(d.Get("name").(string))
	body.SetEmail(d.Get("email").(string))

	_, res, err := apiClient.Client.AccountsAPI.AccountsUpdate(ctx).UpdateAccountCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunAccountReadWithRetries(), ctx, d, meta)
}

func resourceTaikunAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.AccountsAPI.AccountsDelete(ctx, id).Execute()
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunAccount(rawAccount *tkcore.AccountDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"name":                rawAccount.GetName(),
		"email":               rawAccount.GetEmail(),
		"organizations_count": rawAccount.GetOrganizationsCount(),
		"users_count":         rawAccount.GetUsersCount(),
		"groups_count":        rawAccount.GetGroupsCount(),
		"projects_count":      rawAccount.GetProjectsCount(),
		"created_at":          rawAccount.GetCreatedAt(),
		"is_2fa_enabled":      rawAccount.GetIs2FAEnabled(),
	}
}
