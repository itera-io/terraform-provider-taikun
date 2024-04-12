package showback

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunShowbackRuleSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunShowbackRule() *schema.Resource {
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
