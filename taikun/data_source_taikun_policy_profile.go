package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunPolicyProfileSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunPolicyProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunPolicyProfile() *schema.Resource {
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
