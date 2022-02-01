package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackCredentialSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunShowbackCredentialSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a showback credential by its ID.",
		ReadContext: dataSourceTaikunShowbackCredentialRead,
		Schema:      dataSourceTaikunShowbackCredentialSchema(),
	}
}

func dataSourceTaikunShowbackCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunShowbackCredentialReadWithoutRetries()(ctx, d, meta)
}
