package virtual_cluster

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunVirtualClusterSchema() map[string]*schema.Schema {
	projectSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunVirtualClusterSchema())
	utils.AddRequiredFieldsToSchema(projectSchema, "id")
	utils.SetValidateDiagFuncToSchema(projectSchema, "id", utils.StringIsInt)
	return projectSchema
}

func DataSourceTaikunVirtualCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a Virtual Project by its ID.",
		ReadContext: dataSourceTaikunVirtualClusterRead,
		Schema:      dataSourceTaikunVirtualClusterSchema(),
	}
}

func dataSourceTaikunVirtualClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get and set ID (ID is specified by user)
	virtualClusterId, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.Errorf("cannot read Virtual Project ID: %s", err)
	}
	d.SetId(utils.I32toa(virtualClusterId))

	// Get and set parent ID (Parent ID must be queried form API)
	apiClient := meta.(*tk.Client)
	data, response, err := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(virtualClusterId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}
	if data.GetTotalCount() != 1 {
		return diag.Errorf("There should be one, but we found %d virtual projects with specified ID.", data.TotalCount)
	}
	err = d.Set("parent_id", utils.I32toa(data.GetData()[0].GetParentProjectId()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Use normal function to read and flatten the virtual cluster.
	return generateResourceTaikunVirtualClusterRead()(ctx, d, meta)
}
