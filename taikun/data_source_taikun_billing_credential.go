package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingCredentialSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunBillingCredentialSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	return dsSchema
}

func dataSourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing credential by its id.",
		ReadContext: dataSourceTaikunBillingCredentialRead,
		Schema:      dataSourceTaikunBillingCredentialSchema(),
	}
}

func dataSourceTaikunBillingCredentialRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunBillingCredentialRead(ctx, data, meta)
}
