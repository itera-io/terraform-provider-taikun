package catalog_project_binding

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"log"
)

func DataSourceTaikunCatalogProjectBindings() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all projects bound to catalog.",
		ReadContext: dataSourceTaikunCatalogProjectBindingsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
			"catalog_name": {
				Description:      "Name of the catalog.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.StringLenBetween(1, 64),
			},
			"catalog_project_bindings": {
				Description: "List of retrieved projects bound to catalog.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCatalogProjectBindingSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunCatalogProjectBindingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	catalogName := d.Get("catalog_name").(string)
	var orgId int32
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationIDDataConverted, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		orgId = organizationIDDataConverted
	}

	rawCatalog, err := findCatalogByName(apiClient, orgId, catalogName)
	if err != nil {
		return diag.FromErr(err)
	}

	boundProjects := make([]map[string]interface{}, len(rawCatalog.BoundProjects))

	for i, value := range rawCatalog.BoundProjects {
		log.Printf("bound project nr %d: %s", i, utils.I32toa(value.GetId()))
		boundProjects[i] = flattenTaikunCatalogProjectBinding(true, utils.I32toa(value.GetId()), rawCatalog.GetName())
	}

	d.SetId(catalogName + "_" + utils.I32toa(orgId))

	if err := d.Set("catalog_name", d.Get("catalog_name")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", d.Get("organization_id")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("catalog_project_bindings", boundProjects); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
