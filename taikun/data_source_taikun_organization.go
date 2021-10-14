package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/organizations"
)

func dataSourceTaikunOrganizationSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunOrganizationSchema())
	addOptionalFieldsToSchema(dsSchema, "id")
	return dsSchema
}

func dataSourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Organization details",
		ReadContext: dataSourceTaikunOrganizationRead,
		Schema:      dataSourceTaikunOrganizationSchema(),
	}
}

func dataSourceTaikunOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := organizations.NewOrganizationsListParams().WithV(ApiVersion)

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
