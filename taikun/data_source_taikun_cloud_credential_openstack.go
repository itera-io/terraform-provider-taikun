package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialOpenStackSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialOpenStackSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "password", "url")
	return dsSchema
}

func dataSourceTaikunCloudCredentialOpenStack() *schema.Resource {
	return &schema.Resource{
		Description: "Get an OpenStack cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialOpenStackRead,
		Schema:      dataSourceTaikunCloudCredentialOpenStackSchema(),
	}
}

func dataSourceTaikunCloudCredentialOpenStackRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries()(ctx, data, meta)
}
