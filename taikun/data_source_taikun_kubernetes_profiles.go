package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunKubernetesProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Kubernetes profiles.",
		ReadContext: dataSourceTaikunKubernetesProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"kubernetes_profiles": {
				Description: "List of retrieved Kubernetes profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunKubernetesProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunKubernetesProfilesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var kubernetesProfilesListDtos []*models.KubernetesProfilesListDto
	for {
		response, err := apiClient.client.KubernetesProfiles.KubernetesProfilesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		kubernetesProfilesListDtos = append(kubernetesProfilesListDtos, response.GetPayload().Data...)
		if len(kubernetesProfilesListDtos) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(kubernetesProfilesListDtos))
		params = params.WithOffset(&offset)
	}

	kubernetesProfiles := make([]map[string]interface{}, len(kubernetesProfilesListDtos))
	for i, rawKubernetesProfile := range kubernetesProfilesListDtos {
		kubernetesProfiles[i] = flattenTaikunKubernetesProfile(rawKubernetesProfile)
	}
	if err := data.Set("kubernetes_profiles", kubernetesProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
