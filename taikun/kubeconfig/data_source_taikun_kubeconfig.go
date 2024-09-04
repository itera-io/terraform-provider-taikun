package kubeconfig

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubeconfigSchema() map[string]*schema.Schema {
	kubeconfigSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunKubeconfigSchema())
	utils.AddRequiredFieldsToSchema(kubeconfigSchema, "id")
	utils.SetValidateDiagFuncToSchema(kubeconfigSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(kubeconfigSchema, "role")
	return kubeconfigSchema
}

func DataSourceTaikunKubeconfig() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a kubeconfig by its ID.",
		ReadContext: dataSourceTaikunKubeconfigRead,
		Schema:      dataSourceTaikunKubeconfigSchema(),
	}
}

func dataSourceTaikunKubeconfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))
	return generateResourceTaikunKubeconfigReadWithoutRetries()(ctx, d, meta)
}
