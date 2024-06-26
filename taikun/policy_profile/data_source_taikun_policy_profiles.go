package policy_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunPolicyProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Policy profiles.",
		ReadContext: dataSourceTaikunPolicyProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
			"policy_profiles": {
				Description: "List of retrieved Policy profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunPolicyProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunPolicyProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	var offset int32 = 0
	params := apiClient.Client.OpaProfilesAPI.OpaprofilesList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var opaProfilesListDtos []tkcore.OpaProfileListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		opaProfilesListDtos = append(opaProfilesListDtos, response.GetData()...)
		if len(opaProfilesListDtos) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(opaProfilesListDtos))
	}

	opaProfiles := make([]map[string]interface{}, len(opaProfilesListDtos))
	for i, rawOPAProfile := range opaProfilesListDtos {
		opaProfiles[i] = flattenTaikunPolicyProfile(&rawOPAProfile)
	}
	if err := d.Set("policy_profiles", opaProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
