package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubeconfigs() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a project's kubeconfigs.",
		ReadContext: dataSourceTaikunKubeconfigsRead,
		Schema: map[string]*schema.Schema{
			"kubeconfigs": {
				Description: "List of retrieved kubeconfigs.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunKubeconfigSchema(),
				},
			},
			"project_id": {
				Description:      "Project ID filter.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
		},
	}
}

func dataSourceTaikunKubeconfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	var offset int32 = 0

	projectID, err := atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	params := apiClient.Client.KubeConfigAPI.KubeconfigList(context.TODO()).ProjectId(projectID)

	var kubeconfigDTOs []tkcore.KubeConfigForUserDto
	retrievedKubeconfigCount := 0
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		retrievedKubeconfigCount += len(response.Data)
		for _, kubeconfigDTO := range response.Data {
			if kubeconfigDTO.GetProjectId() == projectID {
				kubeconfigDTOs = append(kubeconfigDTOs, kubeconfigDTO)
			}
		}
		if retrievedKubeconfigCount == int(response.GetTotalCount()) {
			break
		}
		offset = int32(retrievedKubeconfigCount)
	}

	kubeconfigs := make([]map[string]interface{}, len(kubeconfigDTOs))
	for i, kubeconfigDTO := range kubeconfigDTOs {
		kubeconfigContent := resourceTaikunKubeconfigGetContent(
			kubeconfigDTO.GetProjectId(),
			kubeconfigDTO.GetId(),
			apiClient,
		)
		kubeconfigs[i] = flattenTaikunKubeconfig(&kubeconfigDTO, kubeconfigContent)
	}

	if err := d.Set("kubeconfigs", kubeconfigs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(projectID))
	return nil
}
