package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
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

func dataSourceTaikunShowbackCredentialsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
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
		response, err := apiClient.client.Showback.ShowbackCredentialsList(params, apiClient)
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
		showbackCredentials[i] = flattenDatasourceTaikunShowbackCredentialItem(rawShowbackCredential)
	}
	if err := data.Set("showback_credentials", showbackCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDatasourceTaikunShowbackCredentialItem(rawShowbackCredential *models.ShowbackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawShowbackCredential.CreatedBy,
		"id":                i32toa(rawShowbackCredential.ID),
		"is_locked":         rawShowbackCredential.IsLocked,
		"last_modified":     rawShowbackCredential.LastModified,
		"last_modified_by":  rawShowbackCredential.LastModifiedBy,
		"name":              rawShowbackCredential.Name,
		"organization_id":   i32toa(rawShowbackCredential.OrganizationID),
		"organization_name": rawShowbackCredential.OrganizationName,
		"password":          rawShowbackCredential.Password,
		"url":               rawShowbackCredential.URL,
		"username":          rawShowbackCredential.Username,
	}
}
