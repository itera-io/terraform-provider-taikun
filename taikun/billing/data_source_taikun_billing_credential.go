package billing

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingCredentialSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunBillingCredentialSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing credential by its ID.",
		ReadContext: dataSourceTaikunBillingCredentialRead,
		Schema:      dataSourceTaikunBillingCredentialSchema(),
	}
}

func dataSourceTaikunBillingCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunBillingCredentialReadWithoutRetries()(ctx, d, meta)
}
