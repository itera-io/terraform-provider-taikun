package project

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunProjectSchema() map[string]*schema.Schema {
	projectSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunProjectSchema())
	utils.AddRequiredFieldsToSchema(projectSchema, "id")
	utils.SetValidateDiagFuncToSchema(projectSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(projectSchema, "taikun_lb_flavor", "router_id_start_range", "router_id_end_range")
	return projectSchema
}

func DataSourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a project by its ID.",
		ReadContext: dataSourceTaikunProjectRead,
		Schema:      dataSourceTaikunProjectSchema(),
	}
}

func dataSourceTaikunProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))
	return generateResourceTaikunProjectReadWithoutRetries()(ctx, d, meta)
}
