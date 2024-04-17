package backup_credential

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBackupCredentialSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunBackupCredentialSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "s3_secret_access_key")
	return dsSchema
}

func DataSourceTaikunBackupCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a backup credential by its ID.",
		ReadContext: dataSourceTaikunBackupCredentialRead,
		Schema:      dataSourceTaikunBackupCredentialSchema(),
	}
}

func dataSourceTaikunBackupCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunBackupCredentialReadWithoutRetries()(ctx, d, meta)
}
