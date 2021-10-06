package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
	"strconv"
)

func resourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Access Profile",
		CreateContext: resourceTaikunAccessProfileCreate,
		ReadContext:   resourceTaikunAccessProfileRead,
		UpdateContext: resourceTaikunAccessProfileUpdate,
		DeleteContext: resourceTaikunAccessProfileDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_proxy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ntp_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ssh_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeString,
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
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	}
}

func resourceTaikunAccessProfileRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := strconv.Atoi(data.Id())
	if err != nil {
		data.SetId("")
		return diag.FromErr(err)
	}
	id32 := int32(id)

	params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id32)

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
		if err := data.Set("id", strconv.Itoa(int(rawAccessProfile.ID))); err != nil {
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
		if err := data.Set("organization_id", strconv.Itoa(int(rawAccessProfile.OrganizationID))); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawAccessProfile.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("projects", projects); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(strconv.Itoa(id))
	} else {
		// Not Found
		data.SetId("")
	}

	return nil
}

func resourceTaikunAccessProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	organizationId, err := strconv.Atoi(data.Get("organization_id").(string))
	if err != nil {
		return diag.Errorf("OrganizationId provider isn't valid: %s", data.Get("organization_id").(string))
	}

	body := &models.UpsertAccessProfileCommand{
		Name:           data.Get("name").(string),
		OrganizationID: int32(organizationId),
	}

	if proxy, isProxySet := data.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	if SSHUsers, isSSHUsersSet := data.GetOk("ssh_users"); isSSHUsersSet {
		rawSSHUsersList := SSHUsers.([]map[string]interface{})
		SSHUsersList := make([]*models.SSHUserCreateDto, 0, len(rawSSHUsersList))
		for i, e := range rawSSHUsersList {
			SSHUsersList[i] = &models.SSHUserCreateDto{
				Name:         e["name"].(string),
				SSHPublicKey: e["public_key"].(string),
			}
		}
		body.SSHUsers = SSHUsersList
	}

	if NtpServers, isNTPServersSet := data.GetOk("ntp_servers"); isNTPServersSet {
		rawNtpServersList := NtpServers.([]map[string]interface{})
		NTPServersList := make([]*models.NtpServerListDto, 0, len(rawNtpServersList))
		for i, e := range rawNtpServersList {
			NTPServersList[i] = &models.NtpServerListDto{
				Address: e["address"].(string),
			}
		}
		body.NtpServers = NTPServersList
	}

	if DNSServers, isDNSServersSet := data.GetOk("dns_servers"); isDNSServersSet {
		rawDNSServersList := DNSServers.([]map[string]interface{})
		DNSServersList := make([]*models.DNSServerListDto, 0, len(rawDNSServersList))
		for i, e := range rawDNSServersList {
			DNSServersList[i] = &models.DNSServerListDto{
				Address: e["address"].(string),
			}
		}
		body.DNSServers = DNSServersList
	}

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	return resourceTaikunAccessProfileRead(ctx, data, meta)
}

func resourceTaikunAccessProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTaikunAccessProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := access_profiles.NewAccessProfilesDeleteParams().WithV(ApiVersion).WithID(int32(id))
	_, _, err = apiClient.client.AccessProfiles.AccessProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
