package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOrganizationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Description: "Address",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"billing_email": {
			Description: "Billing email",
			Type:        schema.TypeString,
			Optional:    true,
		},
		// TODO bound_rules?
		"city": {
			Description: "City",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"cloud_credentials": {
			Description: "Number of associated cloud credentials",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"country": {
			Description: "Country",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"created_at": {
			Description: "Time and date of creation",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"discount_rate": {
			Description:  "Discount rate, must be between 0 and 100 (included)",
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"email": {
			Description: "Email",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"full_name": {
			Description: "Full name",
			Type:        schema.TypeString,
			Required:    true,
		},
		"id": {
			Description: "ID",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"let_managers_change_subscription": {
			Description: "Allow subscription to be changed by managers",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"is_locked": {
			Description: "Whether the organization is locked",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_read_only": {
			Description: "Whether the organization is in read-only mode",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"name": {
			Description:  "Name",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: stringIsLowercase,
		},
		// TODO partner details?
		"partner_id": {
			Description: "ID of the organization's partner",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"partner_name": {
			Description: "Name of the organization's partner",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"phone": {
			Description: "Phone number",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"projects": {
			Description: "Number of associated projects",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"servers": {
			Description: "Number of associated servers",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"users": {
			Description: "Number of associated users",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"vat_number": {
			Description: "VAT number",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func resourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization",
		CreateContext: resourceTaikunOrganizationCreate,
		ReadContext:   resourceTaikunOrganizationRead,
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema:        resourceTaikunOrganizationSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := data.Id()
	id32, _ := atoi32(data.Id())
	data.SetId("")

	var limit int32 = 1
	response, err := apiClient.client.Organizations.OrganizationsList(organizations.NewOrganizationsListParams().WithV(ApiVersion).WithID(&id32).WithLimit(&limit), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(response.GetPayload().Data) != 1 || response.Payload.Data[0].ID != id32 {
		return diag.Errorf("Organization with ID %s not found", id)
	}

	rawOrganization := response.GetPayload().Data[0]

	if err := data.Set("address", rawOrganization.Address); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("billing_email", rawOrganization.BillingEmail); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("city", rawOrganization.City); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("cloud_credentials", rawOrganization.CloudCredentials); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("country", rawOrganization.Country); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("created_at", rawOrganization.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("discount_rate", rawOrganization.DiscountRate); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("email", rawOrganization.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("full_name", rawOrganization.FullName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", i32toa(rawOrganization.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("let_managers_change_subscription", rawOrganization.IsEligibleUpdateSubscription); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_locked", rawOrganization.IsLocked); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("is_read_only", rawOrganization.IsReadOnly); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", rawOrganization.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("partner_id", i32toa(rawOrganization.PartnerID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("partner_name", rawOrganization.PartnerName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("phone", rawOrganization.Phone); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("projects", rawOrganization.Projects); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("servers", rawOrganization.Servers); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("users", rawOrganization.Users); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("vat_number", rawOrganization.VatNumber); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(id)

	return nil
}

func resourceTaikunOrganizationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.OrganizationCreateCommand{
		Address:                      data.Get("address").(string),
		BillingEmail:                 data.Get("billing_email").(string),
		City:                         data.Get("city").(string),
		Country:                      data.Get("country").(string),
		DiscountRate:                 data.Get("discount_rate").(float64),
		Email:                        data.Get("email").(string),
		FullName:                     data.Get("full_name").(string),
		IsEligibleUpdateSubscription: data.Get("let_managers_change_subscription").(bool),
		Name:                         data.Get("name").(string),
		Phone:                        data.Get("phone").(string),
		VatNumber:                    data.Get("vat_number").(string),
	}

	params := organizations.NewOrganizationsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Organizations.OrganizationsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.GetPayload().ID)

	if isLocked, isLockedIsSet := data.GetOk("is_locked"); isLockedIsSet {
		id, _ := atoi32(createResult.GetPayload().ID)
		updateLockBody := &models.UpdateOrganizationCommand{
			Address:                      body.Address,
			BillingEmail:                 body.BillingEmail,
			City:                         body.City,
			Country:                      body.Country,
			DiscountRate:                 body.DiscountRate,
			Email:                        body.Email,
			FullName:                     body.FullName,
			ID:                           id,
			IsEligibleUpdateSubscription: body.IsEligibleUpdateSubscription,
			IsLocked:                     isLocked.(bool),
			Name:                         body.Name,
			Phone:                        body.Phone,
			VatNumber:                    body.VatNumber,
		}
		updateLockParams := organizations.NewOrganizationsUpdateParams().WithV(ApiVersion).WithBody(updateLockBody)
		_, err := apiClient.client.Organizations.OrganizationsUpdate(updateLockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunOrganizationRead(ctx, data, meta)
}

func resourceTaikunOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if !data.HasChanges(
		"address",
		"billing_email",
		"city",
		"country",
		"discount_rate",
		"email",
		"full_name",
		"let_managers_change_subscription",
		"is_locked",
		"name",
		"phone",
		"vat_number",
	) {
		return resourceTaikunOrganizationRead(ctx, data, meta)
	}

	body := &models.UpdateOrganizationCommand{
		Address:                      data.Get("address").(string),
		BillingEmail:                 data.Get("billing_email").(string),
		City:                         data.Get("city").(string),
		Country:                      data.Get("country").(string),
		DiscountRate:                 data.Get("discount_rate").(float64),
		Email:                        data.Get("email").(string),
		FullName:                     data.Get("full_name").(string),
		ID:                           id,
		IsEligibleUpdateSubscription: data.Get("let_managers_change_subscription").(bool),
		IsLocked:                     data.Get("is_locked").(bool),
		Name:                         data.Get("name").(string),
		Phone:                        data.Get("phone").(string),
		VatNumber:                    data.Get("vat_number").(string),
	}

	updateLockParams := organizations.NewOrganizationsUpdateParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.client.Organizations.OrganizationsUpdate(updateLockParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaikunOrganizationRead(ctx, data, meta)
}

func resourceTaikunOrganizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := organizations.NewOrganizationsDeleteParams().WithV(ApiVersion).WithOrganizationID(id)
	_, _, err = apiClient.client.Organizations.OrganizationsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
