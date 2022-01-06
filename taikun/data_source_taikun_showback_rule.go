package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunShowbackRuleSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunShowbackRule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a showback rule by its ID.",
		ReadContext: dataSourceTaikunShowbackRuleRead,
		Schema:      dataSourceTaikunShowbackRuleSchema(),
	}
}

func dataSourceTaikunShowbackRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunShowbackRuleReadWithoutRetries()(ctx, d, meta)
}
