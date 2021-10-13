package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackCredentialSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunShowbackCredentialSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	return dsSchema
}

func dataSourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing credential by its id.",
		ReadContext: dataSourceTaikunShowbackCredentialRead,
		Schema:      dataSourceTaikunShowbackCredentialSchema(),
	}
}

func dataSourceTaikunShowbackCredentialRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunShowbackCredentialRead(ctx, data, meta)
}
