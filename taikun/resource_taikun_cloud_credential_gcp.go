package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTaikunCloudCredentialGoogleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func resourceTaikunCloudCredentialGoogle() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Google Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialGoogleCreate,
		ReadContext:   generateResourceTaikunCloudCredentialGoogleReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialGoogleUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialGoogleSchema(),
	}
}

func resourceTaikunCloudCredentialGoogleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}

func generateResourceTaikunCloudCredentialGoogleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGoogleRead(true)
}
func generateResourceTaikunCloudCredentialGoogleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialGoogleRead(false)
}

func generateResourceTaikunCloudCredentialGoogleRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		// FIXME
		return nil
	}
}

func resourceTaikunCloudCredentialGoogleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}
