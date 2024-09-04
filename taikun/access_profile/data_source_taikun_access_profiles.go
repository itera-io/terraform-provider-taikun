package access_profile

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func DataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all access profiles.",
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema: map[string]*schema.Schema{
			"access_profiles": {
				Description: "List of retrieved access profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunAccessProfileSchema(),
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

func dataSourceTaikunAccessProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.AccessProfilesAPI.AccessprofilesList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var accessProfilesList []tkcore.AccessProfilesListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		accessProfilesList = append(accessProfilesList, response.Data...)
		if len(accessProfilesList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(accessProfilesList))
	}

	accessProfiles := make([]map[string]interface{}, len(accessProfilesList))
	for i, rawAccessProfile := range accessProfilesList {

		sshResponse, res, err := apiClient.Client.SshUsersAPI.SshusersList(context.TODO(), rawAccessProfile.GetId()).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		accessProfiles[i] = flattenTaikunAccessProfile(&rawAccessProfile, sshResponse)
	}
	if err := d.Set("access_profiles", accessProfiles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
