package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunProjectSchema() map[string]*schema.Schema {
	projectSchema := dataSourceSchemaFromResourceSchema(resourceTaikunProjectSchema())
	addRequiredFieldsToSchema(projectSchema, "id")
	setValidateDiagFuncToSchema(projectSchema, "id", stringIsInt)
	return projectSchema
}

func dataSourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a project by its ID.",
		ReadContext: dataSourceTaikunProjectRead,
		Schema:      dataSourceTaikunProjectSchema(),
	}
}

func dataSourceTaikunProjectRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))
	return resourceTaikunProjectRead(ctx, data, meta)
}
