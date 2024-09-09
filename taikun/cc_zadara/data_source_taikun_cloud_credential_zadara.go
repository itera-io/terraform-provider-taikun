package cc_zadara

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialZadaraSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialZadaraSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "secret_access_key", "access_key_id")
	return dsSchema
}

func DataSourceTaikunCloudCredentialZadara() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Zadara cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialZadaraRead,
		Schema:      dataSourceTaikunCloudCredentialZadaraSchema(),
	}
}

func dataSourceTaikunCloudCredentialZadaraRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialZadaraReadWithoutRetries()(ctx, d, meta)
}
