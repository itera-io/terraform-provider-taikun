package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
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

func dataSourceTaikunAccessProfilesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion)
	organizationID, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		params = params.WithOrganizationID(organizationID)
	}

	accessProfilesList := []models.AccessProfilesListDto{}
	for {
		response, err := apiClient.client.AccessProfiles.AccessProfilesList(params, &apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		accessProfilesList = append(accessProfilesList, response.GetPayload().Data...)
		if len(accessProfilesList) == payload.TotalCount {
			break
		}
		params = params.WithOffset(len(accessProfilesList))
	}

	if err := data.Set("access_profiles", accessProfilesList); err != nil {
		return diag.FromErr(err)
	}

	if organizationIDProvided {
		data.SetId(organizationID)
	} else {
		data.SetId(-1)
	}

	return nil
}
