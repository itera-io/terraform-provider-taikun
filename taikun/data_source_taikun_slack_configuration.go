package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunSlackConfigurationSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunSlackConfigurationSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunSlackConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Get a slack configuration by its ID.",
		ReadContext: dataSourceTaikunSlackConfigurationRead,
		Schema:      dataSourceTaikunSlackConfigurationSchema(),
	}
}

func dataSourceTaikunSlackConfigurationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))
	return generateResourceTaikunSlackConfigurationRead(false)(ctx, data, meta)
}
