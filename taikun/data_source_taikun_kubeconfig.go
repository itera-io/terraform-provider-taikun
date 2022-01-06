package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubeconfigSchema() map[string]*schema.Schema {
	kubeconfigSchema := dataSourceSchemaFromResourceSchema(resourceTaikunKubeconfigSchema())
	addRequiredFieldsToSchema(kubeconfigSchema, "id")
	setValidateDiagFuncToSchema(kubeconfigSchema, "id", stringIsInt)
	deleteFieldsFromSchema(kubeconfigSchema, "role")
	return kubeconfigSchema
}

func dataSourceTaikunKubeconfig() *schema.Resource {
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
