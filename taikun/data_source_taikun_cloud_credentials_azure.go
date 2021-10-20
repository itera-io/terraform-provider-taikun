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
		cloudCredentials[i] = flattenDataSourceTaikunCloudCredentialAzureItem(rawCloudCredential)
	}
	if err := data.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunCloudCredentialAzureItem(rawAzureCredential *models.AzureCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawAzureCredential.CreatedBy,
		"id":                i32toa(rawAzureCredential.ID),
		"is_locked":         rawAzureCredential.IsLocked,
		"is_default":        rawAzureCredential.IsDefault,
		"last_modified":     rawAzureCredential.LastModified,
		"last_modified_by":  rawAzureCredential.LastModifiedBy,
		"name":              rawAzureCredential.Name,
		"organization_id":   i32toa(rawAzureCredential.OrganizationID),
		"organization_name": rawAzureCredential.OrganizationName,
		"availability_zone": rawAzureCredential.AvailabilityZone,
		"location":          rawAzureCredential.Location,
		"tenant_id":         rawAzureCredential.TenantID,
	}
}
