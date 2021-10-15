package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/organizations"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunOrganizations() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all organizations (only valid for Partner and Admin roles)",
		ReadContext: dataSourceTaikunOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceSchemaFromResourceSchema(resourceTaikunOrganizationSchema()),
				},
			},
		},
	}
}

func dataSourceTaikunOrganizationsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := organizations.NewOrganizationsListParams().WithV(ApiVersion)

	var organizationsList []*models.OrganizationDetailsDto
	for {
		response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		organizationsList = append(organizationsList, response.Payload.Data...)
		if len(organizationsList) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(organizationsList))
		params = params.WithOffset(&offset)
	}

	organizations := make([]map[string]interface{}, len(organizationsList))
	for i, rawOrganization := range organizationsList {
		organizations[i] = flattenDataSourceTaikunOrganizationsItem(rawOrganization)
	}
	if err := data.Set("organizations", organizations); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunOrganizationsItem(rawOrganization *models.OrganizationDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"address":                          rawOrganization.Address,
		"billing_email":                    rawOrganization.BillingEmail,
		"city":                             rawOrganization.City,
		"cloud_credentials":                rawOrganization.CloudCredentials,
		"country":                          rawOrganization.Country,
		"created_at":                       rawOrganization.CreatedAt,
		"discount_rate":                    rawOrganization.DiscountRate,
		"email":                            rawOrganization.Email,
		"full_name":                        rawOrganization.FullName,
		"id":                               i32toa(rawOrganization.ID),
		"let_managers_change_subscription": rawOrganization.IsEligibleUpdateSubscription,
		"is_locked":                        rawOrganization.IsLocked,
		"is_read_only":                     rawOrganization.IsReadOnly,
		"name":                             rawOrganization.Name,
		"partner_id":                       i32toa(rawOrganization.PartnerID),
		"partner_name":                     rawOrganization.PartnerName,
		"phone":                            rawOrganization.Phone,
		"projects":                         rawOrganization.Projects,
		"servers":                          rawOrganization.Servers,
		"users":                            rawOrganization.Users,
		"vat_number":                       rawOrganization.VatNumber,
	}
}
