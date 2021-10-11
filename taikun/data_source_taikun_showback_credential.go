package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing credential by its id.",
		ReadContext: dataSourceTaikunShowbackCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "The id of the showback credential.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
			},
			"name": {
				Description: "The name of the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"prometheus_username": {
				Description: "The prometheus username.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"prometheus_password": {
				Description: "The prometheus password.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"prometheus_url": {
				Description: "The prometheus url.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_id": {
				Description: "The id of the organization which owns the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_name": {
				Description: "The name of the organization which owns the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_locked": {
				Description: "Indicates whether the showback credential is locked or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_by": {
				Description: "The creator of the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified": {
				Description: "Time of last modification.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The last user who modified the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceTaikunShowbackCredentialRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunShowbackCredentialRead(ctx, data, meta)
}
