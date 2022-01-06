package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/itera-io/taikungoclient/client/kube_config"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_kubeconfig", &resource.Sweeper{
		Name:         "taikun_kubeconfig",
		Dependencies: []string{},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := kube_config.NewKubeConfigListParams().WithV(ApiVersion)

			var kubeconfigDTOs []*models.KubeConfigForUserDto

			for {
				response, err := apiClient.client.KubeConfig.KubeConfigList(params, apiClient)
				if err != nil {
					return err
				}
				kubeconfigDTOs = append(kubeconfigDTOs, response.GetPayload().Data...)
				if len(kubeconfigDTOs) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(kubeconfigDTOs))
				params = params.WithOffset(&offset)
			}

			for _, e := range kubeconfigDTOs {
				if shouldSweep(e.Name) {

					body := models.DeleteKubeConfigCommand{
						ID: e.ID,
					}
					params := kube_config.NewKubeConfigDeleteParams().WithV(ApiVersion).WithBody(&body)
					if _, err := apiClient.client.KubeConfig.KubeConfigDelete(params, apiClient); err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}
