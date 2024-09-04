package cc_gcp

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialGCPSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialGCPSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)

	// config_file & import_project only make sense when declaring a resource
	utils.DeleteFieldsFromSchema(dsSchema, "config_file")
	utils.DeleteFieldsFromSchema(dsSchema, "import_project")

	return dsSchema
}

func DataSourceTaikunCloudCredentialGCP() *schema.Resource {
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
