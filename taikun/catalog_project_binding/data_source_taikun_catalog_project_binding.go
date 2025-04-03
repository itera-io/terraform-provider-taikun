package catalog_project_binding

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunCatalogProjectBindingSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCatalogProjectBindingSchema())

	utils.AddRequiredFieldsToSchema(dsSchema, "catalog_name")
	utils.SetValidateDiagFuncToSchema(dsSchema, "catalog_name", utils.StringLenBetween(3, 30))
	utils.AddRequiredFieldsToSchema(dsSchema, "project_id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "project_id", utils.StringIsInt)
	utils.AddOptionalFieldsToSchema(dsSchema, "organization_id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "organization_id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunCatalogProjectBinding() *schema.Resource {
	return &schema.Resource{
		Description: "Find out if Catalog has specific Project Bound",
		ReadContext: dataSourceTaikunCatalogProjectBindingRead,
		Schema:      dataSourceTaikunCatalogProjectBindingSchema(),
	}
}

func dataSourceTaikunCatalogProjectBindingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := d.Set("catalog_name", d.Get("catalog_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("project_id", d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("organization_id", d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return generateResourceTaikunCatalogProjectBindingReadWithoutRetries()(ctx, d, meta)
}
