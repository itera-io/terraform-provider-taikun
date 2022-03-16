package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialGoogleSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialGoogleSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunCloudCredentialGoogle() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Google cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialGoogleRead,
		Schema:      dataSourceTaikunCloudCredentialGoogleSchema(),
	}
}

func dataSourceTaikunCloudCredentialGoogleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialGoogleReadWithoutRetries()(ctx, d, meta)
}
