package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunShowbackRuleSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunShowbackRule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a showback rule by its id.",
		ReadContext: dataSourceTaikunShowbackRuleRead,
		Schema:      dataSourceTaikunShowbackRuleSchema(),
	}
}

func dataSourceTaikunShowbackRuleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunShowbackRuleRead(ctx, data, meta)
}
