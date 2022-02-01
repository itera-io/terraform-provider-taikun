package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubernetesProfileSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunKubernetesProfileSchema())
	addRequiredFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	return dsSchema
}

func dataSourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Kubernetes profile by its ID.",
		ReadContext: dataSourceTaikunKubernetesProfileRead,
		Schema:      dataSourceTaikunKubernetesProfileSchema(),
	}
}

func dataSourceTaikunKubernetesProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunKubernetesProfileReadWithoutRetries()(ctx, d, meta)
}
