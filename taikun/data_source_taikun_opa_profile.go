package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunOPAProfileSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunOPAProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunOPAProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get an OPA profile by its ID.",
		ReadContext: dataSourceTaikunOPAProfileRead,
		Schema:      dataSourceTaikunOPAProfileSchema(),
	}
}

func dataSourceTaikunOPAProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return generateResourceTaikunOPAProfileReadWithoutRetries()(ctx, data, meta)
}
