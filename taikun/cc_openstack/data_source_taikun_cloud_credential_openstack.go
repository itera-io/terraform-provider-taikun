package cc_openstack

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialOpenStackSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialOpenStackSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "password", "url")
	return dsSchema
}

func DataSourceTaikunCloudCredentialOpenStack() *schema.Resource {
	return &schema.Resource{
		Description: "Get an OpenStack cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialOpenStackRead,
		Schema:      dataSourceTaikunCloudCredentialOpenStackSchema(),
	}
}

func dataSourceTaikunCloudCredentialOpenStackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries()(ctx, d, meta)
}
