package catalog

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunCatalogSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCatalogSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "name")
	utils.SetValidateDiagFuncToSchema(dsSchema, "name", utils.StringLenBetween(3, 30))
	return dsSchema
}

func DataSourceTaikunCatalog() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Catalog by its name.",
		ReadContext: dataSourceTaikunCatalogRead,
		Schema:      dataSourceTaikunCatalogSchema(),
	}
}

func dataSourceTaikunCatalogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := d.Set("name", d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	return generateResourceTaikunCatalogReadWithoutRetries()(ctx, d, meta)
}
