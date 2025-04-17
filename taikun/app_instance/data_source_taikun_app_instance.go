package app_instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunAppInstanceSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunAppInstanceSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunAppInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Application Instance by its ID.",
		ReadContext: dataSourceTaikunAppInstanceRead,
		Schema:      dataSourceTaikunAppInstanceSchema(),
	}
}

func dataSourceTaikunAppInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))
	return generateResourceTaikunAppInstanceReadWithoutRetries()(ctx, d, meta)
}
