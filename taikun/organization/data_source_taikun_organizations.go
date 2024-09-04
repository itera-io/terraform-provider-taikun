package organization

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunOrganizations() *schema.Resource {
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

func dataSourceTaikunOrganizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO())

	var rawOrganizationsList []tkcore.OrganizationDetailsDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		rawOrganizationsList = append(rawOrganizationsList, response.Data...)
		if len(rawOrganizationsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(rawOrganizationsList))
	}

	organizationsList := make([]map[string]interface{}, len(rawOrganizationsList))
	for i, rawOrganization := range rawOrganizationsList {
		organizationsList[i] = flattenTaikunOrganization(&rawOrganization)
		organizationsList[i]["cloud_credentials"] = rawOrganization.GetCloudCredentials()
		organizationsList[i]["projects"] = rawOrganization.GetProjects()
		organizationsList[i]["servers"] = rawOrganization.GetServers()
		organizationsList[i]["users"] = rawOrganization.GetUsers()
	}
	if err := d.Set("organizations", organizationsList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
