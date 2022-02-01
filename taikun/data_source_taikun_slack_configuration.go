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
		Description: "Get a Slack configuration by its ID.",
		ReadContext: dataSourceTaikunSlackConfigurationRead,
		Schema:      dataSourceTaikunSlackConfigurationSchema(),
	}
}

func dataSourceTaikunSlackConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))
	return generateResourceTaikunSlackConfigurationReadWithoutRetries()(ctx, d, meta)
}
