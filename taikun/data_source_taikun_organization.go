package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/organizations"
)

func dataSourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Organization details",
		ReadContext: dataSourceTaikunOrganizationRead,
		Schema: map[string]*schema.Schema{
			"address": {
				Description: "Address",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_email": {
				Description: "Billing email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"city": {
				Description: "City",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cloud_credentials": {
				Description: "Number of associated cloud credentials",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"country": {
				Description: "Country",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "Time and date of creation",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"discount_rate": {
				Description: "Discount rate, must be between 0 and 100 (included)",
				Type:        schema.TypeFloat,
				Computed:    true,
			},
			"email": {
				Description: "Email",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"full_name": {
				Description: "Full name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": {
				Description: "ID",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"let_managers_change_subscription": {
				Description: "Allow subscription to be changed by managers",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_locked": {
				Description: "Whether the organization is locked",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_read_only": {
				Description: "Whether the organization is in read-only mode",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"name": {
				Description: "Name",
				Type:        schema.TypeString,
				Computed:    true,
			},
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
				Computed:    true,
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
				Computed:    true,
			},
		},
	}
}

func dataSourceTaikunOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	var limit int32 = 1
	params := organizations.NewOrganizationsListParams().WithV(ApiVersion).WithLimit(&limit)

	id := data.Get("id").(string)
	id32, _ := atoi32(id)
	if id != "" {
		params = params.WithID(&id32)
	}

	data.SetId("")

	response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(response.GetPayload().Data) != 1 {
		return diag.Errorf("No organization found")
	}
	if id != "" && response.Payload.Data[0].ID != id32 {
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

	if id == "" {
		data.SetId("-1")
	} else {
		data.SetId(id)
	}

	return nil
}
