package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialsAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all AWS cloud credentials.",
		ReadContext: dataSourceTaikunCloudCredentialsAWSRead,
		Schema: map[string]*schema.Schema{
			"cloud_credentials": {
				Description: "List of retrieved AWS cloud credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialAWSSchema(),
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

func dataSourceTaikunCloudCredentialsAWSRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var cloudCredentialsList []tkcore.AmazonCredentialsListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		cloudCredentialsList = append(cloudCredentialsList, response.GetAmazon()...)
		if len(cloudCredentialsList) == int(response.GetTotalCountAws()) {
			break
		}
		offset = int32(len(cloudCredentialsList))
	}

	cloudCredentials := make([]map[string]interface{}, len(cloudCredentialsList))
	for i, rawCloudCredential := range cloudCredentialsList {
		cloudCredentials[i] = flattenTaikunCloudCredentialAWS(&rawCloudCredential)
	}
	if err := d.Set("cloud_credentials", cloudCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
