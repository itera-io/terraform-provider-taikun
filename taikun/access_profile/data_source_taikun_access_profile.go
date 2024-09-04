package access_profile

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunAccessProfileSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get an access profile by its ID.",
		ReadContext: dataSourceTaikunAccessProfileRead,
		Schema:      dataSourceTaikunAccessProfileSchema(),
	}
}

func dataSourceTaikunAccessProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunAccessProfileReadWithoutRetries()(ctx, d, meta)
}
