package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunCloudCredentialGCPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"Zone": {
			Description:  "The zone of the GCP credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
	// FIXME
	return nil
}

func generateResourceTaikunCloudCredentialGCPReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGCPRead(true)
}
func generateResourceTaikunCloudCredentialGCPReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGCPRead(false)
}

func generateResourceTaikunCloudCredentialGCPRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// FIXME
		return nil
	}
}

func resourceTaikunCloudCredentialGCPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}
