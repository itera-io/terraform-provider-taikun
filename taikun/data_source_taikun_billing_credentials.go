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
		Description: "Retrieve all billing credentials.",
		ReadContext: dataSourceTaikunBillingCredentialsRead,
		Schema: map[string]*schema.Schema{
			"billing_credentials": {
				Description: "List of retrieved billing credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunBillingCredentialSchema(),
				},
			},
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
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
		operationCredentials[i] = flattenTaikunBillingCredential(rawOperationCredential)
	}
	if err := data.Set("billing_credentials", operationCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
