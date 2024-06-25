package catalog

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
)

func DataSourceTaikunCatalogs() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all catalogs.",
		ReadContext: dataSourceTaikunCatalogsRead,
		Schema: map[string]*schema.Schema{
			"catalogs": {
				Description: "List of retrieved catalogs.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCatalogSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunCatalogsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.CatalogAPI.CatalogList(context.TODO())

	var catalogsList []tkcore.CatalogListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		catalogsList = append(catalogsList, response.Data...)
		if len(catalogsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(catalogsList))
	}

	catalogs := make([]map[string]interface{}, len(catalogsList))
	for i, rawCatalog := range catalogsList {
		catalogs[i] = flattenTaikunCatalog(&rawCatalog)
	}
	if err := d.Set("catalogs", catalogs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
