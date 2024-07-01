package app_instance

import (
	"context"
	b64 "encoding/base64"
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
		},
	}
}

func dataSourceTaikunAppInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	//params := apiClient.Client.CatalogAPI.CatalogList(context.TODO())
	params := apiClient.Client.ProjectAppsAPI.ProjectappList(context.TODO())

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
		appInstances[i] = flattenTaikunAppInstanceList(&rawAppInstance)
		// Get params for this app
		data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappDetails(context.TODO(), rawAppInstance.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
		appInstances[i]["parameters_yaml"] = b64.URLEncoding.EncodeToString([]byte(data.GetValues()))
	}
	if err := d.Set("application_instances", appInstances); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}

func flattenTaikunAppInstanceList(rawAppInstance *tkcore.InstanceAppListDto) map[string]interface{} {
	return map[string]interface{}{
		"id":             utils.I32toa(rawAppInstance.GetId()),
		"name":           rawAppInstance.GetName(),
		"namespace":      rawAppInstance.GetNamespace(),
		"project_id":     utils.I32toa(rawAppInstance.GetProjectId()),
		"catalog_app_id": utils.I32toa(rawAppInstance.GetCatalogAppId()),
		//"parameters_yaml": b64.URLEncoding.EncodeToString([]byte(rawAppInstance.GetValues())),
		"autosync": rawAppInstance.GetAutoSync(),
	}
}
