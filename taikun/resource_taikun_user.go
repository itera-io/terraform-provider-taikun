package taikun

import (
	"context"
	"regexp"

	"github.com/itera-io/taikungoclient/client/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/models"
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
	apiClient := meta.(*apiClient)

	body := &models.CreateUserCommand{
		Username:    d.Get("user_name").(string),
		DisplayName: d.Get("display_name").(string),
		Email:       d.Get("email").(string),
		Role:        getUserRole(d.Get("role").(string)),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := users.NewUsersCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Users.UsersCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetPayload().ID)

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
		apiClient := meta.(*apiClient)
		id := d.Id()
		d.SetId("")

		response, err := apiClient.client.Users.UsersList(users.NewUsersListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawUser := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunUser(rawUser))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)

		return nil
	}
}

func resourceTaikunUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.UpdateUserCommand{
		ID:                  d.Id(),
		DisplayName:         d.Get("display_name").(string),
		Username:            d.Get("user_name").(string),
		Email:               d.Get("email").(string),
		Role:                getUserRole(d.Get("role").(string)),
		IsApprovedByPartner: true,
	}

	updateUserParams := users.NewUsersUpdateUserParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.client.Users.UsersUpdateUser(updateUserParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunUserReadWithRetries(), ctx, d, meta)
}

func resourceTaikunUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := users.NewUsersDeleteParams().WithV(ApiVersion).WithID(d.Id())
	_, _, err := apiClient.client.Users.UsersDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunUser(rawUser *models.UserForListDto) map[string]interface{} {

	return map[string]interface{}{
		"id":                         rawUser.ID,
		"user_name":                  rawUser.Username,
		"organization_id":            i32toa(rawUser.OrganizationID),
		"organization_name":          rawUser.OrganizationName,
		"role":                       rawUser.Role,
		"email":                      rawUser.Email,
		"display_name":               rawUser.DisplayName,
		"email_confirmed":            rawUser.IsEmailConfirmed,
		"email_notification_enabled": rawUser.IsEmailNotificationEnabled,
		"is_csm":                     rawUser.IsCsm,
		"is_disabled":                rawUser.IsLocked,
		"is_approved_by_partner":     rawUser.IsApprovedByPartner,
		"is_owner":                   rawUser.Owner,
	}
}
