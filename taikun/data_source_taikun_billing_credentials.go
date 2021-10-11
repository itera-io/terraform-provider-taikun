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
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"billing_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The id of the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"prometheus_username": {
							Description: "The prometheus username.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"prometheus_password": {
							Description: "The prometheus password.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"prometheus_url": {
							Description: "The prometheus url.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_id": {
							Description: "The id of the organization which owns the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_name": {
							Description: "The name of the organization which owns the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"is_locked": {
							Description: "Indicates whether the billing credential is locked or not.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"is_default": {
							Description: "Indicates whether the billing credential is the organization's default or not.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"created_by": {
							Description: "The creator of the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified": {
							Description: "Time of last modification.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified_by": {
							Description: "The last user who modified the billing credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunBillingCredentialsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	var organizationID int32 = -1
	if organizationIDProvided {
		organizationID, err := atoi32(organizationIDData.(string))
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

	operationCredentials := make([]map[string]interface{}, len(operationCredentialsList), len(operationCredentialsList))
	for i, rawOperationCredential := range operationCredentialsList {
		operationCredentials[i] = flattenDatasourceTaikunBillingCredentialItem(rawOperationCredential)
	}
	if err := data.Set("billing_credentials", operationCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(organizationID))

	return nil
}

func flattenDatasourceTaikunBillingCredentialItem(rawOperationCredential *models.OperationCredentialsListDto) map[string]interface{} {

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
