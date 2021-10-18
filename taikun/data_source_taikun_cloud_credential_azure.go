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
		Description: "Get an Azure cloud credential by its id.",
		ReadContext: dataSourceTaikunCloudCredentialAzureRead,
		Schema:      dataSourceTaikunCloudCredentialAzureSchema(),
	}
}

func dataSourceTaikunCloudCredentialAzureRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunCloudCredentialAzureRead(ctx, data, meta)
}