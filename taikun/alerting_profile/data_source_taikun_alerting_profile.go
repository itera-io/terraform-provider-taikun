package alerting_profile

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunAlertingProfileSchema() map[string]*schema.Schema {
	alertingProfileSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunAlertingProfileSchema())
	utils.AddRequiredFieldsToSchema(alertingProfileSchema, "id")
	utils.SetValidateDiagFuncToSchema(alertingProfileSchema, "id", utils.StringIsInt)
	return alertingProfileSchema
}

func DataSourceTaikunAlertingProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve an alerting profile by its ID.",
		ReadContext: dataSourceTaikunAlertingProfileRead,
		Schema:      dataSourceTaikunAlertingProfileSchema(),
	}
}

func dataSourceTaikunAlertingProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))
	return generateResourceTaikunAlertingProfileReadWithoutRetries()(ctx, d, meta)
}
