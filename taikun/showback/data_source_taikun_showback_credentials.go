package showback

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkshowback "github.com/itera-io/taikungoclient/showbackclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunShowbackCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all showback credentials.",
		ReadContext: dataSourceTaikunShowbackCredentialsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
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
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	params := apiClient.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsList(context.TODO())
	var offset int32 = 0

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var showbackCredentialsList []tkshowback.ShowbackCredentialsListDto
	for {
		response, resp, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(resp, err))
		}
		showbackCredentialsList = append(showbackCredentialsList, response.GetData()...)
		if len(showbackCredentialsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(showbackCredentialsList))
	}

	showbackCredentials := make([]map[string]interface{}, len(showbackCredentialsList))
	for i, rawShowbackCredential := range showbackCredentialsList {
		showbackCredentials[i] = flattenTaikunShowbackCredential(&rawShowbackCredential)
	}
	if err := d.Set("showback_credentials", showbackCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
