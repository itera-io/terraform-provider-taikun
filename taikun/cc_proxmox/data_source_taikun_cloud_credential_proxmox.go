package cc_proxmox

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialProxmoxSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialProxmoxSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "client_secret")
	return dsSchema
}

func DataSourceTaikunCloudCredentialProxmox() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Proxmox cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialProxmoxRead,
		Schema:      dataSourceTaikunCloudCredentialProxmoxSchema(),
	}
}

func dataSourceTaikunCloudCredentialProxmoxRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialProxmoxReadWithoutRetries()(ctx, d, meta)
}
