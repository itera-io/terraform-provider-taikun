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
		Description: "Retrieve all organizations.",
		ReadContext: dataSourceTaikunOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"organizations": {
				Description: "List of retrieved organizations.",
				Type:        schema.TypeList,
				Computed:    true,
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
		organizations[i] = flattenTaikunOrganization(rawOrganization)
	}
	if err := data.Set("organizations", organizations); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
