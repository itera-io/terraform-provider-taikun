package project

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunProjects() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all projects.",
		ReadContext: dataSourceTaikunProjectsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
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

func dataSourceTaikunProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	params := apiClient.Client.ProjectsAPI.ProjectsList(ctx)

	if organizationIDData, organizationIDProvided := d.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	response, res, err := params.Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	projects := make([]map[string]interface{}, len(response.GetData()))
	for i, projectEntityDTO := range response.GetData() {
		response, res, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectEntityDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		responseVM, res, err := apiClient.Client.StandaloneAPI.StandaloneDetails(ctx, projectEntityDTO.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(ctx, projectEntityDTO.GetId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		boundImageDTOs, err := resourceTaikunProjectGetBoundImageDTOs(ctx, projectEntityDTO.GetId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		project := response.GetProject()
		quotaResponse, res, err := apiClient.Client.ProjectQuotasAPI.ProjectquotasList(ctx).Id(project.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(quotaResponse.GetData()) != 1 {
			return nil
		}

		deleteOnExpiration, err := resourceTaikunProjectGetDeleteOnExpiration(ctx, projectEntityDTO.GetId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		responseProject := response.GetProject()
		projects[i] = flattenTaikunProject(&responseProject, response.GetData(), responseVM.GetData(), boundFlavorDTOs, boundImageDTOs, &quotaResponse.GetData()[0], deleteOnExpiration)
	}
	if err := d.Set("projects", projects); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)
	return nil
}
