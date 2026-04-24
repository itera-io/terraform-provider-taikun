package user

import (
	"context"
	"regexp"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Description:  "The user's display name.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringLenBetween(3, 64),
		},
		"email": {
			Description:      "The email of the user.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: utils.StringIsEmail,
		},
		"email_confirmed": {
			Description: "Indicates whether the email of the user has been confirmed.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"email_notification_enabled": {
			Description: "Indicates whether the user has enabled notifications on their email.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"id": {
			Description: "The UUID of the user.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"account_id": {
			Description:      "The ID of the account the user belongs to.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"owner": {
			Description: "Indicates whether the user is a project owner.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"global_role": {
			Description:  "The role of the user.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"None", "Admin", "AccountAdmin", "AccountOwner"}, false),
		},
		"user_name": {
			Description: "The name of the user.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
	}
}

func ResourceTaikunUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun User",
		CreateContext: resourceTaikunUserCreate,
		ReadContext:   generateResourceTaikunUserReadWithoutRetries(),
		UpdateContext: resourceTaikunUserUpdate,
		DeleteContext: resourceTaikunUserDelete,
		Schema:        resourceTaikunUserSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateUserCommand{}
	body.SetUsername(d.Get("user_name").(string))
	body.SetDisplayName(d.Get("display_name").(string))
	body.SetEmail(d.Get("email").(string))
	body.SetIsAccountAdmin(d.Get("owner").(bool))

	accountIDData, accountIDIsSet := d.GetOk("account_id")
	if accountIDIsSet {
		accountID, err := utils.Atoi32(accountIDData.(string))
		if err != nil {
			return diag.Errorf("account_id isn't valid: %s", d.Get("account_id").(string))
		}
		body.SetAccountId(accountID)
	}

	result, res, err := apiClient.Client.UsersAPI.UsersCreate(context.TODO()).CreateUserCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(result.GetId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunUserReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunUserReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunUserRead(true)
}

func generateResourceTaikunUserReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunUserRead(false)
}

func generateResourceTaikunUserRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id := d.Id()
		d.SetId("")

		searchBody := tkcore.UsersSearchCommand{}
		searchBody.SetSearchTerm(id)
		searchRes, res, err := apiClient.Client.SearchAPI.SearchUsers(ctx).UsersSearchCommand(searchBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		var accountId int32
		var found bool
		for _, user := range searchRes.GetData() {
			if user.GetId() == id {
				accountId = user.GetAccountId()
				found = true
				break
			}
		}

		if !found {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		response, res, err := apiClient.Client.AccountsAPI.AccountsAccountUserDetails(ctx, accountId, id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		err = utils.SetResourceDataFromMap(d, flattenTaikunUser(response, accountId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)
		return nil
	}
}

func resourceTaikunUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.UpdateUserCommand{}
	body.SetId(d.Id())
	body.SetDisplayName(d.Get("display_name").(string))
	body.SetUsername(d.Get("user_name").(string))
	body.SetEmail(d.Get("email").(string))
	body.SetForceToResetPassword(d.Get("force_to_reset_password").(bool))
	body.SetDisable(d.Get("disable").(bool))

	res, err := apiClient.Client.UsersAPI.UsersUpdateUser(ctx).UpdateUserCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunUserReadWithRetries(), ctx, d, meta)
}

func resourceTaikunUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	res, err := apiClient.Client.UsersAPI.UsersDelete(ctx, d.Id()).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunUser(rawUser *tkcore.UserDetailsDto, accountID int32) map[string]interface{} {
	organizations := make([]map[string]interface{}, 0)
	for _, orgContext := range rawUser.GetOrganizations() {
		prjs := make([]map[string]interface{}, 0)
		for _, prj := range orgContext.GetProjects() {
			prjs = append(prjs, map[string]interface{}{
				"id":   prj.GetId(),
				"name": prj.GetName(),
			})
		}
		organizations = append(organizations, map[string]interface{}{
			"id":               orgContext.GetId(),
			"name":             orgContext.GetName(),
			"projects":         prjs,
			"organizationRole": orgContext.GetOrganizationRole(),
			"groupId":          orgContext.GetGroupId(),
			"groupName":        orgContext.GetGroupName(),
		})
	}

	return map[string]interface{}{
		"id":            rawUser.GetId(),
		"user_name":     rawUser.GetName(),
		"display_name":  rawUser.GetDisplayName(),
		"email":         rawUser.GetEmail(),
		"organizations": organizations,
		"account_id":    utils.I32toa(accountID),
		"global_role":   string(*rawUser.GlobalRole),
	}
}
