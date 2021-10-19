package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/s3_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBackupCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The id of the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the backup credential.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"s3_access_key_id": {
			Description:  "The S3 Access Key ID.",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"s3_secret_access_key": {
			Description:  "The S3 Secret Access Key.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"s3_endpoint": {
			Description:  "The S3 endpoint URL.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		"s3_region": {
			Description: "The S3 region.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"organization_id": {
			Description:      "The id of the organization which owns the backup credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "Indicates whether the backup credential is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the backup credential is the organization's default or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunBackupCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Backup Credential",
		CreateContext: resourceTaikunBackupCredentialCreate,
		ReadContext:   resourceTaikunBackupCredentialRead,
		UpdateContext: resourceTaikunBackupCredentialUpdate,
		DeleteContext: resourceTaikunBackupCredentialDelete,
		Schema:        resourceTaikunBackupCredentialSchema(),
	}
}

func resourceTaikunBackupCredentialCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.BackupCredentialsCreateCommand{
		S3Name:        data.Get("name").(string),
		S3AccessKeyID: data.Get("s3_access_key_id").(string),
		S3SecretKey:   data.Get("s3_secret_access_key").(string),
		S3Region:      data.Get("s3_region").(string),
		S3Endpoint:    data.Get("s3_endpoint").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := s3_credentials.NewS3CredentialsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.S3Credentials.S3CredentialsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	locked := data.Get("is_locked").(bool)
	if locked {
		id, err := atoi32(createResult.Payload.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		lockBody := models.BackupLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := s3_credentials.NewS3CredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.S3Credentials.S3CredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunBackupCredentialRead(ctx, data, meta)
}

func resourceTaikunBackupCredentialRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.S3Credentials.S3CredentialsList(s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(response.Payload.Data) == 1 {
		rawBackupCredential := response.GetPayload().Data[0]

		if err := data.Set("created_by", rawBackupCredential.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawBackupCredential.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawBackupCredential.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_default", rawBackupCredential.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawBackupCredential.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawBackupCredential.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawBackupCredential.S3Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawBackupCredential.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawBackupCredential.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("s3_endpoint", rawBackupCredential.S3Endpoint); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("s3_region", rawBackupCredential.S3Region); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("s3_access_key_id", rawBackupCredential.S3AccessKeyID); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunBackupCredentialUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("name") || data.HasChange("s3_access_key_id") || data.HasChange("s3_secret_access_key") {
		updateBody := models.BackupCredentialsUpdateCommand{
			ID:            id,
			S3AccessKeyID: data.Get("s3_access_key_id").(string),
			S3SecretKey:   data.Get("s3_secret_access_key").(string),
			S3Name:        data.Get("name").(string),
		}
		updateParams := s3_credentials.NewS3CredentialsUpdateParams().WithV(ApiVersion).WithBody(&updateBody)
		_, err = apiClient.client.S3Credentials.S3CredentialsUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("is_locked") {
		lockBody := models.BackupLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := s3_credentials.NewS3CredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.S3Credentials.S3CredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunBackupCredentialRead(ctx, data, meta)
}

func resourceTaikunBackupCredentialDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := s3_credentials.NewS3CredentialsDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.S3Credentials.S3CredentialsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
