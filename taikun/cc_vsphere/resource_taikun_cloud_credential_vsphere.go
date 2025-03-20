package cc_vsphere

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunCloudCredentialVsphereSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the vSphere cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the vSphere cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the vSphere cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the vSphere cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the vSphere cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the vSphere cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the vSphere cloud credential.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Computed:    true,
		},
		"name": {
			Description: "The name of the vSphere cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or non alpha numeric (-)",
				),
			),
		},
		"username": {
			Description:  "The vSphere Client ID. (Can be set with env VSPHERE_USERNAME)",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Description:  "The vSphere Client Secret. (Can be set with env VSPHERE_PASSWORD)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"api_host": {
			Description:  "The vSphere authentication URL. (Can be set with env VSPHERE_API_URL)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_API_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"datacenter": {
			Description:  "The chosen vSphere datacenter name. (Can be set with env VSPHERE_DATACENTER)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_DATACENTER", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_pool": {
			Description:  "The chosen vSphere resource pool. (Can be set with env VSPHERE_RESOURCE_POOL)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_RESOURCE_POOL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"data_store": {
			Description:  "The chosen vSphere datastore. (Can be set with env VSPHERE_DATA_STORE)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_DATA_STORE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"drs_enabled": {
			Description: "Do you wish to enable vSphere Distributed Resource Scheduler? (Can be set with env VSPHERE_DRS_ENABLED)",
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("VSPHERE_DRS_ENABLED", nil),
			//AtLeastOneOf: []string{"hypervisors"},
		},
		"hypervisors": {
			Description: "The vSphere hypervisors string array.",
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			//AtLeastOneOf: []string{"drs_enabled"},
		},

		"vm_template_name": {
			Description:  "The vSphere VM template name. (Can be set with env VSPHERE_VM_TEMPLATE)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_VM_TEMPLATE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"continent": {
			Type:        schema.TypeString,
			Description: "The OpenStack continent (e.g., `Africa`, `Asia`, `Europe`, `North America`, `Oceania`, or `South America`).",
			Optional:    true,
			Default:     "Europe",
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Africa",
				"Asia",
				"Europe",
				"North America",
				"Oceania",
				"South America",
			}, false),
		},
		"public_name": {
			Description:  "Public network name. (Can be set with env VSPHERE_PUBLIC_NETWORK_NAME)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_NETWORK_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_ip_address": {
			Description:  "Public network address IP. (Can be set with env VSPHERE_PUBLIC_NETWORK_ADDRESS)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_NETWORK_ADDRESS", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_net_mask": {
			Description:  "Public network mask in number format - CIDR after /. (Can be set with env VSPHERE_PUBLIC_NETMASK)",
			Type:         schema.TypeInt,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_NETMASK", nil),
			ValidateFunc: validation.IntBetween(1, 32),
			Required:     true,
			ForceNew:     true,
		},
		"public_gateway": {
			Description:  "Public network gateway IP, must be inside defined network. (Can be set with env VSPHERE_PUBLIC_GATEWAY)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_GATEWAY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_begin_allocation_range": {
			Description:  "Start of allocation range from public network, must be inside defined network. (Can be set with env VSPHERE_PUBLIC_BEGIN_RANGE)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_BEGIN_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"public_end_allocation_range": {
			Description:  "End of allocation range from public network, must be inside defined network. (Can be set with env VSPHERE_PUBLIC_END_RANGE)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PUBLIC_END_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_name": {
			Description:  "Private network name. (Can be set with env VSPHERE_PRIVATE_NETWORK_NAME)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_NETWORK_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_ip_address": {
			Description:  "Private network address IP. (Can be set with env VSPHERE_PRIVATE_NETWORK_ADDRESS)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_NETWORK_ADDRESS", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_net_mask": {
			Description:  "Private network mask in number format - CIDR after /. (Can be set with env VSPHERE_PRIVATE_NETMASK)",
			Type:         schema.TypeInt,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_NETMASK", nil),
			ValidateFunc: validation.IntBetween(1, 32),
			Required:     true,
			ForceNew:     true,
		},
		"private_gateway": {
			Description:  "Private network gateway IP, must be inside defined network. (Can be set with env VSPHERE_PRIVATE_GATEWAY)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_GATEWAY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_begin_allocation_range": {
			Description:  "Start of allocation range from private network, must be inside defined network. (Can be set with env VSPHERE_PRIVATE_BEGIN_RANGE)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_BEGIN_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
		"private_end_allocation_range": {
			Description:  "End of allocation range from private network, must be inside defined network. (Can be set with env VSPHERE_PRIVATE_END_RANGE)",
			Type:         schema.TypeString,
			DefaultFunc:  schema.EnvDefaultFunc("VSPHERE_PRIVATE_END_RANGE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
			Required:     true,
			ForceNew:     true,
		},
	}
}

func ResourceTaikunCloudCredentialVsphere() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun vSphere Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialVsphereCreate,
		ReadContext:   generateResourceTaikunCloudCredentialVsphereReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialVsphereUpdate,
		DeleteContext: utils.ResourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialVsphereSchema(),
	}
}

func resourceTaikunCloudCredentialVsphereCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateVsphereCommand{}
	body.SetName(d.Get("name").(string))

	body.SetUsername(d.Get("username").(string))
	body.SetPassword(d.Get("password").(string))
	body.SetUrl(d.Get("api_host").(string))
	body.SetDatacenterName(d.Get("datacenter").(string))
	body.SetResourcePoolName(d.Get("resource_pool").(string))
	body.SetDatastoreName(d.Get("data_store").(string))
	body.SetDrsEnabled(d.Get("drs_enabled").(bool))
	body.SetHypervisors(utils.ResourceGetStringList(d.Get("hypervisors")))
	body.SetVmTemplateName(d.Get("vm_template_name").(string))
	body.SetContinent(d.Get("continent").(string))

	datacenterId, err := getDatacenterId(d.Get("datacenter").(string), d.Get("api_host").(string), d.Get("username").(string), d.Get("password").(string), meta)
	if err != nil {
		return diag.FromErr(err)
	}
	body.SetDatacenterId(datacenterId)

	publicNetwork := tkcore.CreateVsphereNetworkDto{}
	publicNetwork.SetName(d.Get("public_name").(string))
	publicNetwork.SetIpAddress(d.Get("public_ip_address").(string))
	publicNetwork.SetNetMask(int32(d.Get("public_net_mask").(int)))
	publicNetwork.SetGateway(d.Get("public_gateway").(string))
	publicNetwork.SetBeginAllocationRange(d.Get("public_begin_allocation_range").(string))
	publicNetwork.SetEndAllocationRange(d.Get("public_end_allocation_range").(string))
	body.SetPublicNetwork(publicNetwork)

	privateNetwork := tkcore.CreateVsphereNetworkDto{}
	privateNetwork.SetName(d.Get("private_name").(string))
	privateNetwork.SetIpAddress(d.Get("private_ip_address").(string))
	privateNetwork.SetNetMask(int32(d.Get("private_net_mask").(int)))
	privateNetwork.SetGateway(d.Get("private_gateway").(string))
	privateNetwork.SetBeginAllocationRange(d.Get("private_begin_allocation_range").(string))
	privateNetwork.SetEndAllocationRange(d.Get("private_end_allocation_range").(string))
	body.SetPrivateNetwork(privateNetwork)

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	continentData, continentIsSet := d.GetOk("continent")
	if continentIsSet {
		body.SetContinent(utils.ContinentShorthand(continentData.(string)))
	}

	createResult, res, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereCreate(context.TODO()).CreateVsphereCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := utils.Atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialVsphereLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCloudCredentialVsphereReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialVsphereReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialVsphereRead(true)
}
func generateResourceTaikunCloudCredentialVsphereReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialVsphereRead(false)
}
func generateResourceTaikunCloudCredentialVsphereRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialVsphere := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunCloudCredentialVsphere(&rawCloudCredentialVsphere))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialVsphereUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialVsphereLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("client_id", "client_secret", "name") {
		updateBody := tkcore.UpdateVsphereCommand{}
		updateBody.SetId(id)
		updateBody.SetName(d.Get("name").(string))
		updateBody.SetUsername(d.Get("username").(string))
		updateBody.SetPassword(d.Get("password").(string))

		res, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereUpdate(context.TODO()).UpdateVsphereCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChanges("hypervisors") {
		updateBody := tkcore.UpdateVsphereHypervisorsCommand{}
		updateBody.SetId(id)
		updateBody.SetHypervisors(utils.ResourceGetStringList(d.Get("hypervisors")))

		res, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereUpdateVsphereHypervisors(context.TODO()).UpdateVsphereHypervisorsCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialVsphereLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCloudCredentialVsphereReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialVsphere(rawVsphereCredential *tkcore.VsphereListDto) map[string]interface{} {

	// Transforming slice of nullable strings to slice of strings
	var hypervisorsSlice []string
	for _, hypervisor := range rawVsphereCredential.GetHypervisors() {
		hypervisorsSlice = append(hypervisorsSlice, *hypervisor.Name.Get())
	}

	// Sort networks to private and public
	var privateNetwork map[string]interface{}
	var publicNetwork map[string]interface{}
	for _, network := range rawVsphereCredential.GetVsphereNetworks() {
		// Get public network
		if !network.GetIsPrivate() {
			publicNetwork = map[string]interface{}{
				"public_name":                   network.Name.Get(),
				"public_ip_address":             network.IpAddress,
				"public_net_mask":               network.NetMask,
				"public_gateway":                network.Gateway,
				"public_begin_allocation_range": network.BeginAllocationRange,
				"public_end_allocation_range":   network.EndAllocationRange,
			}
		}
		// Get private network
		if network.GetIsPrivate() {
			privateNetwork = map[string]interface{}{
				"private_name":                   network.Name.Get(),
				"private_ip_address":             network.IpAddress,
				"private_net_mask":               network.NetMask,
				"private_gateway":                network.Gateway,
				"private_begin_allocation_range": network.BeginAllocationRange,
				"private_end_allocation_range":   network.EndAllocationRange,
			}
		}
	}

	return map[string]interface{}{
		"created_by":        rawVsphereCredential.GetCreatedBy(),
		"id":                utils.I32toa(rawVsphereCredential.GetId()),
		"is_default":        rawVsphereCredential.GetIsDefault(),
		"last_modified":     rawVsphereCredential.GetLastModified(),
		"last_modified_by":  rawVsphereCredential.GetLastModifiedBy(),
		"lock":              rawVsphereCredential.GetIsLocked(),
		"organization_id":   utils.I32toa(rawVsphereCredential.GetOrganizationId()),
		"organization_name": rawVsphereCredential.GetOrganizationName(),
		"name":              rawVsphereCredential.GetName(),
		"username":          rawVsphereCredential.GetUsername(),
		"api_host":          rawVsphereCredential.GetUrl(),
		"datacenter":        rawVsphereCredential.GetDatacenterName(),
		"resource_pool":     rawVsphereCredential.GetResourcePool(),
		"data_store":        rawVsphereCredential.GetDatastore(),
		"drs_enabled":       rawVsphereCredential.GetDrsEnabled(),
		"hypervisors":       hypervisorsSlice,
		"vm_template_name":  rawVsphereCredential.GetVmTemplateName(),
		"continent":         rawVsphereCredential.GetContinentName(),

		"public_name":                   publicNetwork["public_name"],
		"public_ip_address":             publicNetwork["public_ip_address"],
		"public_net_mask":               publicNetwork["public_net_mask"],
		"public_gateway":                publicNetwork["public_gateway"],
		"public_begin_allocation_range": publicNetwork["public_begin_allocation_range"],
		"public_end_allocation_range":   publicNetwork["public_end_allocation_range"],

		"private_name":                   privateNetwork["private_name"],
		"private_ip_address":             privateNetwork["private_ip_address"],
		"private_net_mask":               privateNetwork["private_net_mask"],
		"private_gateway":                privateNetwork["private_gateway"],
		"private_begin_allocation_range": privateNetwork["private_begin_allocation_range"],
		"private_end_allocation_range":   privateNetwork["private_end_allocation_range"],
	}
}

func resourceTaikunCloudCredentialVsphereLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}

func getDatacenterId(datacenterName string, url string, username string, password string, meta interface{}) (string, error) {
	// Get a list of datacenters for this vSphere credential
	apiClient := meta.(*tk.Client)
	body := tkcore.DatacenterListCommand{
		Url:            *tkcore.NewNullableString(&url),
		Username:       *tkcore.NewNullableString(&username),
		Password:       *tkcore.NewNullableString(&password),
		DatacenterName: *tkcore.NewNullableString(&datacenterName),
	}
	data, response, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereDatacenterList(context.TODO()).DatacenterListCommand(body).Execute()
	if err != nil {
		return "", tk.CreateError(response, err)
	}

	// Iterate over the list and find the ID for our datacenter name
	for _, v := range data {
		if v.GetName() == datacenterName {
			// Return the id.
			return v.GetDatacenter(), nil
		}
	}

	// Not found
	return "", fmt.Errorf("Could not find Datacenter ID for the specified Datacenter Name.")
}
