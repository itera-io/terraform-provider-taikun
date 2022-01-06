package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingCredentialSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunBillingCredentialSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing credential by its ID.",
		ReadContext: dataSourceTaikunBillingCredentialRead,
		Schema:      dataSourceTaikunBillingCredentialSchema(),
	}
}

func dataSourceTaikunBillingCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunBillingCredentialReadWithoutRetries()(ctx, d, meta)
}
