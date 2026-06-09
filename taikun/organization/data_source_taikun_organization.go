package organization

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunOrganizationSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunOrganizationSchema())
	utils.AddOptionalFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	return dsSchema
}

func DataSourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an organization.",
		ReadContext: dataSourceTaikunOrganizationRead,
		Schema:      dataSourceTaikunOrganizationSchema(),
	}
}

func dataSourceTaikunOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	var limit int32 = 1

	params := apiClient.Client.OrganizationsAPI.OrganizationsList(ctx).Limit(limit)
	id := d.Get("id").(string)
	if id != "" {
		id32, err := utils.Atoi32(id)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.Id(id32)
	}

	d.SetId("")

	response, res, err := params.Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(response.Data) != 1 {
		return diag.Errorf("organization not found")
	}

	rawOrganization := response.Data[0]

	err = utils.SetResourceDataFromMap(d, flattenTaikunOrganization(&rawOrganization))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.I32toa(rawOrganization.GetId()))

	return nil
}
