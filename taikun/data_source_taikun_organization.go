package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/organizations"
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
	apiClient := meta.(*taikungoclient.Client)

	var limit int32 = 1
	params := organizations.NewOrganizationsListParams().WithV(ApiVersion).WithLimit(&limit)

	id := d.Get("id").(string)
	id32, _ := atoi32(id)
	if id != "" {
		params = params.WithID(&id32)
	}

	d.SetId("")

	response, err := apiClient.Client.Organizations.OrganizationsList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(response.Payload.Data) != 1 {
		return diag.Errorf("organization with ID %d not found", id32)
	}

	rawOrganization := response.GetPayload().Data[0]

	organizationMap := flattenTaikunOrganization(rawOrganization)
	organizationMap["cloud_credentials"] = rawOrganization.CloudCredentials
	organizationMap["projects"] = rawOrganization.Projects
	organizationMap["servers"] = rawOrganization.Servers
	organizationMap["users"] = rawOrganization.Users

	err = setResourceDataFromMap(d, organizationMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(rawOrganization.ID))

	return nil
}
