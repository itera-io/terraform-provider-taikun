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

	//var projectAppDTOs []tkcore.InstanceAppListDto
	//for {
	//	response, res, err := params.Offset(offset).Execute()
	//	if err != nil {
	//		return diag.FromErr(tk.CreateError(res, err))
	//	}
	//	projectAppDTOs = append(projectAppDTOs, response.Data...)
	//	if len(projectAppDTOs) == int(response.GetTotalCount()) {
	//		break
	//	}
	//	offset = int32(len(projectAppDTOs))
	//}
	//
	//appInstances := make([]map[string]interface{}, len(projectAppDTOs))
	//for i, appInstanceDTO := range projectAppDTOs {
	//
	//	// Load all the found data to the local object
	//	paramsInFile := paramsSpecifiedAsFile(d)
	//	if paramsInFile {
	//		// File parameters were used
	//		err := utils.SetResourceDataFromMap(d, flattenTaikunAppInstance(true, appInstanceDTO))
	//		if err != nil {
	//			return diag.FromErr(err)
	//		}
	//	} else {
	//		// Base64 parameters were used
	//		err := utils.SetResourceDataFromMap(d, flattenTaikunAppInstance(false, appInstanceDTO))
	//		if err != nil {
	//			return diag.FromErr(err)
	//		}
	//	}
	//
	//	appInstances[i] = flattenTaikunAppInstance(&alertingProfileDTO, alertingIntegrationsResponse)
	//}
	//
	//if err := d.Set("app_instances", appInstances); err != nil {
	//	return diag.FromErr(err)
	//}

	//appInstances := make([]map[string]interface{}, len(appInstancesList))
	//for i, rawAppInstance := range appInstancesList {
	//	appInstances[i] = flattenTaikunAppInstanceList(&rawAppInstance)
	//	// Get params for this app
	//	data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappDetails(context.TODO(), rawAppInstance.GetId()).Execute()
	//	if err != nil {
	//		return diag.FromErr(tk.CreateError(response, err))
	//	}
	//	appInstances[i]["parameters_yaml"] = b64.URLEncoding.EncodeToString([]byte(data.GetValues()))
	//}
	//if err := d.Set("application_instances", appInstances); err != nil {
	//	return diag.FromErr(err)
	//}

	d.SetId(dataSourceID)

	return nil
}

func flattenTaikunAppInstanceList(rawAppInstance *tkcore.ProjectAppDetailsDto) map[string]interface{} {
	var taikun_link_enabled = false
	if rawAppInstance.GetTaikunLinkUrl() != "" {
		taikun_link_enabled = true
	}
	return map[string]interface{}{
		"id":                utils.I32toa(rawAppInstance.GetId()),
		"name":              rawAppInstance.GetName(),
		"namespace":         rawAppInstance.GetNamespace(),
		"project_id":        utils.I32toa(rawAppInstance.GetProjectId()),
		"catalog_app_id":    utils.I32toa(rawAppInstance.GetCatalogAppId()),
		"parameters_base64": b64.URLEncoding.EncodeToString([]byte(rawAppInstance.GetValues())),
		"autosync":          rawAppInstance.GetAutoSync(),
		"taikun_link":       taikun_link_enabled,
		"taikun_link_url":   rawAppInstance.GetTaikunLinkUrl(),
	}
}
