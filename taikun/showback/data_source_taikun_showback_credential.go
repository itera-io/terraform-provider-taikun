package showback

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackCredentialSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunShowbackCredentialSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a showback credential by its ID.",
		ReadContext: dataSourceTaikunShowbackCredentialRead,
		Schema:      dataSourceTaikunShowbackCredentialSchema(),
	}
}

func dataSourceTaikunShowbackCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunShowbackCredentialReadWithoutRetries()(ctx, d, meta)
}
