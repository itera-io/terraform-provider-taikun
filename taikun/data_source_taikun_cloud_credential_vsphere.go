package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialVsphereSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialVsphereSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "password")
	return dsSchema
}

func dataSourceTaikunCloudCredentialVsphere() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Vsphere cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialVsphereRead,
		Schema:      dataSourceTaikunCloudCredentialVsphereSchema(),
	}
}

func dataSourceTaikunCloudCredentialVsphereRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialVsphereReadWithoutRetries()(ctx, d, meta)
}
