package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialProxmoxSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialProxmoxSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	deleteFieldsFromSchema(dsSchema, "client_secret")
	return dsSchema
}

func dataSourceTaikunCloudCredentialProxmox() *schema.Resource {
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
