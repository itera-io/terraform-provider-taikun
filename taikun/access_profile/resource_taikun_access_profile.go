package access_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			ValidateDiagFunc: utils.StringIsInt,
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

func ResourceTaikunAccessProfile() *schema.Resource {
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
	apiClient := meta.(*tk.Client)

	body := &tkcore.CreateAccessProfileCommand{}
	body.SetName(d.Get("name").(string))

	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		orgId, _ := utils.Atoi32(organizationIDData.(string))
		body.SetOrganizationId(orgId)
	}
	if proxy, isProxySet := d.GetOk("http_proxy"); isProxySet {
		body.SetHttpProxy(proxy.(string))
	}

	resourceTaikunAccessProfileCreateAllowedHosts(d, body)
	resourceTaikunAccessProfileCreateDnsServers(d, body)
	resourceTaikunAccessProfileCreateNtpServers(d, body)
	resourceTaikunAccessProfileCreateSshUsers(d, body)

	response, res, err := apiClient.Client.AccessProfilesAPI.AccessprofilesCreate(context.TODO()).CreateAccessProfileCommand(*body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	id, _ := utils.Atoi32(response.GetId())

	utils.SetResourceDataId(d, id)

	err = resourceTaikunAccessProfileCreateLock(d, id, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}

// set allowed hosts in create request's body
func resourceTaikunAccessProfileCreateAllowedHosts(d *schema.ResourceData, body *tkcore.CreateAccessProfileCommand) {
	if allowedHostData, allowedHostIsSet := d.GetOk("allowed_host"); allowedHostIsSet {
		allowedHosts := allowedHostData.([]interface{})
		body.AllowedHosts = make([]tkcore.AllowedHostCreateDto, len(allowedHosts))
		for i, rawAllowedHost := range allowedHosts {
			allowedHost := rawAllowedHost.(map[string]interface{})
			body.AllowedHosts[i] = tkcore.AllowedHostCreateDto{}
			body.AllowedHosts[i].SetDescription(allowedHost["description"].(string))
			body.AllowedHosts[i].SetIpAddress(allowedHost["address"].(string))
			body.AllowedHosts[i].SetMaskBits(int32(allowedHost["mask_bits"].(int)))
		}
	}
}

// set DNS servers in create request's body
func resourceTaikunAccessProfileCreateDnsServers(d *schema.ResourceData, body *tkcore.CreateAccessProfileCommand) {
	if dnsServerData, dnsServerIsSet := d.GetOk("dns_server"); dnsServerIsSet {
		dnsServers := dnsServerData.([]interface{})
		body.SetDnsServers(make([]tkcore.DnsServerCreateDto, len(dnsServers)))
		for i, rawDnsServer := range dnsServers {
			dnsServer := rawDnsServer.(map[string]interface{})
			body.DnsServers[i] = tkcore.DnsServerCreateDto{}
			body.DnsServers[i].SetAddress(dnsServer["address"].(string))
		}
	}
}

// set NTP servers in create request's body
func resourceTaikunAccessProfileCreateNtpServers(d *schema.ResourceData, body *tkcore.CreateAccessProfileCommand) {
	if ntpServerData, ntpServerIsSet := d.GetOk("ntp_server"); ntpServerIsSet {
		ntpServers := ntpServerData.([]interface{})
		body.NtpServers = make([]tkcore.NtpServerCreateDto, len(ntpServers))
		for i, rawNtpServer := range ntpServers {
			ntpServer := rawNtpServer.(map[string]interface{})
			body.NtpServers[i] = tkcore.NtpServerCreateDto{}
			body.NtpServers[i].SetAddress(ntpServer["address"].(string))
		}
	}
}

// set SSH users in create request's body
func resourceTaikunAccessProfileCreateSshUsers(d *schema.ResourceData, body *tkcore.CreateAccessProfileCommand) {
	if sshUserData, sshUserIsSet := d.GetOk("ssh_user"); sshUserIsSet {
		sshUsers := sshUserData.([]interface{})
		body.SetSshUsers(make([]tkcore.SshUserCreateDto, len(sshUsers)))
		for i, rawSshUser := range sshUsers {
			sshUser := rawSshUser.(map[string]interface{})
			body.SshUsers[i] = tkcore.SshUserCreateDto{}
			body.SshUsers[i].SetName(sshUser["name"].(string))
			body.SshUsers[i].SetSshPublicKey(sshUser["public_key"].(string))
		}
	}
}

// lock access profile after creation
func resourceTaikunAccessProfileCreateLock(d *schema.ResourceData, id int32, apiClient *tk.Client) (err error) {
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
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.AccessProfilesAPI.AccessprofilesList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		sshResponse, res, err := apiClient.Client.SshUsersAPI.SshusersList(context.TODO(), id).Execute()
		if err != nil {
			/*
				if _, ok := err.(*ssh_users.SSHUsersListNotFound); ok && withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
			*/
			return diag.FromErr(tk.CreateError(res, err))
		}

		rawAccessProfile := response.Data[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunAccessProfile(&rawAccessProfile, sshResponse))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

// Update an access profile resource.
// If there have been changes to the allowed hosts, DNS servers, NTP servers or
// SSH users, these will be deleted and recreated as there is no easy way to
// tell which of them are new and which have been modified.
func resourceTaikunAccessProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
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

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunAccessProfileReadWithRetries(), ctx, d, meta)
}

// Update the access profile's HTTP proxy
func resourceTaikunAccessProfileUpdateHttpProxy(d *schema.ResourceData, id int32, apiClient *tk.Client) (err error) {
	if d.HasChange("http_proxy") {

		body := tkcore.UpdateAccessProfileDto{}
		body.SetName(d.Get("name").(string))

		if newHttpProxy, newHttpProxyIsSet := d.GetOk("http_proxy"); newHttpProxyIsSet {
			body.SetHttpProxy(newHttpProxy.(string))
		} else {
			body.SetHttpProxy("")
		}
		res, newErr := apiClient.Client.AccessProfilesAPI.AccessprofilesUpdate(context.TODO(), id).UpdateAccessProfileDto(body).Execute()
		if newErr != nil {
			return tk.CreateError(res, newErr)
		}
	}
	return err
}

// Update the access profile's allowed hosts
func resourceTaikunAccessProfileUpdateAllowedHosts(d *schema.ResourceData, accessProfileId int32, apiClient *tk.Client) (err error) {
	if !d.HasChange("allowed_host") {
		return
	}

	// Delete old allowed hosts
	oldAllowedHostData, newAllowedHostData := d.GetChange("allowed_host")
	oldAllowedHosts := oldAllowedHostData.([]interface{})
	for _, rawOldAllowedHost := range oldAllowedHosts {
		oldAllowedHost := rawOldAllowedHost.(map[string]interface{})
		id, _ := utils.Atoi32(oldAllowedHost["id"].(string))
		if res, err := apiClient.Client.AllowedHostAPI.AllowedhostDelete(context.TODO(), id).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	// Add new allowed hosts
	newAllowedHosts := newAllowedHostData.([]interface{})
	for _, rawNewAllowedHost := range newAllowedHosts {
		newAllowedHost := rawNewAllowedHost.(map[string]interface{})
		body := tkcore.CreateAllowedHostCommand{}
		body.SetAccessProfileId(accessProfileId)
		body.SetDescription(newAllowedHost["description"].(string))
		body.SetIpAddress(newAllowedHost["address"].(string))
		body.SetMaskBits(int32(newAllowedHost["mask_bits"].(int)))

		if _, res, err := apiClient.Client.AllowedHostAPI.AllowedhostCreate(context.TODO()).CreateAllowedHostCommand(body).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	return
}

// Update the access profile's DNS servers
func resourceTaikunAccessProfileUpdateDnsServers(d *schema.ResourceData, accessProfileId int32, apiClient *tk.Client) (err error) {
	if !d.HasChange("dns_server") {
		return
	}

	// Delete old servers
	oldDnsServerData, newDnsServerData := d.GetChange("dns_server")
	oldDnsServers := oldDnsServerData.([]interface{})
	for _, rawOldDnsServer := range oldDnsServers {
		oldDnsServer := rawOldDnsServer.(map[string]interface{})
		id, _ := utils.Atoi32(oldDnsServer["id"].(string))
		if res, err2 := apiClient.Client.DnsServersAPI.DnsserversDelete(context.TODO(), id).Execute(); err != nil {
			err = tk.CreateError(res, err2)
			return err
		}
	}

	// Add new servers
	newDnsServers := newDnsServerData.([]interface{})
	for _, rawNewDnsServer := range newDnsServers {
		newDnsServer := rawNewDnsServer.(map[string]interface{})
		body := tkcore.CreateDnsServerCommand{}
		body.SetAccessProfileId(accessProfileId)
		body.SetAddress(newDnsServer["address"].(string))
		if _, _, err = apiClient.Client.DnsServersAPI.DnsserversCreate(context.TODO()).CreateDnsServerCommand(body).Execute(); err != nil {
			return
		}
	}

	return
}

// Update the access profile's NTP servers
func resourceTaikunAccessProfileUpdateNtpServers(d *schema.ResourceData, accessProfileId int32, apiClient *tk.Client) (err error) {
	if !d.HasChange("ntp_server") {
		return
	}

	// Delete old servers
	oldNtpServerData, newNtpServerData := d.GetChange("ntp_server")
	oldNtpServers := oldNtpServerData.([]interface{})
	for _, rawOldNtpServer := range oldNtpServers {
		oldNtpServer := rawOldNtpServer.(map[string]interface{})
		id, _ := utils.Atoi32(oldNtpServer["id"].(string))
		if res, err := apiClient.Client.NtpServersAPI.NtpserversDelete(context.TODO(), id).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	// Add new servers
	newNtpServers := newNtpServerData.([]interface{})
	for _, rawNewNtpServer := range newNtpServers {
		newNtpServer := rawNewNtpServer.(map[string]interface{})
		body := tkcore.CreateNtpServerCommand{}
		body.SetAccessProfileId(accessProfileId)
		body.SetAddress(newNtpServer["address"].(string))
		if _, res, err := apiClient.Client.NtpServersAPI.NtpserversCreate(context.TODO()).CreateNtpServerCommand(body).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	return
}

// Update the access profile's SSH users
func resourceTaikunAccessProfileUpdateSshUsers(d *schema.ResourceData, accessProfileId int32, apiClient *tk.Client) (err error) {
	if !d.HasChange("ssh_user") {
		return
	}

	// Delete old SSH users
	oldSshUserData, newSshUserData := d.GetChange("ssh_user")
	oldSshUsers := oldSshUserData.([]interface{})
	for _, rawOldSshUser := range oldSshUsers {
		oldSshUser := rawOldSshUser.(map[string]interface{})
		id, _ := utils.Atoi32(oldSshUser["id"].(string))
		body := tkcore.DeleteSshUserCommand{}
		body.SetId(id)
		if res, err := apiClient.Client.SshUsersAPI.SshusersDelete(context.TODO()).DeleteSshUserCommand(body).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	// Add new SSH users
	newSshUsers := newSshUserData.([]interface{})
	for _, rawNewSshUser := range newSshUsers {
		newSshUser := rawNewSshUser.(map[string]interface{})
		body := tkcore.CreateSshUserCommand{}
		body.SetAccessProfileId(accessProfileId)
		body.SetName(newSshUser["name"].(string))
		body.SetSshPublicKey(newSshUser["public_key"].(string))

		if _, res, err := apiClient.Client.SshUsersAPI.SshusersCreate(context.TODO()).CreateSshUserCommand(body).Execute(); err != nil {
			err = tk.CreateError(res, err)
			return err
		}
	}

	return
}

func resourceTaikunAccessProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.AccessProfilesAPI.AccessprofilesDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunAccessProfile(rawAccessProfile *tkcore.AccessProfilesListDto, sshResponse []tkcore.SshUsersListDto) map[string]interface{} {
	AllowedHosts := make([]map[string]interface{}, len(rawAccessProfile.AllowedHosts))
	for i, rawAllowedHost := range rawAccessProfile.AllowedHosts {
		AllowedHosts[i] = map[string]interface{}{
			"description": rawAllowedHost.GetDescription(),
			"id":          utils.I32toa(rawAllowedHost.GetId()),
			"address":     rawAllowedHost.GetIpAddress(),
			"mask_bits":   rawAllowedHost.GetMaskBits(),
		}
	}

	DNSServers := make([]map[string]interface{}, len(rawAccessProfile.GetDnsServers()))
	for i, rawDNSServer := range rawAccessProfile.GetDnsServers() {
		DNSServers[i] = map[string]interface{}{
			"address": rawDNSServer.GetAddress(),
			"id":      utils.I32toa(rawDNSServer.GetId()),
		}
	}

	NTPServers := make([]map[string]interface{}, len(rawAccessProfile.NtpServers))
	for i, rawNTPServer := range rawAccessProfile.NtpServers {
		NTPServers[i] = map[string]interface{}{
			"address": rawNTPServer.GetAddress(),
			"id":      utils.I32toa(rawNTPServer.GetId()),
		}
	}

	SSHUsers := make([]map[string]interface{}, len(sshResponse))
	for i, rawSSHUser := range sshResponse {
		SSHUsers[i] = map[string]interface{}{
			"id":         utils.I32toa(rawSSHUser.GetId()),
			"name":       rawSSHUser.GetName(),
			"public_key": rawSSHUser.GetSshPublicKey(),
		}
	}

	return map[string]interface{}{
		"allowed_host":      AllowedHosts,
		"created_by":        rawAccessProfile.GetCreatedBy(),
		"dns_server":        DNSServers,
		"http_proxy":        rawAccessProfile.GetHttpProxy(),
		"id":                utils.I32toa(rawAccessProfile.GetId()),
		"last_modified":     rawAccessProfile.GetLastModified(),
		"last_modified_by":  rawAccessProfile.GetLastModifiedBy(),
		"lock":              rawAccessProfile.GetIsLocked(),
		"name":              rawAccessProfile.GetName(),
		"ntp_server":        NTPServers,
		"organization_id":   utils.I32toa(rawAccessProfile.GetOrganizationId()),
		"organization_name": rawAccessProfile.GetOrganizationName(),
		"ssh_user":          SSHUsers,
	}
}

func resourceTaikunAccessProfileLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.AccessProfilesLockManagementCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	res, err := apiClient.Client.AccessProfilesAPI.AccessprofilesLockManager(context.TODO()).AccessProfilesLockManagementCommand(body).Execute()
	return tk.CreateError(res, err)
}
