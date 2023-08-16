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
			Optional:     true,
			Default:      100,
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
		"is_read_only": {
			Description: "Whether the organization is in read-only mode.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the organization.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"managers_can_change_subscription": {
			Description: "Allow subscription to be changed by managers.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
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
		ReadContext:   generateResourceTaikunOrganizationReadWithoutRetries(),
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema:        resourceTaikunOrganizationSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.OrganizationCreateCommand{}
	body.SetAddress(d.Get("address").(string))
	body.SetBillingEmail(d.Get("billing_email").(string))
	body.SetCity(d.Get("city").(string))
	body.SetCountry(d.Get("country").(string))
	body.SetDiscountRate(d.Get("discount_rate").(float64))
	body.SetEmail(d.Get("email").(string))
	body.SetFullName(d.Get("full_name").(string))
	body.SetIsEligibleUpdateSubscription(d.Get("managers_can_change_subscription").(bool))
	body.SetName(d.Get("name").(string))
	body.SetPhone(d.Get("phone").(string))
	body.SetVatNumber(d.Get("vat_number").(string))

	createResult, res, err := apiClient.Client.OrganizationsApi.OrganizationsCreate(context.TODO()).OrganizationCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if isLocked, isLockedIsSet := d.GetOk("lock"); isLockedIsSet {
		updateLockBody := tkcore.UpdateOrganizationCommand{}
		updateLockBody.SetAddress(body.GetAddress())
		updateLockBody.SetBillingEmail(body.GetBillingEmail())
		updateLockBody.SetCity(body.GetCity())
		updateLockBody.SetCountry(body.GetCountry())
		updateLockBody.SetDiscountRate(body.GetDiscountRate())
		updateLockBody.SetEmail(body.GetEmail())
		updateLockBody.SetFullName(body.GetFullName())
		updateLockBody.SetId(id)
		updateLockBody.SetIsEligibleUpdateSubscription(body.GetIsEligibleUpdateSubscription())
		updateLockBody.SetIsLocked(isLocked.(bool))
		updateLockBody.SetName(body.GetName())
		updateLockBody.SetPhone(body.GetPhone())
		updateLockBody.SetVatNumber(body.GetVatNumber())

		res, err := apiClient.Client.OrganizationsApi.OrganizationsUpdate(ctx).UpdateOrganizationCommand(updateLockBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunOrganizationReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(true)
}
func generateResourceTaikunOrganizationReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(false)
}
func generateResourceTaikunOrganizationRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id := d.Id()
		id32, _ := atoi32(d.Id())
		d.SetId("")

		response, res, err := apiClient.Client.OrganizationsApi.OrganizationsList(context.TODO()).Id(id32).Execute()

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

		rawOrganization := response.Data[0]

		err = setResourceDataFromMap(d, flattenTaikunOrganization(&rawOrganization))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)

		return nil
	}
}

func resourceTaikunOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.UpdateOrganizationCommand{}
	body.SetAddress(d.Get("address").(string))
	body.SetBillingEmail(d.Get("billing_email").(string))
	body.SetCity(d.Get("city").(string))
	body.SetCountry(d.Get("country").(string))
	body.SetDiscountRate(d.Get("discount_rate").(float64))
	body.SetEmail(d.Get("email").(string))
	body.SetFullName(d.Get("full_name").(string))
	body.SetId(id)
	body.SetIsEligibleUpdateSubscription(d.Get("managers_can_change_subscription").(bool))
	body.SetIsLocked(d.Get("lock").(bool))
	body.SetName(d.Get("name").(string))
	body.SetPhone(d.Get("phone").(string))
	body.SetVatNumber(d.Get("vat_number").(string))

	res, err := apiClient.Client.OrganizationsApi.OrganizationsUpdate(context.TODO()).UpdateOrganizationCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return readAfterUpdateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}

func resourceTaikunOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.OrganizationsApi.OrganizationsDelete(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunOrganization(rawOrganization *tkcore.OrganizationDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"address":                          rawOrganization.GetAddress(),
		"billing_email":                    rawOrganization.GetBillingEmail(),
		"city":                             rawOrganization.GetCity(),
		"country":                          rawOrganization.GetCountry(),
		"created_at":                       rawOrganization.GetCreatedAt(),
		"discount_rate":                    rawOrganization.GetDiscountRate(),
		"email":                            rawOrganization.GetEmail(),
		"full_name":                        rawOrganization.GetFullName(),
		"id":                               i32toa(rawOrganization.GetId()),
		"managers_can_change_subscription": rawOrganization.GetIsEligibleUpdateSubscription(),
		"lock":                             rawOrganization.GetIsLocked(),
		"is_read_only":                     rawOrganization.GetIsReadOnly(),
		"name":                             rawOrganization.GetName(),
		"partner_id":                       i32toa(rawOrganization.GetPartnerId()),
		"partner_name":                     rawOrganization.GetPartnerName(),
		"phone":                            rawOrganization.GetPhone(),
		"vat_number":                       rawOrganization.GetVatNumber(),
	}
}
