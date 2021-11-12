package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBackupCredentialSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunBackupCredentialSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "s3_secret_access_key")
	return dsSchema
}

func dataSourceTaikunBackupCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a backup credential by its ID.",
		ReadContext: dataSourceTaikunBackupCredentialRead,
		Schema:      dataSourceTaikunBackupCredentialSchema(),
	}
}

func dataSourceTaikunBackupCredentialRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return generateResourceTaikunBackupCredentialReadWithoutRetries()(ctx, data, meta)
}
