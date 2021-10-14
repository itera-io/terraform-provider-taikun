package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubernetesProfileSchema() map[string]*schema.Schema {
	dsSchema := datasourceSchemaFromResourceSchema(resourceTaikunKubernetesProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Kubernetes profile by its id.",
		ReadContext: dataSourceTaikunKubernetesProfileRead,
		Schema:      dataSourceTaikunKubernetesProfileSchema(),
	}
}

func dataSourceTaikunKubernetesProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunKubernetesProfileRead(ctx, data, meta)
}
