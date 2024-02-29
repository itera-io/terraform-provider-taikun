package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunCloudCredentialProxmoxSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the Proxmox cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the Proxmox cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the Proxmox cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the Proxmox cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Proxmox cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the Proxmox cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Proxmox cloud credential.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Computed:    true,
		},
		"name": {
			Description: "The name of the Proxmox cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
			//ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or non alpha numeric (-)",
				),
			),
		},
		"continent": {
			Description: "The Proxmox continent (`Asia`, `Europe` or `America`).",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "Europe",
			ValidateFunc: validation.StringInSlice([]string{
				"Asia",
				"Europe",
				"America",
			}, false),
		},
		"api_host": {
			Description:  "The Proxmox authentication URL.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_API_HOST", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_id": {
			Description: "The Proxmox Client ID.",
			Type:        schema.TypeString,
			Required:    true,
			//ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description: "The Proxmox Client Secret.",
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			//ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_CLIENT_SECRET", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"storage": {
			Description:  "The Proxmox storage option.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_STORAGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"vm_template_name": {
			Description:  "The Proxmox VM template name",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_VM_TEMPLATE_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"hypervisors": {
			Description: "The Proxmox hypervisors string array",
			Type:        schema.TypeList,
			Required:    true,
			//ForceNew:    true,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"public_ip_address": {
			Description:  "Public network address IP",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_NETWORK", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_net_mask": {
			Description:  "Public network mask in number format (CIDR after /)",
			Type:         schema.TypeInt,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_NETMASK", nil),
			ValidateFunc: validation.IntBetween(1, 32),
			Required:     true,
			ForceNew:     true,
		},
		"public_gateway": {
			Description:  "Public network gateway IP (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_GATEWAY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_begin_allocation_range": {
			Description:  "Start of allocation range from public network (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_BEGIN_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_end_allocation_range": {
			Description:  "End of allocation range from public network (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_END_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_bridge": {
			Description:  "Bridge interace name for private network",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PUBLIC_BRIDGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_ip_address": {
			Description:  "Private network address IP",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_NETWORK", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_net_mask": {
			Description:  "Private network mask in number format (CIDR after /)",
			Type:         schema.TypeInt,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_NETMASK", nil),
			ValidateFunc: validation.IntBetween(1, 32),
			Required:     true,
			ForceNew:     true,
		},
		"private_gateway": {
			Description:  "Private network gateway IP (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_GATEWAY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_begin_allocation_range": {
			Description:  "Start of allocation range from private network (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_BEGIN_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_end_allocation_range": {
			Description:  "End of allocation range from private network (must be inside defined network)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_END_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_bridge": {
			Description:  "Bridge interace name for private network",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("PROXMOX_PRIVATE_BRIDGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
	}
}

func resourceTaikunCloudCredentialProxmox() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Proxmox Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialProxmoxCreate,
		ReadContext:   generateResourceTaikunCloudCredentialProxmoxReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialProxmoxUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialProxmoxSchema(),
	}
}

func resourceTaikunCloudCredentialProxmoxCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateProxmoxCommand{}
	body.SetName(d.Get("name").(string))
	body.SetUrl(d.Get("api_host").(string))
	body.SetTokenId(d.Get("client_id").(string))
	body.SetTokenSecret(d.Get("client_secret").(string))
	body.SetStorage(d.Get("storage").(string))
	body.SetVmTemplateName(d.Get("vm_template_name").(string))
	body.SetHypervisors(resourceGetStringList(d.Get("hypervisors")))

	publicNetwork := tkcore.CreateProxmoxNetworkDto{}
	publicNetwork.SetGateway(d.Get("public_gateway").(string))
	publicNetwork.SetIpAddress(d.Get("public_ip_address").(string))
	publicNetwork.SetNetMask(int32(d.Get("public_net_mask").(int)))
	publicNetwork.SetBeginAllocationRange(d.Get("public_begin_allocation_range").(string))
	publicNetwork.SetEndAllocationRange(d.Get("public_end_allocation_range").(string))
	publicNetwork.SetBridge(d.Get("public_bridge").(string))
	body.SetPublicNetwork(publicNetwork)

	privateNetwork := tkcore.CreateProxmoxNetworkDto{}
	privateNetwork.SetGateway(d.Get("private_gateway").(string))
	privateNetwork.SetIpAddress(d.Get("private_ip_address").(string))
	privateNetwork.SetNetMask(int32(d.Get("private_net_mask").(int)))
	privateNetwork.SetBeginAllocationRange(d.Get("private_begin_allocation_range").(string))
	privateNetwork.SetEndAllocationRange(d.Get("private_end_allocation_range").(string))
	privateNetwork.SetBridge(d.Get("private_bridge").(string))
	body.SetPrivateNetwork(privateNetwork)

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	continentData, continentIsSet := d.GetOk("continent")
	if continentIsSet {
		body.SetContinent(continentShorthand(continentData.(string)))
	}

	createResult, res, err := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxCreate(context.TODO()).CreateProxmoxCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialProxmoxLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunCloudCredentialProxmoxReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialProxmoxReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialProxmoxRead(true)
}
func generateResourceTaikunCloudCredentialProxmoxReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialProxmoxRead(false)
}
func generateResourceTaikunCloudCredentialProxmoxRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialProxmox := response.GetData()[0]

		err = setResourceDataFromMap(d, flattenTaikunCloudCredentialProxmox(&rawCloudCredentialProxmox))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialProxmoxUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialProxmoxLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("client_id", "client_secret", "name") {
		updateBody := tkcore.UpdateProxmoxCommand{}
		updateBody.SetId(id)
		updateBody.SetName(d.Get("name").(string))
		updateBody.SetTokenId(d.Get("client_id").(string))
		updateBody.SetTokenSecret(d.Get("client_secret").(string))

		res, err := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxUpdate(context.TODO()).UpdateProxmoxCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChanges("hypervisors") {
		updateBody := tkcore.UpdateHypervisorsCommand{}
		updateBody.SetId(id)
		updateBody.SetHypervisors(resourceGetStringList(d.Get("hypervisors")))

		res, err := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxUpdateHypervisors(context.TODO()).UpdateHypervisorsCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialProxmoxLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunCloudCredentialProxmoxReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialProxmox(rawProxmoxCredential *tkcore.ProxmoxListDto) map[string]interface{} {

	// Transforming slice of nullable strings to slice of strings
	var hypervisorsSlice []string
	for _, hypervisor := range rawProxmoxCredential.GetHypervisors() {
		hypervisorsSlice = append(hypervisorsSlice, *hypervisor.Name.Get())
	}

	// Sort networks to private and public
	var privateNetwork map[string]interface{}
	var publicNetwork map[string]interface{}
	for _, network := range rawProxmoxCredential.GetProxmoxNetworks() {
		// Get private network
		if network.GetIsPrivate() {
			privateNetwork = map[string]interface{}{
				"private_ip_address":             network.IpAddress.Get(),
				"private_net_mask":               network.NetMask,
				"private_gateway":                network.Gateway.Get(),
				"private_begin_allocation_range": network.BeginAllocationRange.Get(),
				"private_end_allocation_range":   network.EndAllocationRange.Get(),
				"private_bridge":                 network.Bridge.Get(),
			}
		}
		// Get public network
		if !network.GetIsPrivate() {
			publicNetwork = map[string]interface{}{
				"public_ip_address":             network.IpAddress.Get(),
				"public_net_mask":               network.NetMask,
				"public_gateway":                network.Gateway.Get(),
				"public_begin_allocation_range": network.BeginAllocationRange.Get(),
				"public_end_allocation_range":   network.EndAllocationRange.Get(),
				"public_bridge":                 network.Bridge.Get(),
			}
		}
	}

	return map[string]interface{}{
		"created_by":                     rawProxmoxCredential.GetCreatedBy(),
		"id":                             i32toa(rawProxmoxCredential.GetId()),
		"is_default":                     rawProxmoxCredential.GetIsDefault(),
		"last_modified":                  rawProxmoxCredential.GetLastModified(),
		"last_modified_by":               rawProxmoxCredential.GetLastModifiedBy(),
		"lock":                           rawProxmoxCredential.GetIsLocked(),
		"organization_id":                i32toa(rawProxmoxCredential.GetOrganizationId()),
		"organization_name":              rawProxmoxCredential.GetOrganizationName(),
		"name":                           rawProxmoxCredential.GetName(),
		"continent":                      rawProxmoxCredential.GetContinentName(),
		"api_host":                       rawProxmoxCredential.GetUrl(),
		"client_id":                      rawProxmoxCredential.GetTokenId(),
		"storage":                        rawProxmoxCredential.GetStorage(),
		"vm_template_name":               rawProxmoxCredential.GetVmTemplateName(),
		"hypervisors":                    hypervisorsSlice,
		"private_ip_address":             privateNetwork["private_ip_address"],
		"private_net_mask":               privateNetwork["private_net_mask"],
		"private_gateway":                privateNetwork["private_gateway"],
		"private_begin_allocation_range": privateNetwork["private_begin_allocation_range"],
		"private_end_allocation_range":   privateNetwork["private_end_allocation_range"],
		"private_bridge":                 privateNetwork["private_bridge"],
		"public_ip_address":              publicNetwork["public_ip_address"],
		"public_net_mask":                publicNetwork["public_net_mask"],
		"public_gateway":                 publicNetwork["public_gateway"],
		"public_begin_allocation_range":  publicNetwork["public_begin_allocation_range"],
		"public_end_allocation_range":    publicNetwork["public_end_allocation_range"],
		"public_bridge":                  publicNetwork["public_bridge"],
	}
}

func resourceTaikunCloudCredentialProxmoxLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
