package taikun

import (
	"context"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunCloudCredentialsOpenStack() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of OpenStack cloud credentials, optionally filtered by organization.",
		ReadContext: dataSourceTaikunCloudCredentialsOpenStackRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization id filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"cloud_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialOpenStackSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunCloudCredentialsOpenStackRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var cloudCredentialsList []*models.OpenstackCredentialsListDto
	for {
		response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudCredentialsList = append(cloudCredentialsList, response.GetPayload().Openstack...)
		if len(cloudCredentialsList) == int(response.GetPayload().TotalCountOpenstack) {
			break
		}
		offset := int32(len(cloudCredentialsList))
		params = params.WithOffset(&offset)
	}

	cloudCredentials := make([]map[string]interface{}, len(cloudCredentialsList))
	for i, rawCloudCredential := range cloudCredentialsList {
		cloudCredentials[i] = flattenDataSourceTaikunCloudCredentialOpenStackItem(rawCloudCredential)
	}
	if err := data.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunCloudCredentialOpenStackItem(rawOpenStackCredential *models.OpenstackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":                 rawOpenStackCredential.CreatedBy,
		"id":                         i32toa(rawOpenStackCredential.ID),
		"is_locked":                  rawOpenStackCredential.IsLocked,
		"is_default":                 rawOpenStackCredential.IsDefault,
		"last_modified":              rawOpenStackCredential.LastModified,
		"last_modified_by":           rawOpenStackCredential.LastModifiedBy,
		"name":                       rawOpenStackCredential.Name,
		"user":                       rawOpenStackCredential.User,
		"project_name":               rawOpenStackCredential.Project,
		"project_id":                 rawOpenStackCredential.TenantID,
		"organization_id":            i32toa(rawOpenStackCredential.OrganizationID),
		"organization_name":          rawOpenStackCredential.OrganizationName,
		"public_network_name":        rawOpenStackCredential.PublicNetwork,
		"availability_zone":          rawOpenStackCredential.AvailabilityZone,
		"domain":                     rawOpenStackCredential.Domain,
		"region":                     rawOpenStackCredential.Region,
		"volume_type_name":           rawOpenStackCredential.VolumeType,
		"imported_network_subnet_id": rawOpenStackCredential.InternalSubnetID,
	}
}
