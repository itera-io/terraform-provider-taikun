package taikun

import (
	"context"

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

func dataSourceTaikunBackupCredentialsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
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
		response, err := apiClient.client.S3Credentials.S3CredentialsList(params, apiClient)
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
	if err := data.Set("backup_credentials", backupCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
