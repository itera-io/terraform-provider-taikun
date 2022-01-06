package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialAzureSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialAzureSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "subscription_id", "client_id", "client_secret")
	return dsSchema
}

func dataSourceTaikunCloudCredentialAzure() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Azure cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialAzureRead,
		Schema:      dataSourceTaikunCloudCredentialAzureSchema(),
	}
}

func dataSourceTaikunCloudCredentialAzureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialAzureReadWithoutRetries()(ctx, d, meta)
}
