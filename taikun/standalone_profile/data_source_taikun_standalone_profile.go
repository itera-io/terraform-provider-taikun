package standalone_profile

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunStandaloneProfileSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunStandaloneProfileSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunStandaloneProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a standalone profile by its ID.",
		ReadContext: dataSourceTaikunStandaloneProfileRead,
		Schema:      dataSourceTaikunStandaloneProfileSchema(),
	}
}

func dataSourceTaikunStandaloneProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunStandaloneProfileReadWithoutRetries()(ctx, d, meta)
}
