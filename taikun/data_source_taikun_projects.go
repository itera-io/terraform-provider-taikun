package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/project_quotas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
)

func dataSourceTaikunProjects() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all projects.",
		ReadContext: dataSourceTaikunProjectsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"projects": {
				Description: "List of retrieved projects.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunProjectSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunProjectsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := projects.NewProjectsListSelectorParams().WithV(ApiVersion)

	if organizationIDData, organizationIDProvided := data.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	response, err := apiClient.client.Projects.ProjectsListSelector(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	projects := make([]map[string]interface{}, len(response.Payload))
	for i, projectEntityDTO := range response.Payload {
		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectEntityDTO.ID)
		response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(projectEntityDTO.ID, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		quotaParams := project_quotas.NewProjectQuotasListParams().WithV(ApiVersion).WithID(&response.Payload.Project.QuotaID)
		quotaResponse, err := apiClient.client.ProjectQuotas.ProjectQuotasList(quotaParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(quotaResponse.Payload.Data) != 1 {
			return nil
		}

		projects[i] = flattenTaikunProject(response.Payload.Project, response.Payload.Data, boundFlavorDTOs, quotaResponse.Payload.Data[0])
	}
	if err := data.Set("projects", projects); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)
	return nil
}
