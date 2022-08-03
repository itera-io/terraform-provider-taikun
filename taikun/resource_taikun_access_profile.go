package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/ssh_users"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the access profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"dns_server": {
			Description: "List of DNS servers.",
			Type:        schema.TypeList,
			Optional:    true,
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
		"http_proxy": {
			Description:  "HTTP proxy of the access profile.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
		"id": {
			Description: "The ID of the access profile.",
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
		"lock": {
			Description: "Indicates whether to lock the access profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the access profile.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"ntp_server": {
			Description: "List of NTP servers.",
			Type:        schema.TypeList,
			Optional:    true,
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
		"ssh_user": {
			Description: "List of SSH users.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "ID of the SSH user.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "Name of the SSH user.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(3, 30),
							validation.StringMatch(
								regexp.MustCompile("^[a-z_][a-z0-9_-]*[$]?$"),
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
				},
			},
		},
	}
}

func resourceTaikunAccessProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Access Profile",
		CreateContext: resourceTaikunAccessProfileCreate,
		ReadContext:   generateResourceTaikunAccessProfileReadWithoutRetries(),
		UpdateContext: resourceTaikunAccessProfileUpdate,
		DeleteContext: resourceTaikunAccessProfileDelete,
		Schema:        resourceTaikunAccessProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunAccessProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	body := &models.CreateAccessProfileCommand{
		Name: d.Get("name").(string),
	}
	resourceTaikunAccessProfileUpsertSetBody(d, body)

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.Client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunAccessProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunAccessProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunAccessProfileRead(true)
}
func generateResourceTaikunAccessProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunAccessProfileRead(false)
}
func generateResourceTaikunAccessProfileRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.AccessProfiles.AccessProfilesList(access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		sshResponse, err := apiClient.Client.SSHUsers.SSHUsersList(ssh_users.NewSSHUsersListParams().WithV(ApiVersion).WithAccessProfileID(id), apiClient)
		if err != nil {
			if _, ok := err.(*ssh_users.SSHUsersListNotFound); ok && withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return diag.FromErr(err)
		}

		rawAccessProfile := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunAccessProfile(rawAccessProfile, sshResponse))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunAccessProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunAccessProfileLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := resourceTaikunAccessProfileUpdateDeleteOldDNSServers(d, apiClient); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceTaikunAccessProfileUpdateDeleteOldNTPServers(d, apiClient); err != nil {
		return diag.FromErr(err)
	}
	if err := resourceTaikunAccessProfileUpdateDeleteOldSSHUsers(d, apiClient); err != nil {
		return diag.FromErr(err)
	}

	body := &models.UpsertAccessProfileCommand{
		ID:   id,
		Name: d.Get("name").(string),
	}
	resourceTaikunAccessProfileUpsertSetBody(d, body)

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	if _, err := apiClient.Client.AccessProfiles.AccessProfilesCreate(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunAccessProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunAccessProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := access_profiles.NewAccessProfilesDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.Client.AccessProfiles.AccessProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunAccessProfile(rawAccessProfile *models.AccessProfilesListDto, sshResponse *ssh_users.SSHUsersListOK) map[string]interface{} {

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

	SSHUsers := make([]map[string]interface{}, len(sshResponse.Payload))
	for i, rawSSHUser := range sshResponse.Payload {
		SSHUsers[i] = map[string]interface{}{
			"id":         i32toa(rawSSHUser.ID),
			"name":       rawSSHUser.Name,
			"public_key": rawSSHUser.SSHPublicKey,
		}
	}

	return map[string]interface{}{
		"created_by":        rawAccessProfile.CreatedBy,
		"dns_server":        DNSServers,
		"http_proxy":        rawAccessProfile.HTTPProxy,
		"id":                i32toa(rawAccessProfile.ID),
		"lock":              rawAccessProfile.IsLocked,
		"last_modified":     rawAccessProfile.LastModified,
		"last_modified_by":  rawAccessProfile.LastModifiedBy,
		"name":              rawAccessProfile.Name,
		"ntp_server":        NTPServers,
		"organization_id":   i32toa(rawAccessProfile.OrganizationID),
		"organization_name": rawAccessProfile.OrganizationName,
		"ssh_user":          SSHUsers,
	}
}

func resourceTaikunAccessProfileUpdateDeleteOldDNSServers(d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	oldDNSServersData, _ := d.GetChange("dns_server")
	oldDNSServers := oldDNSServersData.([]interface{})
	for _, oldDNSServerData := range oldDNSServers {
		oldDNSServer := oldDNSServerData.(map[string]interface{})
		oldDNSServerID, _ := atoi32(oldDNSServer["id"].(string))
		params := access_profiles.NewAccessProfilesDeleteDNSServerParams().WithV(ApiVersion).WithBody(&models.DNSServerDeleteCommand{ID: oldDNSServerID})
		_, err := apiClient.client.AccessProfiles.AccessProfilesDeleteDNSServer(params, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunAccessProfileUpdateDeleteOldNTPServers(d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	oldNTPServersData, _ := d.GetChange("ntp_server")
	oldNTPServers := oldNTPServersData.([]interface{})
	for _, oldNTPServerData := range oldNTPServers {
		oldNTPServer := oldNTPServerData.(map[string]interface{})
		oldNTPServerID, _ := atoi32(oldNTPServer["id"].(string))
		params := access_profiles.NewAccessProfilesDeleteNtpServerParams().WithV(ApiVersion).WithBody(&models.NtpServerDeleteCommand{ID: oldNTPServerID})
		_, err := apiClient.client.AccessProfiles.AccessProfilesDeleteNtpServer(params, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunAccessProfileUpdateDeleteOldSSHUsers(d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	oldSSHUsersData, _ := d.GetChange("ssh_user")
	oldSSHUsers := oldSSHUsersData.([]interface{})
	for _, oldSSHUserData := range oldSSHUsers {
		oldSSHUser := oldSSHUserData.(map[string]interface{})
		oldSSHUserID, _ := atoi32(oldSSHUser["id"].(string))
		params := ssh_users.NewSSHUsersDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteSSHUserCommand{ID: oldSSHUserID})
		_, err := apiClient.client.SSHUsers.SSHUsersDelete(params, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunAccessProfileUpsertSetBody(d *schema.ResourceData, body *models.CreateAccessProfileCommand) {
	if DNSServers, isDNSServersSet := d.GetOk("dns_server"); isDNSServersSet {
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
	if NtpServers, isNTPServersSet := d.GetOk("ntp_server"); isNTPServersSet {
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
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		body.OrganizationID, _ = atoi32(organizationIDData.(string))
	}
	if proxy, isProxySet := d.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}
	if SSHUsers, isSSHUsersSet := d.GetOk("ssh_user"); isSSHUsersSet {
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
}

func resourceTaikunAccessProfileLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.AccessProfilesLockManagementCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := access_profiles.NewAccessProfilesLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.AccessProfiles.AccessProfilesLockManager(params, apiClient)
	return err
}
