package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunAlertingProfileSchema() map[string]*schema.Schema {
	alertingProfileSchema := dataSourceSchemaFromResourceSchema(resourceTaikunAlertingProfileSchema())
	addRequiredFieldsToSchema(alertingProfileSchema, "id")
	setValidateDiagFuncToSchema(alertingProfileSchema, "id", stringIsInt)
	return alertingProfileSchema
}

func dataSourceTaikunAlertingProfile() *schema.Resource {
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
