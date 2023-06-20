package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"
	"regexp"

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
			ValidateDiagFunc: stringIsEmail,
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
		"is_approved_by_partner": {
			Description: "Indicates whether the user account is approved by its Partner. If it isn't, the user won't be able to login.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"is_csm": {
			Description: "Indicates whether the user is a Customer Success Manager.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"is_disabled": {
			Description: "Indicates whether the user is locked.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"is_owner": {
			Description: "Indicates whether the user is the Owner of their organization.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"organization_id": {
			Description:      "The ID of the user's organization.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the user's organization.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"role": {
			Description:  "The role of the user: `Manager` or `User`.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"User", "Manager"}, false),
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

func resourceTaikunUser() *schema.Resource {
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
	body.SetRole(tkcore.UserRole(d.Get("role").(string)))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	result, res, err := apiClient.Client.UsersApi.UsersCreate(context.TODO()).CreateUserCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(result.GetId())

	return readAfterCreateWithRetries(generateResourceTaikunUserReadWithRetries(), ctx, d, meta)
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

		response, res, err := apiClient.Client.UsersApi.UsersList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawUser := response.Data[0]

		err = setResourceDataFromMap(d, flattenTaikunUser(rawUser))
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
	body.SetRole(tkcore.UserRole(d.Get("role").(string)))
	body.SetIsApprovedByPartner(true)

	res, err := apiClient.Client.UsersApi.UsersUpdateUser(context.TODO()).UpdateUserCommand(body).Execute()

	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return readAfterUpdateWithRetries(generateResourceTaikunUserReadWithRetries(), ctx, d, meta)
}

func resourceTaikunUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	res, err := apiClient.Client.UsersApi.UsersDelete(context.TODO(), d.Id()).Execute()

	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunUser(rawUser tkcore.UserForListDto) map[string]interface{} {

	return map[string]interface{}{
		"id":                         rawUser.GetId(),
		"user_name":                  rawUser.GetUsername(),
		"organization_id":            i32toa(rawUser.GetOrganizationId()),
		"organization_name":          rawUser.GetOrganizationName(),
		"role":                       rawUser.GetRole(),
		"email":                      rawUser.GetEmail(),
		"display_name":               rawUser.GetDisplayName(),
		"email_confirmed":            rawUser.GetIsEmailConfirmed(),
		"email_notification_enabled": rawUser.GetIsEmailNotificationEnabled(),
		"is_csm":                     rawUser.GetIsCsm(),
		"is_disabled":                rawUser.GetIsLocked(),
		"is_approved_by_partner":     rawUser.GetIsApprovedByPartner(),
		"is_owner":                   rawUser.GetOwner(),
	}
}
