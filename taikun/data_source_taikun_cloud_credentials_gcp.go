package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialsGCP() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Google Cloud Platform credentials.",
		ReadContext: dataSourceTaikunCloudCredentialsGCPRead,
		Schema: map[string]*schema.Schema{
			"cloud_credentials": {
				Description: "List of retrieved Google Cloud Platform credentials.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunCloudCredentialGCPSchema(),
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

func dataSourceTaikunCloudCredentialsGCPRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// FIXME
	return nil
}
