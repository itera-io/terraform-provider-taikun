package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/s3_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBackupCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the backup credential is the organization's default.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the backup credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the backup credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the backup credential.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the backup credential.",
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
		"s3_access_key_id": {
			Description:  "The S3 access key ID.",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
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
		"s3_secret_access_key": {
			Description:  "The S3 secret access key.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceTaikunBackupCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Backup Credential",
		CreateContext: resourceTaikunBackupCredentialCreate,
		ReadContext:   generateResourceTaikunBackupCredentialReadWithoutRetries(),
		UpdateContext: resourceTaikunBackupCredentialUpdate,
		DeleteContext: resourceTaikunBackupCredentialDelete,
		Schema:        resourceTaikunBackupCredentialSchema(),
	}
}

func resourceTaikunBackupCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.BackupCredentialsCreateCommand{
		S3Name:        d.Get("name").(string),
		S3AccessKeyID: d.Get("s3_access_key_id").(string),
		S3SecretKey:   d.Get("s3_secret_access_key").(string),
		S3Region:      d.Get("s3_region").(string),
		S3Endpoint:    d.Get("s3_endpoint").(string),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := s3_credentials.NewS3CredentialsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.S3Credentials.S3CredentialsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunBackupCredentialLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunBackupCredentialReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunBackupCredentialReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunBackupCredentialRead(true)
}
func generateResourceTaikunBackupCredentialReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunBackupCredentialRead(false)
}
func generateResourceTaikunBackupCredentialRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.S3Credentials.S3CredentialsList(s3_credentials.NewS3CredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBackupCredential := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunBackupCredential(rawBackupCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBackupCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunBackupCredentialLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("name", "s3_access_key_id", "s3_secret_access_key") {
		updateBody := models.BackupCredentialsUpdateCommand{
			ID:            id,
			S3AccessKeyID: d.Get("s3_access_key_id").(string),
			S3SecretKey:   d.Get("s3_secret_access_key").(string),
			S3Name:        d.Get("name").(string),
		}
		updateParams := s3_credentials.NewS3CredentialsUpdateParams().WithV(ApiVersion).WithBody(&updateBody)
		_, err = apiClient.client.S3Credentials.S3CredentialsUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunBackupCredentialLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunBackupCredentialReadWithRetries(), ctx, d, meta)
}

func resourceTaikunBackupCredentialDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := s3_credentials.NewS3CredentialsDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.S3Credentials.S3CredentialsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunBackupCredential(rawBackupCredential *models.BackupCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawBackupCredential.CreatedBy,
		"id":                i32toa(rawBackupCredential.ID),
		"lock":              rawBackupCredential.IsLocked,
		"is_default":        rawBackupCredential.IsDefault,
		"last_modified":     rawBackupCredential.LastModified,
		"last_modified_by":  rawBackupCredential.LastModifiedBy,
		"name":              rawBackupCredential.S3Name,
		"organization_id":   i32toa(rawBackupCredential.OrganizationID),
		"organization_name": rawBackupCredential.OrganizationName,
		"s3_access_key_id":  rawBackupCredential.S3AccessKeyID,
		"s3_region":         rawBackupCredential.S3Region,
		"s3_endpoint":       rawBackupCredential.S3Endpoint,
	}
}

func resourceTaikunBackupCredentialLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.BackupLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := s3_credentials.NewS3CredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.S3Credentials.S3CredentialsLockManager(params, apiClient)
	return err
}
