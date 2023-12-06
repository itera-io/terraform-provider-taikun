package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunOrganizationSchema() map[string]*schema.Schema {
	dsSchema := dataSourceSchemaFromResourceSchema(resourceTaikunOrganizationSchema())
	addOptionalFieldsToSchema(dsSchema, "id")
	setValidateDiagFuncToSchema(dsSchema, "id", stringIsInt)
	setFieldInSchema(dsSchema, "cloud_credentials", &schema.Schema{
		Description: "Number of associated cloud credentials.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	setFieldInSchema(dsSchema, "users", &schema.Schema{
		Description: "Number of associated users.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	setFieldInSchema(dsSchema, "projects", &schema.Schema{
		Description: "Number of associated projects.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	setFieldInSchema(dsSchema, "servers", &schema.Schema{
		Description: "Number of associated servers.",
		Type:        schema.TypeInt,
		Computed:    true,
	})
	return dsSchema
}

func dataSourceTaikunOrganization() *schema.Resource {
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
	id32, _ := atoi32(id)
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

	err = setResourceDataFromMap(d, organizationMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(rawOrganization.GetId()))

	return nil
}
