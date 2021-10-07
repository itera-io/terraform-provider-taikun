package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/organizations"
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_read_only": {
				Type:     schema.TypeString,
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

	// TODO use Organizations/Details endpoint instead
	var limit int32 = 1
	response, err := apiClient.client.Organizations.OrganizationsList(organizations.NewOrganizationsListParams().WithV(ApiVersion).WithSearchID(&id).WithLimit(&limit), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.GetPayload().TotalCount != 1 {
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
	if err := data.Set("partner_id", rawOrganization.PartnerID); err != nil {
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
	return nil
}

func resourceTaikunOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTaikunOrganizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
