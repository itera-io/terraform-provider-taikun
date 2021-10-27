package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialAWSSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialAWSSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "secret_access_key", "access_key_id")
	return dsSchema
}

func dataSourceTaikunCloudCredentialAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Get an AWS cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialAWSRead,
		Schema:      dataSourceTaikunCloudCredentialAWSSchema(),
	}
}

func dataSourceTaikunCloudCredentialAWSRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return generateResourceTaikunCloudCredentialAWSRead(false)(ctx, data, meta)
}
