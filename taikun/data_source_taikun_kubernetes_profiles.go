package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubernetesProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Kubernetes profiles.",
		ReadContext: dataSourceTaikunKubernetesProfilesRead,
		Schema: map[string]*schema.Schema{
			"kubernetes_profiles": {
				Description: "List of retrieved Kubernetes profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunKubernetesProfileSchema(),
				},
			},
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
		},
	}
}

func dataSourceTaikunKubernetesProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var kubernetesProfilesListDtos []tkcore.KubernetesProfilesListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		kubernetesProfilesListDtos = append(kubernetesProfilesListDtos, response.GetData()...)
		if len(kubernetesProfilesListDtos) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(kubernetesProfilesListDtos))
	}

	kubernetesProfiles := make([]map[string]interface{}, len(kubernetesProfilesListDtos))
	for i, rawKubernetesProfile := range kubernetesProfilesListDtos {
		kubernetesProfiles[i] = flattenTaikunKubernetesProfile(&rawKubernetesProfile)
	}
	if err := d.Set("kubernetes_profiles", kubernetesProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
