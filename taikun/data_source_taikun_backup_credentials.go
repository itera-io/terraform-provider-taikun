package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/s3_credentials"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
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
	apiClient := meta.(*taikungoclient.Client)
	dataSourceID := "all"

	params := s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var backupCredentialsList []*models.BackupCredentialsListDto
	for {
		response, err := apiClient.Client.S3Credentials.S3CredentialsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		backupCredentialsList = append(backupCredentialsList, response.GetPayload().Data...)
		if len(backupCredentialsList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(backupCredentialsList))
		params = params.WithOffset(&offset)
	}

	backupCredentials := make([]map[string]interface{}, len(backupCredentialsList))
	for i, rawBackupCredential := range backupCredentialsList {
		backupCredentials[i] = flattenTaikunBackupCredential(rawBackupCredential)
	}
	if err := d.Set("backup_credentials", backupCredentials); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
