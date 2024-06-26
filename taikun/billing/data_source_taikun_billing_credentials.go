package billing

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunBillingCredentials() *schema.Resource {
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
				ValidateDiagFunc: utils.StringIsInt,
			},
		},
	}
}

func dataSourceTaikunBillingCredentialsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.OperationCredentialsAPI.OpscredentialsList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var operationCredentialsList []tkcore.OperationCredentialsListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		operationCredentialsList = append(operationCredentialsList, response.Data...)
		if len(operationCredentialsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(operationCredentialsList))
	}

	operationCredentials := make([]map[string]interface{}, len(operationCredentialsList))
	for i, rawOperationCredential := range operationCredentialsList {
		operationCredentials[i] = flattenTaikunBillingCredential(&rawOperationCredential)
	}
	if err := d.Set("billing_credentials", operationCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
