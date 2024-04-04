package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	apiClient := meta.(*tk.Client)

	params := apiClient.Client.GoogleAPI.GooglecloudCreate(context.TODO())

	configFile, err := os.Open(d.Get("config_file").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	params = params.Config(configFile)

	name := d.Get("name").(string)
	params = params.Name(name)
	region := d.Get("region").(string)
	params = params.Region(region)

	azCount := int32(d.Get("az_count").(int))
	/*
		if err != nil {
			return diag.FromErr(err)
		} else if azCount < 1 || azCount > 3 {
			return diag.Errorf("The az_count value must be between 1 and 3 inclusive.")
		}
	*/
	params = params.AzCount(azCount)

	importProject := d.Get("import_project").(bool)
	params = params.ImportProject(importProject)
	if !importProject {
		billingAccountID := d.Get("billing_account_id").(string)
		params = params.BillingAccountId(billingAccountID)
		folderID := d.Get("folder_id").(string)
		params = params.FolderId(folderID)
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, newErr := atoi32(organizationIDData.(string))
		if newErr != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		params = params.OrganizationId(organizationId)
	}

	createResult, res, err := params.Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

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
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, _, err := apiClient.Client.GoogleAPI.GooglecloudList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialGCP := response.GetData()[0]

		err = setResourceDataFromMap(d, flattenTaikunCloudCredentialGCP(&rawCloudCredentialGCP))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialGCPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
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

func flattenTaikunCloudCredentialGCP(rawGCPCredential *tkcore.GoogleCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"billing_account_id":   rawGCPCredential.GetBillingAccountId(),
		"billing_account_name": rawGCPCredential.GetBillingAccountName(),
		"folder_id":            rawGCPCredential.GetFolderId(),
		"id":                   i32toa(rawGCPCredential.GetId()),
		"is_default":           rawGCPCredential.GetIsDefault(),
		"lock":                 rawGCPCredential.GetIsLocked(),
		"name":                 rawGCPCredential.GetName(),
		"organization_id":      i32toa(rawGCPCredential.GetOrganizationId()),
		"organization_name":    rawGCPCredential.GetOrganizationName(),
		"region":               rawGCPCredential.GetRegion(),
		"zones":                rawGCPCredential.GetZones(),
	}
}

func resourceTaikunCloudCredentialGCPLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	_, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return err
}
