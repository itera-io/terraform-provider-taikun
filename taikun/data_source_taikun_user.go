package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunUserSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunUserSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunUser() *schema.Resource {
	return &schema.Resource{
		Description: "Get a user by its id.",
		ReadContext: dataSourceTaikunUserRead,
		Schema:      dataSourceTaikunUserSchema(),
	}
}

func dataSourceTaikunUserRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunUserRead(ctx, data, meta)
}
