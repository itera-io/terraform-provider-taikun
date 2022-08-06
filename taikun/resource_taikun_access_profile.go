package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/dns_servers"
	"github.com/itera-io/taikungoclient/client/ntp_servers"
	"github.com/itera-io/taikungoclient/client/ssh_users"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunAccessProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_host": {
			Description: "List of allowed hosts.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"description": {
						Description: "Description of the host.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"id": {
						Description: "ID of the host.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"address": {
						Description:  "IPv4 address of the host",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"mask_bits": {
						Description:  "Number of bits in the network mask.",
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 32),
					},
				},
			},
		},
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
			ForceNew:     true,
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

	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		body.OrganizationID, _ = atoi32(organizationIDData.(string))
	}
	if proxy, isProxySet := d.GetOk("http_proxy"); isProxySet {
		body.HTTPProxy = proxy.(string)
	}

	resourceTaikunAccessProfileCreateAllowedHosts(d, body)
	resourceTaikunAccessProfileCreateDnsServers(d, body)
	resourceTaikunAccessProfileCreateNtpServers(d, body)
	resourceTaikunAccessProfileCreateSshUsers(d, body)

	params := access_profiles.NewAccessProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	id, err := resourceTaikunAccessProfileCreateSendRequest(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	setResourceDataId(d, id)

	err = resourceTaikunAccessProfileCreateLock(d, id, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterCreateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}

// set allowed hosts in create request body
func resourceTaikunAccessProfileCreateAllowedHosts(d *schema.ResourceData, body *models.CreateAccessProfileCommand) {
	if allowedHostData, allowedHostIsSet := d.GetOk("allowed_host"); allowedHostIsSet {
		allowedHosts := allowedHostData.([]interface{})
		body.AllowedHosts = make([]*models.AllowedHostCreateDto, len(allowedHosts))
		for i, rawAllowedHost := range allowedHosts {
			allowedHost := rawAllowedHost.(map[string]interface{})
			body.AllowedHosts[i] = &models.AllowedHostCreateDto{
				Description: allowedHost["description"].(string),
				IPAddress:   allowedHost["address"].(string),
				MaskBits:    allowedHost["mask_bits"].(int32),
			}
		}
	}
}

// set DNS servers in create request body
func resourceTaikunAccessProfileCreateDnsServers(d *schema.ResourceData, body *models.CreateAccessProfileCommand) {
	if dnsServerData, dnsServerIsSet := d.GetOk("dns_server"); dnsServerIsSet {
		dnsServers := dnsServerData.([]interface{})
		body.DNSServers = make([]*models.DNSServerCreateDto, len(dnsServers))
		for i, rawDnsServer := range dnsServers {
			dnsServer := rawDnsServer.(map[string]interface{})
			body.DNSServers[i] = &models.DNSServerCreateDto{
				Address: dnsServer["address"].(string),
			}
		}
	}
}

// set NTP servers in create request body
func resourceTaikunAccessProfileCreateNtpServers(d *schema.ResourceData, body *models.CreateAccessProfileCommand) {
	if ntpServerData, ntpServerIsSet := d.GetOk("ntp_server"); ntpServerIsSet {
		ntpServers := ntpServerData.([]interface{})
		body.NtpServers = make([]*models.NtpServerCreateDto, len(ntpServers))
		for i, rawNtpServer := range ntpServers {
			ntpServer := rawNtpServer.(map[string]interface{})
			body.NtpServers[i] = &models.NtpServerCreateDto{
				Address: ntpServer["address"].(string),
			}
		}
	}
}

// set SSH users in create request body
func resourceTaikunAccessProfileCreateSshUsers(d *schema.ResourceData, body *models.CreateAccessProfileCommand) {
	if sshUserData, sshUserIsSet := d.GetOk("ssh_user"); sshUserIsSet {
		sshUsers := sshUserData.([]interface{})
		body.SSHUsers = make([]*models.SSHUserCreateDto, len(sshUsers))
		for i, rawSshUser := range sshUsers {
			sshUser := rawSshUser.(map[string]interface{})
			body.SSHUsers[i] = &models.SSHUserCreateDto{
				Name:         sshUser["name"].(string),
				SSHPublicKey: sshUser["public_key"].(string),
			}
		}
	}
}

// send access profile creation request
// returns the ID of the new resource or an error
func resourceTaikunAccessProfileCreateSendRequest(params *access_profiles.AccessProfilesCreateParams, apiClient *taikungoclient.Client) (int32, error) {
	createResult, err := apiClient.Client.AccessProfiles.AccessProfilesCreate(params, apiClient)
	if err != nil {
		return 0, err
	}
	return atoi32(createResult.Payload.ID)
}

// lock access profile after creation
func resourceTaikunAccessProfileCreateLock(d *schema.ResourceData, id int32, apiClient *taikungoclient.Client) (err error) {
	if d.Get("lock").(bool) {
		err = resourceTaikunAccessProfileLock(id, true, apiClient)
	}
	return
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

// Update an access profile resource.
// If there have been changes to the allowed hosts, DNS servers, NTP servers or
// SSH users, these will be deleted and recreated as there is no easy way to
// tell which of them are new and which have been modified.
func resourceTaikunAccessProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if isLocked, _ := d.GetChange("lock"); isLocked.(bool) {
		if err := resourceTaikunAccessProfileLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := resourceTaikunAccessProfileUpdateHttpProxy(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunAccessProfileUpdateAllowedHosts(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunAccessProfileUpdateDnsServers(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunAccessProfileUpdateNtpServers(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunAccessProfileUpdateSshUsers(d, id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunAccessProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}

// Update the access profile's HTTP proxy
func resourceTaikunAccessProfileUpdateHttpProxy(d *schema.ResourceData, id int32, apiClient *taikungoclient.Client) (err error) {
	if d.HasChange("http_proxy") {
		body := models.UpdateAccessProfileDto{
			Name:      d.Get("name").(string),
			HTTPProxy: "",
		}

		if newHttpProxy, newHttpProxyIsSet := d.GetOk("http_proxy"); newHttpProxyIsSet {
			body.HTTPProxy = newHttpProxy.(string)
		}

		params := access_profiles.NewAccessProfilesUpdateParams().WithV(ApiVersion).WithID(id)
		_, err = apiClient.Client.AccessProfiles.AccessProfilesUpdate(params, apiClient)
	}
	return
}

// Update the access profile's allowed hosts
func resourceTaikunAccessProfileUpdateAllowedHosts(d *schema.ResourceData, accessProfileId int32, apiClient *taikungoclient.Client) (err error) {
	if !d.HasChange("allowed_host") {
		return
	}

	// Delete old allowed hosts
	// FIXME

	// Add new allowed hosts
	// FIXME

	return
}

// Update the access profile's DNS servers
func resourceTaikunAccessProfileUpdateDnsServers(d *schema.ResourceData, accessProfileId int32, apiClient *taikungoclient.Client) (err error) {
	if !d.HasChange("dns_server") {
		return
	}

	// Delete old servers
	oldDnsServerData, newDnsServerData := d.GetChange("dns_server")
	oldDnsServers := oldDnsServerData.([]interface{})
	for _, rawOldDnsServer := range oldDnsServers {
		oldDnsServer := rawOldDnsServer.(map[string]interface{})
		id, _ := atoi32(oldDnsServer["id"].(string))
		params := dns_servers.NewDNSServersDeleteParams().WithV(ApiVersion).WithID(id)
		if _, _, err = apiClient.Client.DNSServers.DNSServersDelete(params, apiClient); err != nil {
			return
		}
	}

	// Add new servers
	newDnsServers := newDnsServerData.([]interface{})
	for _, rawNewDnsServer := range newDnsServers {
		newDnsServer := rawNewDnsServer.(map[string]interface{})
		body := models.CreateDNSServerCommand{
			AccessProfileID: accessProfileId,
			Address:         newDnsServer["address"].(string),
		}
		params := dns_servers.NewDNSServersCreateParams().WithV(ApiVersion).WithBody(&body)
		if _, err = apiClient.Client.DNSServers.DNSServersCreate(params, apiClient); err != nil {
			return
		}
	}

	return
}

// Update the access profile's NTP servers
func resourceTaikunAccessProfileUpdateNtpServers(d *schema.ResourceData, accessProfileId int32, apiClient *taikungoclient.Client) (err error) {
	if !d.HasChange("ntp_server") {
		return
	}

	// Delete old servers
	oldNtpServerData, newNtpServerData := d.GetChange("ntp_server")
	oldNtpServers := oldNtpServerData.([]interface{})
	for _, rawOldNtpServer := range oldNtpServers {
		oldNtpServer := rawOldNtpServer.(map[string]interface{})
		id, _ := atoi32(oldNtpServer["id"].(string))
		params := ntp_servers.NewNtpServersDeleteParams().WithV(ApiVersion).WithID(id)
		if _, _, err = apiClient.Client.NtpServers.NtpServersDelete(params, apiClient); err != nil {
			return
		}
	}

	// Add new servers
	newNtpServers := newNtpServerData.([]interface{})
	for _, rawNewNtpServer := range newNtpServers {
		newNtpServer := rawNewNtpServer.(map[string]interface{})
		body := models.CreateNtpServerCommand{
			AccessProfileID: accessProfileId,
			Address:         newNtpServer["address"].(string),
		}
		params := ntp_servers.NewNtpServersCreateParams().WithV(ApiVersion).WithBody(&body)
		if _, err = apiClient.Client.NtpServers.NtpServersCreate(params, apiClient); err != nil {
			return
		}
	}

	return
}

// Update the access profile's SSH users
func resourceTaikunAccessProfileUpdateSshUsers(d *schema.ResourceData, accessProfileId int32, apiClient *taikungoclient.Client) (err error) {
	if !d.HasChange("ssh_user") {
		return
	}

	// Delete old SSH users
	oldSshUserData, newSshUserData := d.GetChange("ssh_user")
	oldSshUsers := oldSshUserData.([]interface{})
	for _, rawOldSshUser := range oldSshUsers {
		oldSshUser := rawOldSshUser.(map[string]interface{})
		id, _ := atoi32(oldSshUser["id"].(string))
		params := ssh_users.NewSSHUsersDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteSSHUserCommand{ID: id})
		if _, err = apiClient.Client.SSHUsers.SSHUsersDelete(params, apiClient); err != nil {
			return
		}
	}

	// Add new SSH users
	newSshUsers := newSshUserData.([]interface{})
	for _, rawNewSshUser := range newSshUsers {
		newSshUser := rawNewSshUser.(map[string]interface{})
		body := models.CreateSSHUserCommand{
			AccessProfileID: accessProfileId,
			Name:            newSshUser["name"].(string),
			SSHPublicKey:    newSshUser["public_key"].(string),
		}
		params := ssh_users.NewSSHUsersCreateParams().WithV(ApiVersion).WithBody(&body)
		if _, err = apiClient.Client.SSHUsers.SSHUsersCreate(params, apiClient); err != nil {
			return
		}
	}

	return
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

	AllowedHosts := make([]map[string]interface{}, len(rawAccessProfile.AllowedHosts))
	for i, rawAllowedHost := range rawAccessProfile.AllowedHosts {
		AllowedHosts[i] = map[string]interface{}{
			"description": rawAllowedHost.Description,
			"id":          i32toa(rawAllowedHost.ID),
			"ip_address":  rawAllowedHost.IPAddress,
			"mask_bits":   rawAllowedHost.MaskBits,
		}
	}

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
		"allowed_host":      AllowedHosts,
		"created_by":        rawAccessProfile.CreatedBy,
		"dns_server":        DNSServers,
		"http_proxy":        rawAccessProfile.HTTPProxy,
		"id":                i32toa(rawAccessProfile.ID),
		"last_modified":     rawAccessProfile.LastModified,
		"last_modified_by":  rawAccessProfile.LastModifiedBy,
		"lock":              rawAccessProfile.IsLocked,
		"name":              rawAccessProfile.Name,
		"ntp_server":        NTPServers,
		"organization_id":   i32toa(rawAccessProfile.OrganizationID),
		"organization_name": rawAccessProfile.OrganizationName,
		"ssh_user":          SSHUsers,
	}
}

func resourceTaikunAccessProfileLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.AccessProfilesLockManagementCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := access_profiles.NewAccessProfilesLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.Client.AccessProfiles.AccessProfilesLockManager(params, apiClient)
	return err
}
