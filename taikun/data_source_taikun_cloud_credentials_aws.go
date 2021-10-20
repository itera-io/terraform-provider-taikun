package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/cloud_credentials"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunCloudCredentialsAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all AWS cloud credentials.",
		ReadContext: dataSourceTaikunCloudCredentialsAWSRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"cloud_credentials": {
				Description: "List of retrieved AWS cloud credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialAWSSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunCloudCredentialsAWSRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var cloudCredentialsList []*models.AmazonCredentialsListDto
	for {
		response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudCredentialsList = append(cloudCredentialsList, response.GetPayload().Amazon...)
		if len(cloudCredentialsList) == int(response.GetPayload().TotalCountAws) {
			break
		}
		offset := int32(len(cloudCredentialsList))
		params = params.WithOffset(&offset)
	}

	cloudCredentials := make([]map[string]interface{}, len(cloudCredentialsList))
	for i, rawCloudCredential := range cloudCredentialsList {
		cloudCredentials[i] = flattenDataSourceTaikunCloudCredentialAWSItem(rawCloudCredential)
	}
	if err := data.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunCloudCredentialAWSItem(rawAWSCredential *models.AmazonCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawAWSCredential.CreatedBy,
		"id":                i32toa(rawAWSCredential.ID),
		"is_locked":         rawAWSCredential.IsLocked,
		"is_default":        rawAWSCredential.IsDefault,
		"last_modified":     rawAWSCredential.LastModified,
		"last_modified_by":  rawAWSCredential.LastModifiedBy,
		"name":              rawAWSCredential.Name,
		"organization_id":   i32toa(rawAWSCredential.OrganizationID),
		"organization_name": rawAWSCredential.OrganizationName,
		"availability_zone": rawAWSCredential.AvailabilityZone,
		"region":            rawAWSCredential.Region,
	}
}
