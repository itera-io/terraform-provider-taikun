package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOrganizationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Description: "Address.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"billing_email": {
			Description: "Billing email.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		// TODO bound_rules?
		"city": {
			Description: "City.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"country": {
			Description: "Country.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"created_at": {
			Description: "Time and date of creation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"discount_rate": {
			Description:  "Discount rate, must be between 0 and 100 (included).",
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"email": {
			Description: "Email.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"full_name": {
			Description:  "Full name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"id": {
			Description: "Organization's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"managers_can_change_subscription": {
			Description: "Allow subscription to be changed by managers.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"lock": {
			Description: "Indicates whether to lock the organization.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_read_only": {
			Description: "Whether the organization is in read-only mode.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"name": {
			Description: "Organization's name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"partner_id": {
			Description: "ID of the organization's partner.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"partner_name": {
			Description: "Name of the organization's partner.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"phone": {
			Description: "Phone number.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"projects": {
			Description: "Number of associated projects.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"servers": {
			Description: "Number of associated servers.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"vat_number": {
			Description: "VAT number.",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func resourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization",
		CreateContext: resourceTaikunOrganizationCreate,
		ReadContext:   generateResourceTaikunOrganizationRead(false),
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema:        resourceTaikunOrganizationSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
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
		IsEligibleUpdateSubscription: data.Get("managers_can_change_subscription").(bool),
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

	if isLocked, isLockedIsSet := data.GetOk("lock"); isLockedIsSet {
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

	return readAfterCreateWithRetries(generateResourceTaikunOrganizationRead(true), ctx, data, meta)
}

func generateResourceTaikunOrganizationRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id := data.Id()
		id32, _ := atoi32(data.Id())
		data.SetId("")

		response, err := apiClient.client.Organizations.OrganizationsList(organizations.NewOrganizationsListParams().WithV(ApiVersion).WithID(&id32), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if isAfterUpdateOrCreate {
				data.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawOrganization := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunOrganization(rawOrganization))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(id)

		return nil
	}
}

func resourceTaikunOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
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
		IsEligibleUpdateSubscription: data.Get("managers_can_change_subscription").(bool),
		IsLocked:                     data.Get("lock").(bool),
		Name:                         data.Get("name").(string),
		Phone:                        data.Get("phone").(string),
		VatNumber:                    data.Get("vat_number").(string),
	}

	updateLockParams := organizations.NewOrganizationsUpdateParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.client.Organizations.OrganizationsUpdate(updateLockParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunOrganizationRead(true), ctx, data, meta)
}

func resourceTaikunOrganizationDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func flattenTaikunOrganization(rawOrganization *models.OrganizationDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"address":                          rawOrganization.Address,
		"billing_email":                    rawOrganization.BillingEmail,
		"city":                             rawOrganization.City,
		"country":                          rawOrganization.Country,
		"created_at":                       rawOrganization.CreatedAt,
		"discount_rate":                    rawOrganization.DiscountRate,
		"email":                            rawOrganization.Email,
		"full_name":                        rawOrganization.FullName,
		"id":                               i32toa(rawOrganization.ID),
		"managers_can_change_subscription": rawOrganization.IsEligibleUpdateSubscription,
		"lock":                             rawOrganization.IsLocked,
		"is_read_only":                     rawOrganization.IsReadOnly,
		"name":                             rawOrganization.Name,
		"partner_id":                       i32toa(rawOrganization.PartnerID),
		"partner_name":                     rawOrganization.PartnerName,
		"phone":                            rawOrganization.Phone,
		"projects":                         rawOrganization.Projects,
		"servers":                          rawOrganization.Servers,
		"vat_number":                       rawOrganization.VatNumber,
	}
}
