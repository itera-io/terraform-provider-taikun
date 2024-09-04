package slack

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunSlackConfigurationSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunSlackConfigurationSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunSlackConfiguration() *schema.Resource {
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
