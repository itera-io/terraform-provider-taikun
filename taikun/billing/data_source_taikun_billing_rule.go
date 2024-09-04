package billing

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingRuleSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunBillingRuleSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunBillingRule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing rule by its ID.",
		ReadContext: dataSourceTaikunBillingRuleRead,
		Schema:      dataSourceTaikunBillingRuleSchema(),
	}
}

func dataSourceTaikunBillingRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunBillingRuleReadWithoutRetries()(ctx, d, meta)
}
