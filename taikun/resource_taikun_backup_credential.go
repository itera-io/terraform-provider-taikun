package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	apiClient := meta.(*tk.Client)

	body := tkcore.BackupCredentialsCreateCommand{}
	body.SetS3Name(d.Get("name").(string))
	body.SetS3AccessKeyId(d.Get("s3_access_key_id").(string))
	body.SetS3SecretKey(d.Get("s3_secret_access_key").(string))
	body.SetS3Region(d.Get("s3_region").(string))
	body.SetS3Endpoint(d.Get("s3_endpoint").(string))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.S3CredentialsAPI.S3credentialsCreate(context.TODO()).BackupCredentialsCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

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
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.S3CredentialsAPI.S3credentialsList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBackupCredential := response.Data[0]

		err = setResourceDataFromMap(d, flattenTaikunBackupCredential(&rawBackupCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBackupCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
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
		body := tkcore.BackupCredentialsUpdateCommand{}
		body.SetId(id)
		body.SetS3SecretKey(d.Get("s3_secret_access_key").(string))
		body.SetS3AccessKeyId(d.Get("s3_access_key_id").(string))
		body.SetS3Name(d.Get("name").(string))

		res, err := apiClient.Client.S3CredentialsAPI.S3credentialsUpdate(ctx).BackupCredentialsUpdateCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
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
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.S3CredentialsAPI.S3credentialsDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunBackupCredential(rawBackupCredential *tkcore.BackupCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawBackupCredential.GetCreatedBy(),
		"id":                i32toa(rawBackupCredential.GetId()),
		"lock":              rawBackupCredential.GetIsLocked(),
		"is_default":        rawBackupCredential.GetIsDefault(),
		"last_modified":     rawBackupCredential.GetLastModified(),
		"last_modified_by":  rawBackupCredential.GetLastModifiedBy(),
		"name":              rawBackupCredential.GetS3Name(),
		"organization_id":   i32toa(rawBackupCredential.GetOrganizationId()),
		"organization_name": rawBackupCredential.GetOrganizationName(),
		"s3_access_key_id":  rawBackupCredential.GetS3AccessKeyId(),
		"s3_region":         rawBackupCredential.GetS3Region(),
		"s3_endpoint":       rawBackupCredential.GetS3Endpoint(),
	}
}

func resourceTaikunBackupCredentialLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.BackupLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))
	res, err := apiClient.Client.S3CredentialsAPI.S3credentialsLockManagement(context.TODO()).BackupLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
