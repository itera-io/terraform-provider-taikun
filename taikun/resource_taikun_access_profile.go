package taikun

import (
	"context"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/ssh_users"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the access profile.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the access profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"http_proxy": {
			Description:  "HTTP Proxy of the access profile.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		"ntp_server": {
			Description: "List of NTP servers.",
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    2,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"address": {
						Description: "Address of the NTP server.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"id": {
						Description: "ID of the NTP server.",
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
			MaxItems:    2,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"address": {
						Description: "Address of the DNS server.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"id": {
						Description: "ID of the DNS server.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"ssh_user": {
			Description: "List of SSH users.",
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description: "Name of the SSH user.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(3, 30),
							validation.StringMatch(
								regexp.MustCompile("^[a-z_][a-z0-9_-]*[$]?"),
								"expect a valid linux user",
							),
							validation.StringNotInSlice([]string{"ubuntu"}, true),
						),
					},
					"public_key": {
						Description:  "Public key of the SSH user.",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"id": {
						Description: "ID of the SSH user.",
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
			Description: "The time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "Indicates whether the access profile is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"project": {
			Description: "List of associated projects.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "ID of associated project.",
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
	}
}

func resourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Access Profile",
		CreateContext: resourceTaikunAccessProfileCreate,
		ReadContext:   resourceTaikunAccessProfileRead,
		UpdateContext: resourceTaikunAccessProfileUpdate,
		DeleteContext: resourceTaikunAccessProfileDelete,
		Schema:        resourceTaikunAccessProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunAccessProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.UpsertAccessProfileCommand{
		Name: data.Get("name").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	if proxy, isProxySet := data.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	if SSHUsers, isSSHUsersSet := data.GetOk("ssh_user"); isSSHUsersSet {
		rawSSHUsersList := SSHUsers.([]interface{})
		SSHUsersList := make([]*models.SSHUserCreateDto, len(rawSSHUsersList))
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
		NTPServersList := make([]*models.NtpServerListDto, len(rawNtpServersList))
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
		DNSServersList := make([]*models.DNSServerListDto, len(rawDNSServersList))
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

	data.SetId(createResult.Payload.ID)

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

	return resourceTaikunAccessProfileRead(ctx, data, meta)
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
	if len(response.Payload.Data) != 1 {
		return diag.Errorf("access profile with ID %d not found", id)
	}

	rawAccessProfile := response.GetPayload().Data[0]

	DNSServers := make([]map[string]interface{}, len(rawAccessProfile.DNSServers))
	for i, rawDNSServer := range rawAccessProfile.DNSServers {
		DNSServers[i] = map[string]interface{}{
			"address": rawDNSServer.Address,
			"id":      i32toa(rawDNSServer.ID),
		}
	}

	NTPServers := make([]map[string]interface{}, len(rawAccessProfile.NtpServers))
	for i, rawNTPServer := range rawAccessProfile.NtpServers {
		NTPServers[i] = map[string]interface{}{
			"address": rawNTPServer.Address,
			"id":      i32toa(rawNTPServer.ID),
		}
	}

	projects := make([]map[string]interface{}, len(rawAccessProfile.Projects))
	for i, rawProject := range rawAccessProfile.Projects {
		projects[i] = map[string]interface{}{
			"id":   i32toa(rawProject.ID),
			"name": rawProject.Name,
		}
	}

	SSHUsers := make([]map[string]interface{}, len(sshResponse.Payload))
	for i, rawSSHUser := range sshResponse.Payload {
		SSHUsers[i] = map[string]interface{}{
			"id":         i32toa(rawSSHUser.ID),
			"name":       rawSSHUser.Name,
			"public_key": rawSSHUser.SSHPublicKey,
		}
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
	if err := data.Set("id", i32toa(rawAccessProfile.ID)); err != nil {
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
	if err := data.Set("organization_id", i32toa(rawAccessProfile.OrganizationID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_name", rawAccessProfile.OrganizationName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("project", projects); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("ssh_user", SSHUsers); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(id))

	return nil
}

func resourceTaikunAccessProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())

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

	if !data.HasChangeExcept("is_locked") {
		return resourceTaikunAccessProfileRead(ctx, data, meta)
	}

	body := &models.UpsertAccessProfileCommand{
		ID:   id,
		Name: data.Get("name").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	if proxy, isProxySet := data.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	updateResponse, err := apiClient.client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(updateResponse.Payload.ID)

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
