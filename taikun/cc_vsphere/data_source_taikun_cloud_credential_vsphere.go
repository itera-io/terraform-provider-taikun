package cc_vsphere

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialVsphereSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialVsphereSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "password")
	return dsSchema
}

func DataSourceTaikunCloudCredentialVsphere() *schema.Resource {
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
