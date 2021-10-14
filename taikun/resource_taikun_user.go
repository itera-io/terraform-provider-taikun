package taikun

import (
	"context"
	"github.com/itera-io/taikungoclient/client/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The UUID of the user.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"user_name": {
			Description: "The name of the user.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"organization_id": {
			Description:  "The id of the organization to which the user belongs.",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization to which the user belongs.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"role": {
			Description:  "The role of the user.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"User", "Manager"}, false),
		},
		"email": {
			Description: "The email of the user.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"display_name": {
			Description: "The name of the user displayed in the upper right corner.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
		},
		"email_confirmed": {
			Description: "Indicates whether the email of the user has been confirmed or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"email_notification_enabled": {
			Description: "Indicates whether the user has chosen to receive notifications on his email or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"is_csm": {
			Description: "Indicates whether the user is a Customer Success Manager or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"is_owner": {
			Description: "Indicates whether the user is the Owner of his organization.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"user_disabled": {
			Description: "Indicates whether the user is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"approved_by_partner": {
			Description: "Indicates whether the user account has been approved by a partner. If new user is not approved by partner, he won't be able to login.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
	}
}

func resourceTaikunUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun User",
		CreateContext: resourceTaikunUserCreate,
		ReadContext:   resourceTaikunUserRead,
		UpdateContext: resourceTaikunUserUpdate,
		DeleteContext: resourceTaikunUserDelete,
		Schema:        resourceTaikunUserSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunUserCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateUserCommand{
		Username:    data.Get("user_name").(string),
		DisplayName: data.Get("display_name").(string),
		Email:       data.Get("email").(string),
		Role:        getUserRole(data.Get("role").(string)),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := users.NewUsersCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Users.UsersCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.GetPayload().ID)

	updateUserBody := &models.UpdateUserCommand{
		ID:                  createResult.GetPayload().ID,
		Email:               data.Get("email").(string),
		Role:                getUserRole(data.Get("role").(string)),
		Username:            data.Get("user_name").(string),
		DisplayName:         data.Get("display_name").(string),
		Disable:             data.Get("user_disabled").(bool),
		IsApprovedByPartner: data.Get("approved_by_partner").(bool),
	}

	updateUserParams := users.NewUsersUpdateUserParams().WithV(ApiVersion).WithBody(updateUserBody)
	_, err = apiClient.client.Users.UsersUpdateUser(updateUserParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaikunUserRead(ctx, data, meta)
}

func resourceTaikunUserRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := data.Id()
	data.SetId("")

	response, err := apiClient.client.Users.UsersList(users.NewUsersListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawUser := response.GetPayload().Data[0]

		if err := data.Set("id", rawUser.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("user_name", rawUser.Username); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawUser.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawUser.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("role", rawUser.Role); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("email", rawUser.Email); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("display_name", rawUser.DisplayName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("email_confirmed", rawUser.IsEmailConfirmed); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("email_notification_enabled", rawUser.IsEmailNotificationEnabled); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_csm", rawUser.IsCsm); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_owner", rawUser.Owner); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("user_disabled", rawUser.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("approved_by_partner", rawUser.IsApprovedByPartner); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(id)
	}

	return nil
}

func resourceTaikunUserUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.UpdateUserCommand{
		ID:                  data.Id(),
		DisplayName:         data.Get("display_name").(string),
		Username:            data.Get("user_name").(string),
		Email:               data.Get("email").(string),
		Role:                getUserRole(data.Get("role").(string)),
		Disable:             data.Get("user_disabled").(bool),
		IsApprovedByPartner: data.Get("approved_by_partner").(bool),
	}

	updateUserParams := users.NewUsersUpdateUserParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.client.Users.UsersUpdateUser(updateUserParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaikunUserRead(ctx, data, meta)
}

func resourceTaikunUserDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := users.NewUsersDeleteParams().WithV(ApiVersion).WithID(data.Id())
	_, _, err := apiClient.client.Users.UsersDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
