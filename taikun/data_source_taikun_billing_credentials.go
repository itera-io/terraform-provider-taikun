package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/ops_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunBillingCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of billing credentials, optionally filtered by organization.",
		ReadContext: dataSourceTaikunBillingCredentialsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization id filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"billing_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunBillingCredentialSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunBillingCredentialsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var operationCredentialsList []*models.OperationCredentialsListDto
	for {
		response, err := apiClient.client.OpsCredentials.OpsCredentialsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		operationCredentialsList = append(operationCredentialsList, response.GetPayload().Data...)
		if len(operationCredentialsList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(operationCredentialsList))
		params = params.WithOffset(&offset)
	}

	operationCredentials := make([]map[string]interface{}, len(operationCredentialsList))
	for i, rawOperationCredential := range operationCredentialsList {
		operationCredentials[i] = flattenDataSourceTaikunBillingCredentialItem(rawOperationCredential)
	}
	if err := data.Set("billing_credentials", operationCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDataSourceTaikunBillingCredentialItem(rawOperationCredential *models.OperationCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":          rawOperationCredential.CreatedBy,
		"id":                  i32toa(rawOperationCredential.ID),
		"is_locked":           rawOperationCredential.IsLocked,
		"is_default":          rawOperationCredential.IsDefault,
		"last_modified":       rawOperationCredential.LastModified,
		"last_modified_by":    rawOperationCredential.LastModifiedBy,
		"name":                rawOperationCredential.Name,
		"organization_id":     i32toa(rawOperationCredential.OrganizationID),
		"organization_name":   rawOperationCredential.OrganizationName,
		"prometheus_password": rawOperationCredential.PrometheusPassword,
		"prometheus_url":      rawOperationCredential.PrometheusURL,
		"prometheus_username": rawOperationCredential.PrometheusUsername,
	}
}
