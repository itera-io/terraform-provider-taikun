package cc_proxmox

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunCloudCredentialsProxmox() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Proxmox cloud credentials.",
		ReadContext: dataSourceTaikunCloudCredentialsProxmoxRead,
		Schema: map[string]*schema.Schema{
			"cloud_credentials": {
				Description: "List of retrieved Proxmox cloud credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialProxmoxSchema(),
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

func dataSourceTaikunCloudCredentialsProxmoxRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var cloudCredentialsList []tkcore.ProxmoxListDto

	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		cloudCredentialsList = append(cloudCredentialsList, response.GetData()...)
		if len(cloudCredentialsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(cloudCredentialsList))
	}

	cloudCredentials := make([]map[string]interface{}, len(cloudCredentialsList))
	for i, rawCloudCredential := range cloudCredentialsList {
		cloudCredentials[i] = flattenTaikunCloudCredentialProxmox(&rawCloudCredential)
	}
	if err := d.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
