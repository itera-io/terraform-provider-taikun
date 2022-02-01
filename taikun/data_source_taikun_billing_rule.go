package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingRuleSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunBillingRuleSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunBillingRule() *schema.Resource {
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
