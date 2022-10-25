package taikun

import (
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/google_cloud"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunCloudCredentialGCPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"az_count": {
			Description:  "The number of GCP availability zone expected for the region.",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
			Default:      1,
		},
		"billing_account_id": {
			Description:   "The ID of the GCP credential's billing account.",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"import_project"},
		},
		"billing_account_name": {
			Description: "The name of the GCP credential's billing account.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"config_file": {
			Description:      "The path of the GCP credential's configuration file.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsFilePath,
		},
		"folder_id": {
			Description:   "The folder ID of the GCP credential.",
			Optional:      true,
			Type:          schema.TypeString,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"import_project"},
		},
		"id": {
			Description: "The ID of the GCP credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the GCP cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"import_project": {
			Description:   "Whether to import a project or not",
			Type:          schema.TypeBool,
			Default:       false,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"billing_account_id", "folder_id"},
		},
		"lock": {
			Description: "Indicates whether to lock the GCP cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description: "The name of the GCP credential.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or '-'",
				),
			),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the GCP credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the GCP credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"region": {
			Description:  "The region of the GCP credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"zones": {
			Description: "The given zones of the GCP credential.",
			Type:        schema.TypeSet,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func resourceTaikunCloudCredentialGCP() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Google Cloud Platform Credential",
		CreateContext: resourceTaikunCloudCredentialGCPCreate,
		ReadContext:   generateResourceTaikunCloudCredentialGCPReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialGCPUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialGCPSchema(),
	}
}

func resourceTaikunCloudCredentialGCPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	params := google_cloud.NewGoogleCloudCreateParams().WithV(ApiVersion)

	configFile, err := os.Open(d.Get("config_file").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	params = params.WithConfig(configFile)

	name := d.Get("name").(string)
	params = params.WithName(&name)
	region := d.Get("region").(string)
	params = params.WithRegion(&region)

	azCount := int32(d.Get("az_count").(int))
	/*
		if err != nil {
			return diag.FromErr(err)
		} else if azCount < 1 || azCount > 3 {
			return diag.Errorf("The az_count value must be between 1 and 3 inclusive.")
		}
	*/
	params = params.WithAzCount(&azCount)

	importProject := d.Get("import_project").(bool)
	params = params.WithImportProject(&importProject)
	if !importProject {
		billingAccountID := d.Get("billing_account_id").(string)
		params = params.WithBillingAccountID(&billingAccountID)
		folderID := d.Get("folder_id").(string)
		params = params.WithFolderID(&folderID)
	}

	organizationID, err := getOrganizationFromDataOrElseDefault(d, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	params = params.WithOrganizationID(&organizationID)

	createResult, err := apiClient.Client.GoogleCloud.GoogleCloudCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialGCPLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunCloudCredentialGCPReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunCloudCredentialGCPReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGCPRead(true)
}
func generateResourceTaikunCloudCredentialGCPReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGCPRead(false)
}

func generateResourceTaikunCloudCredentialGCPRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Google) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialGCP := response.GetPayload().Google[0]

		err = setResourceDataFromMap(d, flattenTaikunCloudCredentialGCP(rawCloudCredentialGCP))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialGCPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("lock") {
		if err := resourceTaikunCloudCredentialGCPLock(id, d.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunCloudCredentialGCPReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialGCP(rawGCPCredential *models.GoogleCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"billing_account_id":   rawGCPCredential.BillingAccountID,
		"billing_account_name": rawGCPCredential.BillingAccountName,
		"folder_id":            rawGCPCredential.FolderID,
		"id":                   i32toa(rawGCPCredential.ID),
		"is_default":           rawGCPCredential.IsDefault,
		"lock":                 rawGCPCredential.IsLocked,
		"name":                 rawGCPCredential.Name,
		"organization_id":      i32toa(rawGCPCredential.OrganizationID),
		"organization_name":    rawGCPCredential.OrganizationName,
		"region":               rawGCPCredential.Region,
		"zones":                rawGCPCredential.Zones,
	}
}

func resourceTaikunCloudCredentialGCPLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.CloudLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.Client.CloudCredentials.CloudCredentialsLockManager(params, apiClient)

	return err
}
