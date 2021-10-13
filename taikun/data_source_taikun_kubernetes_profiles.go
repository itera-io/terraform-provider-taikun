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
		Description: "Get the list of Kubernetes profiles, optionally filtered by organization.",
		ReadContext: dataSourceTaikunKubernetesProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"kubernetes_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bastion_proxy_enabled": {
							Description: "Exposes the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"created_by": {
							Description: "The creator of the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cni": {
							Description: "Container Network Interface(CNI) of the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The id of the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"is_locked": {
							Description: "Indicates whether the Kubernetes profile is locked or not.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"last_modified": {
							Description: "Time of last modification.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified_by": {
							Description: "The last user who modified the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"load_balancing_solution": {
							Description: "Load-balancing solution: `None`, `Octavia` or `Taikun`. `Octavia` and `Taikun` are only available for OpenStack cloud.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_id": {
							Description: "The id of the organization which owns the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_name": {
							Description: "The name of the organization which owns the Kubernetes profile.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
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

	kubernetesProfiles := make([]map[string]interface{}, len(kubernetesProfilesListDtos), len(kubernetesProfilesListDtos))
	for i, rawKubernetesProfile := range kubernetesProfilesListDtos {
		kubernetesProfiles[i] = flattenDataSourceTaikunKubernetesProfilesItem(rawKubernetesProfile)
	}
	if err := data.Set("kubernetes_profiles", kubernetesProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunKubernetesProfilesItem(rawKubernetesProfile *models.KubernetesProfilesListDto) map[string]interface{} {

	return map[string]interface{}{
		"bastion_proxy_enabled":   rawKubernetesProfile.ExposeNodePortOnBastion,
		"created_by":              rawKubernetesProfile.CreatedBy,
		"cni":                     rawKubernetesProfile.Cni,
		"id":                      i32toa(rawKubernetesProfile.ID),
		"is_locked":               rawKubernetesProfile.IsLocked,
		"last_modified":           rawKubernetesProfile.LastModified,
		"last_modified_by":        rawKubernetesProfile.LastModifiedBy,
		"load_balancing_solution": getLoadBalancingSolution(rawKubernetesProfile.OctaviaEnabled, rawKubernetesProfile.TaikunLBEnabled),
		"name":                    rawKubernetesProfile.Name,
		"organization_id":         i32toa(rawKubernetesProfile.OrganizationID),
		"organization_name":       rawKubernetesProfile.OrganizationName,
	}
}
