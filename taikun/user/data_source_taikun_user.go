package user

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunUserSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunUserSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsUUID)
	return dsSchema
}

func DataSourceTaikunUser() *schema.Resource {
	return &schema.Resource{
		Description: "Get a user by its ID.",
		ReadContext: dataSourceTaikunUserRead,
		Schema:      dataSourceTaikunUserSchema(),
	}
}

func dataSourceTaikunUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunUserReadWithoutRetries()(ctx, d, meta)
}
