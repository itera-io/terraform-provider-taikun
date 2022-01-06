package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunAccessProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get an access profile by its ID.",
		ReadContext: dataSourceTaikunAccessProfileRead,
		Schema:      dataSourceTaikunAccessProfileSchema(),
	}
}

func dataSourceTaikunAccessProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunAccessProfileReadWithoutRetries()(ctx, d, meta)
}
