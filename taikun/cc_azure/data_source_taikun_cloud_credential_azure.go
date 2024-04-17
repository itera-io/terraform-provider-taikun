package cc_azure

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialAzureSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialAzureSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "subscription_id", "client_id", "client_secret")
	return dsSchema
}

func DataSourceTaikunCloudCredentialAzure() *schema.Resource {
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
