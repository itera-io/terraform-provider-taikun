package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description:  "The id of the access profile.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: stringIsInt,
		},
		"created_by": {
			Description: "The creator of the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"dns_server": {
			Description: "List of DNS servers.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"address": {
						Description: "Address of DNS server.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"id": {
						Description: "Id of DNS server.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"http_proxy": {
			Description: "HTTP Proxy of the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "Indicates whether the access profile is locked or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"ntp_server": {
			Description: "List of NTP servers.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"address": {
						Description: "Address of NTP server.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"id": {
						Description: "Id of NTP server.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"organization_id": {
			Description: "The id of the organization which owns the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"project": {
			Description: "List of associated projects.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "Id of associated project.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "Name of associated project.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"ssh_user": {
			Description: "List of SSH Users.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description: "Name of SSH User.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"public_key": {
						Description: "Public key of SSH User.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"id": {
						Description: "Id of SSH User.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}

func dataSourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get an access profiles by its id.",
		ReadContext: dataSourceTaikunAccessProfileRead,
		Schema:      dataSourceTaikunAccessProfileSchema(),
	}
}

func dataSourceTaikunAccessProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunAccessProfileRead(ctx, data, meta)
}
