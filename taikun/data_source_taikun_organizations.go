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
					Schema: dataSourceTaikunOrganizationSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunOrganizationsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := organizations.NewOrganizationsListParams().WithV(ApiVersion)

	var rawOrganizationsList []*models.OrganizationDetailsDto
	for {
		response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		rawOrganizationsList = append(rawOrganizationsList, response.Payload.Data...)
		if len(rawOrganizationsList) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(rawOrganizationsList))
		params = params.WithOffset(&offset)
	}

	organizationsList := make([]map[string]interface{}, len(rawOrganizationsList))
	for i, rawOrganization := range rawOrganizationsList {
		organizationsList[i] = flattenTaikunOrganization(rawOrganization)
		organizationsList[i]["cloud_credentials"] = rawOrganization.CloudCredentials
		organizationsList[i]["projects"] = rawOrganization.Projects
		organizationsList[i]["servers"] = rawOrganization.Servers
		organizationsList[i]["users"] = rawOrganization.Users
	}
	if err := data.Set("organizations", organizationsList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
