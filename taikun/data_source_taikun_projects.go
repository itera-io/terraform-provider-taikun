package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceTaikunProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	params := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO())

	if organizationIDData, organizationIDProvided := d.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	response, _, err := params.Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	projects := make([]map[string]interface{}, len(response.GetData()))
	for i, projectEntityDTO := range response.GetData() {
		response, res, err := apiClient.Client.ServersAPI.ServersDetails(context.TODO(), projectEntityDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		responseVM, res, err := apiClient.Client.StandaloneAPI.StandaloneDetails(context.TODO(), projectEntityDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(projectEntityDTO.GetId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		boundImageDTOs, err := resourceTaikunProjectGetBoundImageDTOs(projectEntityDTO.GetId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		project := response.GetProject()
		quotaResponse, res, err := apiClient.Client.ProjectQuotasAPI.ProjectquotasList(context.TODO()).Id(project.GetProjectId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(quotaResponse.GetData()) != 1 {
			return nil
		}

		responseProject := response.GetProject()
		projects[i] = flattenTaikunProject(&responseProject, response.GetData(), responseVM.GetData(), boundFlavorDTOs, boundImageDTOs, &quotaResponse.GetData()[0])
	}
	if err := d.Set("projects", projects); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)
	return nil
}
