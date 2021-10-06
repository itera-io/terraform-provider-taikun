package taikun

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
)

func dataSourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get an access profiles by its id",
		ReadContext: dataSourceTaikunAccessProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
						"id": {
							Type:     schema.TypeInt,
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
	}
}

//TODO Should be moved to the AccessProfile Resource
func dataSourceTaikunAccessProfileRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := int32(data.Get("id").(int))

	params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id)

	response, err := apiClient.client.AccessProfiles.AccessProfilesList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawAccessProfile := response.GetPayload().Data[0]

		DNSServers := make([]map[string]interface{}, 0, len(rawAccessProfile.DNSServers))
		for _, rawDNSServer := range rawAccessProfile.DNSServers {
			DNSServers = append(DNSServers, map[string]interface{}{
				"address": rawDNSServer.Address,
				"id":      rawDNSServer.ID,
			})
		}

		NTPServers := make([]map[string]interface{}, 0, len(rawAccessProfile.NtpServers))
		for _, rawNTPServer := range rawAccessProfile.NtpServers {
			NTPServers = append(NTPServers, map[string]interface{}{
				"address": rawNTPServer.Address,
				"id":      rawNTPServer.ID,
			})
		}

		projects := make([]map[string]interface{}, 0, len(rawAccessProfile.Projects))
		for _, rawProject := range rawAccessProfile.Projects {
			projects = append(projects, map[string]interface{}{
				"id":   rawProject.ID,
				"name": rawProject.Name,
			})
		}

		if err := data.Set("created_by", rawAccessProfile.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("dns_servers", DNSServers); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("http_proxy", rawAccessProfile.HTTPProxy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", rawAccessProfile.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawAccessProfile.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawAccessProfile.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawAccessProfile.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawAccessProfile.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("ntp_servers", NTPServers); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", rawAccessProfile.OrganizationID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawAccessProfile.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("projects", projects); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(strconv.Itoa(int(id)))
	} else {
		// Not Found
		data.SetId("")
	}

	return nil
}
