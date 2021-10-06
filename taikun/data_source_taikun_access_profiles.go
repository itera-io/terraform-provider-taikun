package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of access profiles, optionally filtered by organization",
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:     schema.TypeString,
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
										Type:     schema.TypeString,
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
							Type:     schema.TypeString,
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
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"organization_id": {
							Type:     schema.TypeString,
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
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
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

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	var organizationID int32 = -1
	if organizationIDProvided {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	accessProfilesList := []*models.AccessProfilesListDto{}
	for {
		response, err := apiClient.client.AccessProfiles.AccessProfilesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		accessProfilesList = append(accessProfilesList, response.GetPayload().Data...)
		if len(accessProfilesList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(accessProfilesList))
		params = params.WithOffset(&offset)
	}

	flattenedAccessProfilesList := flattenDatasourceTaikunAccessProfilesList(accessProfilesList)
	if err := data.Set("access_profiles", flattenedAccessProfilesList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(organizationID))

	return nil
}

func flattenDatasourceTaikunAccessProfilesList(rawAccessProfiles []*models.AccessProfilesListDto) []map[string]interface{} {
	accessProfiles := make([]map[string]interface{}, 0, len(rawAccessProfiles))
	for _, rawAccessProfile := range rawAccessProfiles {

		DNSServers := make([]map[string]interface{}, 0, len(rawAccessProfile.DNSServers))
		for _, rawDNSServer := range rawAccessProfile.DNSServers {
			DNSServers = append(DNSServers, map[string]interface{}{
				"address": rawDNSServer.Address,
				"id":      i32toa(rawDNSServer.ID),
			})
		}

		NTPServers := make([]map[string]interface{}, 0, len(rawAccessProfile.NtpServers))
		for _, rawNTPServer := range rawAccessProfile.NtpServers {
			NTPServers = append(NTPServers, map[string]interface{}{
				"address": rawNTPServer.Address,
				"id":      i32toa(rawNTPServer.ID),
			})
		}

		projects := make([]map[string]interface{}, 0, len(rawAccessProfile.Projects))
		for _, rawProject := range rawAccessProfile.Projects {
			projects = append(projects, map[string]interface{}{
				"id":   i32toa(rawProject.ID),
				"name": rawProject.Name,
			})
		}

		accessProfiles = append(accessProfiles, map[string]interface{}{
			"created_by":        rawAccessProfile.CreatedBy,
			"dns_servers":       DNSServers,
			"http_proxy":        rawAccessProfile.HTTPProxy,
			"id":                i32toa(rawAccessProfile.ID),
			"is_locked":         rawAccessProfile.IsLocked,
			"last_modified":     rawAccessProfile.LastModified,
			"last_modified_by":  rawAccessProfile.LastModifiedBy,
			"name":              rawAccessProfile.Name,
			"ntp_servers":       NTPServers,
			"organization_id":   i32toa(rawAccessProfile.OrganizationID),
			"organization_name": rawAccessProfile.OrganizationName,
			"projects":          projects,
		})
	}
	return accessProfiles
}
