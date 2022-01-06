package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunUserSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunUserSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsUUID)
	return dsSchema
}

func dataSourceTaikunUser() *schema.Resource {
	return &schema.Resource{
		Description: "Get a user by its ID.",
		ReadContext: dataSourceTaikunUserRead,
		Schema:      dataSourceTaikunUserSchema(),
	}
}

func dataSourceTaikunUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunUserReadWithoutRetries()(ctx, d, meta)
}
