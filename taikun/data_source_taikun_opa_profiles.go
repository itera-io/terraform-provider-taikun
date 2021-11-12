package taikun

import (
	"context"
	"github.com/itera-io/taikungoclient/client/opa_profiles"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunOPAProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all OPA profiles.",
		ReadContext: dataSourceTaikunOPAProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"opa_profiles": {
				Description: "List of retrieved OPA profiles.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunOPAProfileSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunOPAProfilesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var opaProfilesListDtos []*models.OpaProfileListDto
	for {
		response, err := apiClient.client.OpaProfiles.OpaProfilesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		opaProfilesListDtos = append(opaProfilesListDtos, response.GetPayload().Data...)
		if len(opaProfilesListDtos) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(opaProfilesListDtos))
		params = params.WithOffset(&offset)
	}

	opaProfiles := make([]map[string]interface{}, len(opaProfilesListDtos))
	for i, rawOPAProfile := range opaProfilesListDtos {
		opaProfiles[i] = flattenTaikunOPAProfile(rawOPAProfile)
	}
	if err := data.Set("opa_profiles", opaProfiles); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
