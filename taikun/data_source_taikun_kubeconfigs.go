package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/kube_config"
	"github.com/itera-io/taikungoclient/models"
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
	apiClient := meta.(*taikungoclient.Client)

	projectID, err := atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	params := kube_config.NewKubeConfigListParams().WithV(ApiVersion).WithProjectID(&projectID)

	var kubeconfigDTOs []*models.KubeConfigForUserDto
	retrievedKubeconfigCount := 0
	for {
		response, err := apiClient.Client.KubeConfig.KubeConfigList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		retrievedKubeconfigCount += len(response.Payload.Data)
		for _, kubeconfigDTO := range response.Payload.Data {
			if kubeconfigDTO.ProjectID == projectID {
				kubeconfigDTOs = append(kubeconfigDTOs, kubeconfigDTO)
			}
		}
		if retrievedKubeconfigCount == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(retrievedKubeconfigCount)
		params = params.WithOffset(&offset)
	}

	kubeconfigs := make([]map[string]interface{}, len(kubeconfigDTOs))
	for i, kubeconfigDTO := range kubeconfigDTOs {
		kubeconfigContent := resourceTaikunKubeconfigGetContent(
			kubeconfigDTO.ProjectID,
			kubeconfigDTO.ID,
			apiClient,
		)
		kubeconfigs[i] = flattenTaikunKubeconfig(kubeconfigDTO, kubeconfigContent)
	}

	if err := d.Set("kubeconfigs", kubeconfigs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(projectID))
	return nil
}
