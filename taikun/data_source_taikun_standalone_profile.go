package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunStandaloneProfileSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunStandaloneProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunStandaloneProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a standalone profile by its ID.",
		ReadContext: dataSourceTaikunStandaloneProfileRead,
		Schema:      dataSourceTaikunStandaloneProfileSchema(),
	}
}

func dataSourceTaikunStandaloneProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunStandaloneProfileReadWithoutRetries()(ctx, d, meta)
}
