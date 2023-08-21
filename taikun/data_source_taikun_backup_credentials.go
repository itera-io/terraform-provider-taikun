package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBackupCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all backup credentials.",
		ReadContext: dataSourceTaikunBackupCredentialsRead,
		Schema: map[string]*schema.Schema{
			"backup_credentials": {
				Description: "List of retrieved backup credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunBackupCredentialSchema(),
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

func dataSourceTaikunBackupCredentialsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	params := apiClient.Client.S3CredentialsApi.S3credentialsList(context.TODO())
	var offset int32 = 0

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var backupCredentialsList []tkcore.BackupCredentialsListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		backupCredentialsList = append(backupCredentialsList, response.Data...)
		if len(backupCredentialsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(backupCredentialsList))
	}

	backupCredentials := make([]map[string]interface{}, len(backupCredentialsList))
	for i, rawBackupCredential := range backupCredentialsList {
		backupCredentials[i] = flattenTaikunBackupCredential(&rawBackupCredential)
	}
	if err := d.Set("backup_credentials", backupCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
