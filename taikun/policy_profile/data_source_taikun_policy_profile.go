package policy_profile

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunPolicyProfileSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunPolicyProfileSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunPolicyProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Policy profile by its ID.",
		ReadContext: dataSourceTaikunPolicyProfileRead,
		Schema:      dataSourceTaikunPolicyProfileSchema(),
	}
}

func dataSourceTaikunPolicyProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunPolicyProfileReadWithoutRetries()(ctx, d, meta)
}
