package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/cloud_credentials"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunCloudCredentialsAzure() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Azure cloud credentials.",
		ReadContext: dataSourceTaikunCloudCredentialsAzureRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"cloud_credentials": {
				Description: "List of retrieved Azure cloud credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialAzureSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunCloudCredentialsAzureRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var cloudCredentialsList []*models.AzureCredentialsListDto
	for {
		response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudCredentialsList = append(cloudCredentialsList, response.GetPayload().Azure...)
		if len(cloudCredentialsList) == int(response.GetPayload().TotalCountAzure) {
			break
		}
		offset := int32(len(cloudCredentialsList))
		params = params.WithOffset(&offset)
	}

	cloudCredentials := make([]map[string]interface{}, len(cloudCredentialsList))
	for i, rawCloudCredential := range cloudCredentialsList {
		cloudCredentials[i] = flattenTaikunCloudCredentialAzure(rawCloudCredential)
	}
	if err := data.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
