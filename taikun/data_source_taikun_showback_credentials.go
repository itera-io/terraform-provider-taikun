package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
	"github.com/itera-io/taikungoclient/showbackclient/showback_credentials"
)

func dataSourceTaikunShowbackCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all showback credentials.",
		ReadContext: dataSourceTaikunShowbackCredentialsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"showback_credentials": {
				Description: "List of retrieved showback credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunShowbackCredentialSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunShowbackCredentialsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := showback_credentials.NewShowbackCredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var showbackCredentialsList []*models.ShowbackCredentialsListDto
	for {
		response, err := apiClient.ShowbackClient.Showback.ShowbackCredentialsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		showbackCredentialsList = append(showbackCredentialsList, response.GetPayload().Data...)
		if len(showbackCredentialsList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(showbackCredentialsList))
		params = params.WithOffset(&offset)
	}

	showbackCredentials := make([]map[string]interface{}, len(showbackCredentialsList))
	for i, rawShowbackCredential := range showbackCredentialsList {
		showbackCredentials[i] = flattenTaikunShowbackCredential(rawShowbackCredential)
	}
	if err := d.Set("showback_credentials", showbackCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
