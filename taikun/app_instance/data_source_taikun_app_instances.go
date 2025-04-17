package app_instance

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func DataSourceTaikunAppInstances() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all application instances.",
		ReadContext: dataSourceTaikunAppInstancesRead,
		Schema: map[string]*schema.Schema{
			"application_instances": {
				Description: "List of retrieved application instances.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAppInstanceSchema(),
				},
			},
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
		},
	}
}

func dataSourceTaikunAppInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.ProjectAppsAPI.ProjectappList(context.TODO())

	if organizationIDData, organizationIDProvided := d.GetOk("organization_id"); organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var appInstancesList []tkcore.InstanceAppListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		appInstancesList = append(appInstancesList, response.Data...)
		if len(appInstancesList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(appInstancesList))
	}

	appInstances := make([]map[string]interface{}, len(appInstancesList))
	for i, rawAppInstance := range appInstancesList {

		data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappDetails(ctx, rawAppInstance.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}

		// Base64 parameters were used
		//appInstances[i] = flattenTaikunAppInstanceList(data)

		appInstances[i] = flattenTaikunAppInstance(false, data)
	}
	if err := d.Set("application_instances", appInstances); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
