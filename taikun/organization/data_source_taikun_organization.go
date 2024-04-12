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
	utils.SetFieldInSchema(dsSchema, "cloud_credentials", &schema.Schema{
		Description: "Number of associated cloud credentials.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	utils.SetFieldInSchema(dsSchema, "users", &schema.Schema{
		Description: "Number of associated users.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	utils.SetFieldInSchema(dsSchema, "projects", &schema.Schema{
		Description: "Number of associated projects.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	utils.SetFieldInSchema(dsSchema, "servers", &schema.Schema{
		Description: "Number of associated servers.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	return dsSchema
}

func DataSourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Get the details of an organization.",
		ReadContext: dataSourceTaikunOrganizationRead,
		Schema:      dataSourceTaikunOrganizationSchema(),
	}
}

func dataSourceTaikunOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	var limit int32 = 1

	params := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Limit(limit)
	id := d.Get("id").(string)
	id32, _ := utils.Atoi32(id)
	if id != "" {
		params = params.Id(id32)
	}

	d.SetId("")

	response, res, err := params.Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(response.Data) != 1 {
		return diag.Errorf("organization with ID %d not found", id32)
	}

	rawOrganization := response.Data[0]

	organizationMap := flattenTaikunOrganization(&rawOrganization)
	organizationMap["cloud_credentials"] = rawOrganization.GetCloudCredentials()
	organizationMap["projects"] = rawOrganization.GetProjects()
	organizationMap["servers"] = rawOrganization.GetServers()
	organizationMap["users"] = rawOrganization.GetUsers()

	err = utils.SetResourceDataFromMap(d, organizationMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.I32toa(rawOrganization.GetId()))

	return nil
}
