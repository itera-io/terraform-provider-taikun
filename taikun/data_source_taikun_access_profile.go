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

func dataSourceTaikunAccessProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return generateResourceTaikunAccessProfileRead(false)(ctx, data, meta)
}
