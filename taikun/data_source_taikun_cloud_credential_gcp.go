package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialGCPSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialGCPSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)

	// config_file & import_project only make sense when declaring a resource
	deleteFieldsFromSchema(dsSchema, "config_file")
	deleteFieldsFromSchema(dsSchema, "import_project")

	return dsSchema
}

func dataSourceTaikunCloudCredentialGCP() *schema.Resource {
	return &schema.Resource{
		Description: "Get a GCP credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialGCPRead,
		Schema:      dataSourceTaikunCloudCredentialGCPSchema(),
	}
}

func dataSourceTaikunCloudCredentialGCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialGCPReadWithoutRetries()(ctx, d, meta)
}
