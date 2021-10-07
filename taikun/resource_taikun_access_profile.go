package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/ssh_users"
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
				Description: "The id of the access profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the access profile.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"organization_id": {
				Description:  "The id of the organization which owns the access profile.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
			},
			"organization_name": {
				Description: "The name of the organization which owns the access profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"http_proxy": {
				Description: "HTTP Proxy of the access profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ntp_server": {
				Description: "List of NTP servers.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description: "Address of NTP server.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Id of NTP server.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"dns_server": {
				Description: "List of DNS servers.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Description: "Address of DNS server.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Id of DNS server.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"ssh_user": {
				Description: "List of SSH Users.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of SSH User.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"public_key": {
							Description: "Public key of SSH User.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Id of SSH User.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"created_by": {
				Description: "The creator of the access profile.",
				Type:        schema.TypeString,
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
			"is_locked": {
				Description: "Indicates whether the access profile is locked or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"projects": {
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
		},
	}
}

func resourceTaikunAccessProfileRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.AccessProfiles.AccessProfilesList(access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	sshResponse, err := apiClient.client.SSHUsers.SSHUsersList(ssh_users.NewSSHUsersListParams().WithV(ApiVersion).WithAccessProfileID(id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawAccessProfile := response.GetPayload().Data[0]

		DNSServers := make([]map[string]interface{}, 0, len(rawAccessProfile.DNSServers))
		for _, rawDNSServer := range rawAccessProfile.DNSServers {
			DNSServers = append(DNSServers, map[string]interface{}{
				"address": rawDNSServer.Address,
				"id":      strconv.Itoa(int(rawDNSServer.ID)),
			})
		}

		NTPServers := make([]map[string]interface{}, 0, len(rawAccessProfile.NtpServers))
		for _, rawNTPServer := range rawAccessProfile.NtpServers {
			NTPServers = append(NTPServers, map[string]interface{}{
				"address": rawNTPServer.Address,
				"id":      strconv.Itoa(int(rawNTPServer.ID)),
			})
		}

		projects := make([]map[string]interface{}, 0, len(rawAccessProfile.Projects))
		for _, rawProject := range rawAccessProfile.Projects {
			projects = append(projects, map[string]interface{}{
				"id":   strconv.Itoa(int(rawProject.ID)),
				"name": rawProject.Name,
			})
		}

		SSHUsers := make([]map[string]interface{}, 0, len(sshResponse.Payload))
		for _, rawSSHUser := range sshResponse.Payload {
			SSHUsers = append(SSHUsers, map[string]interface{}{
				"id":         strconv.Itoa(int(rawSSHUser.ID)),
				"name":       rawSSHUser.Name,
				"public_key": rawSSHUser.SSHPublicKey,
			})
		}

		if err := data.Set("created_by", rawAccessProfile.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("dns_server", DNSServers); err != nil {
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
		if err := data.Set("ntp_server", NTPServers); err != nil {
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
		if err := data.Set("ssh_user", SSHUsers); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunAccessProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	organizationId, err := atoi32(data.Get("organization_id").(string))
	if err != nil {
		return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
	}

	body := &models.UpsertAccessProfileCommand{
		Name:           data.Get("name").(string),
		OrganizationID: organizationId,
	}

	if proxy, isProxySet := data.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	if SSHUsers, isSSHUsersSet := data.GetOk("ssh_user"); isSSHUsersSet {
		rawSSHUsersList := SSHUsers.([]interface{})
		SSHUsersList := make([]*models.SSHUserCreateDto, len(rawSSHUsersList), len(rawSSHUsersList))
		for i, e := range rawSSHUsersList {
			rawSSHUser := e.(map[string]interface{})
			SSHUsersList[i] = &models.SSHUserCreateDto{
				Name:         rawSSHUser["name"].(string),
				SSHPublicKey: rawSSHUser["public_key"].(string),
			}
		}
		body.SSHUsers = SSHUsersList
	}

	if NtpServers, isNTPServersSet := data.GetOk("ntp_server"); isNTPServersSet {
		rawNtpServersList := NtpServers.([]interface{})
		NTPServersList := make([]*models.NtpServerListDto, len(rawNtpServersList), len(rawNtpServersList))
		for i, e := range rawNtpServersList {
			rawNtpServer := e.(map[string]interface{})
			NTPServersList[i] = &models.NtpServerListDto{
				Address: rawNtpServer["address"].(string),
			}
		}
		body.NtpServers = NTPServersList
	}

	if DNSServers, isDNSServersSet := data.GetOk("dns_server"); isDNSServersSet {
		rawDNSServersList := DNSServers.([]interface{})
		DNSServersList := make([]*models.DNSServerListDto, len(rawDNSServersList), len(rawDNSServersList))
		for i, e := range rawDNSServersList {
			rawDNSServer := e.(map[string]interface{})
			DNSServersList[i] = &models.DNSServerListDto{
				Address: rawDNSServer["address"].(string),
			}
		}
		body.DNSServers = DNSServersList
	}

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	locked := data.Get("is_locked").(bool)
	if locked {
		id, err := atoi32(createResult.Payload.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		lockBody := models.AccessProfilesLockManagementCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := access_profiles.NewAccessProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.AccessProfiles.AccessProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	data.SetId(createResult.Payload.ID)

	return resourceTaikunAccessProfileRead(ctx, data, meta)
}

func resourceTaikunAccessProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())

	organizationId, err := atoi32(data.Get("organization_id").(string))
	if err != nil {
		return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
	}

	body := &models.UpsertAccessProfileCommand{
		ID:             id,
		Name:           data.Get("name").(string),
		OrganizationID: organizationId,
	}

	if proxy, isProxySet := data.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("is_locked") {
		lockBody := models.AccessProfilesLockManagementCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := access_profiles.NewAccessProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.AccessProfiles.AccessProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	data.SetId(createResult.Payload.ID)

	return resourceTaikunAccessProfileRead(ctx, data, meta)
}

func resourceTaikunAccessProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := access_profiles.NewAccessProfilesDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.AccessProfiles.AccessProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
