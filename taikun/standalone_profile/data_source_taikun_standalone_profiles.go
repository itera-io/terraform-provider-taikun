package standalone_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunStandaloneProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all standalone profiles.",
		ReadContext: dataSourceTaikunStandaloneProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
			"standalone_profiles": {
				Description: "List of retrieved standalone profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunStandaloneProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunStandaloneProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.StandaloneProfileAPI.StandaloneprofileList(ctx)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var standaloneProfilesListDtos []tkcore.StandAloneProfilesListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		standaloneProfilesListDtos = append(standaloneProfilesListDtos, response.GetData()...)
		if len(standaloneProfilesListDtos) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(standaloneProfilesListDtos))
	}

	standaloneProfiles := make([]map[string]interface{}, len(standaloneProfilesListDtos))
	for i, rawStandaloneProfile := range standaloneProfilesListDtos {

		securityGroupResponse, res, err := apiClient.Client.SecurityGroupAPI.SecuritygroupList(ctx, rawStandaloneProfile.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		standaloneProfiles[i] = flattenTaikunStandaloneProfile(&rawStandaloneProfile, securityGroupResponse)
	}
	if err := d.Set("standalone_profiles", standaloneProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
