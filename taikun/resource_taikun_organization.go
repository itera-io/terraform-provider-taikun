package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization",
		CreateContext: resourceTaikunOrganizationCreate,
		ReadContext:   resourceTaikunOrganizationRead,
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO bound_rules?
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_credentials": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"discount_rate": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(0, 100),
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"let_managers_change_subscription": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// TODO partner details?
			"partner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partner_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servers": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vat_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTaikunOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := data.Id()
	data.SetId("")

	var limit int32 = 1
	response, err := apiClient.client.Organizations.OrganizationsList(organizations.NewOrganizationsListParams().WithV(ApiVersion).WithSearchID(&id).WithLimit(&limit), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(response.GetPayload().Data) != 1 {
		return nil
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
