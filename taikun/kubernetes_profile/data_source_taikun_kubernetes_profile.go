package kubernetes_profile

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubernetesProfileSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(ResourceTaikunKubernetesProfileSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Kubernetes profile by its ID.",
		ReadContext: dataSourceTaikunKubernetesProfileRead,
		Schema:      dataSourceTaikunKubernetesProfileSchema(),
	}
}

func dataSourceTaikunKubernetesProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return GenerateResourceTaikunKubernetesProfileReadWithoutRetries()(ctx, d, meta)
}
