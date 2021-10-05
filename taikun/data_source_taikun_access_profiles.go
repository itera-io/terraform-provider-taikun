package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of access profiles, optionally filtered by organization",
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"access_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"http_proxy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_locked": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modified_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ntp_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"organization_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"organization_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"projects": {
							Description: "List of associated projects",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
