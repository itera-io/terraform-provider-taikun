package cc_openstack

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"
)

func resourceTaikunCloudCredentialOpenStackSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"availability_zone": {
			Description: "The OpenStack availability zone.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"created_by": {
			Description: "The creator of the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"continent": {
			Description: "The OpenStack continent (`Asia`, `Europe` or `America`).",
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
		"domain": {
			Description:      "The OpenStack domain. Required if you authenticate with username + password. (Can be set with env OS_USER_DOMAIN_NAME)",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.IgnoreChangeFromEmpty,
			DefaultFunc:      schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", nil),
			ValidateFunc:     validation.StringIsNotEmpty,
		},
		"id": {
			Description: "The ID of the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"imported_network_subnet_id": {
			Description: "The OpenStack network subnet ID to import a network.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"is_default": {
			Description: "Indicates whether the OpenStack cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the OpenStack cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description: "The name of the OpenStack cloud credential.",
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
		"organization_id": {
			Description:      "The ID of the organization which owns the OpenStack cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"password": {
			Description:  "The OpenStack password or Application Credential secret. (Can be set with env OS_PASSWORD)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PASSWORD", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Description: "The OpenStack project ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"project_name": {
			Description:  "The OpenStack project name. (Can be set with env OS_PROJECT_NAME)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PROJECT_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"public_network_name": {
			Description:  "The name of the public OpenStack network to use. (Can be set with env OS_INTERFACE)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_INTERFACE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"region": {
			Description:  "The OpenStack region. (Can be set with env OS_REGION_NAME)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_REGION_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"url": {
			Description:  "The OpenStack authentication URL. (Can be set with env OS_AUTH_URL)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user": {
			Description:  "The OpenStack user or Application credential Id. (Can be set with env OS_USERNAME)",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"using_application_credentials": {
			Description: "Indicate whether you are specifying OpenStack application credential [id + secret] or [Username + Password + Domain] (default).",
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"volume_type_name": {
			Description: "The OpenStack type of volume.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func ResourceTaikunCloudCredentialOpenStack() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun OpenStack Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialOpenStackCreate,
		ReadContext:   generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialOpenStackUpdate,
		DeleteContext: utils.ResourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialOpenStackSchema(),
	}
}

func resourceTaikunCloudCredentialOpenStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateOpenstackCloudCommand{}
	body.SetName(d.Get("name").(string))
	body.SetOpenStackUrl(d.Get("url").(string))
	body.SetOpenStackProject(d.Get("project_name").(string))
	body.SetOpenStackPublicNetwork(d.Get("public_network_name").(string))
	body.SetOpenStackRegion(d.Get("region").(string))

	username := d.Get("user").(string)
	password := d.Get("password").(string)
	domain := d.Get("domain").(string)
	appCredEnabledBool := d.Get("using_application_credentials").(bool)
	//appcredenabled := d.Get("application_credential_enabled").(bool)

	if appCredEnabledBool {
		body.SetOpenStackUser(username)
		body.SetOpenStackPassword(password)
		body.SetApplicationCredEnabled(true)
	} else {
		body.SetOpenStackUser(username)
		body.SetOpenStackPassword(password)
		body.SetOpenStackDomain(domain)
		body.SetApplicationCredEnabled(false)
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	importedNetworkSubnetIDData, importedNetworkSubnetIDDataIsSet := d.GetOk("imported_network_subnet_id")
	if importedNetworkSubnetIDDataIsSet {
		body.SetOpenStackImportNetwork(true)
		body.SetOpenStackInternalSubnetId(importedNetworkSubnetIDData.(string))
	}

	volumeTypeNameData, volumeTypeNameIsSet := d.GetOk("volume_type_name")
	if volumeTypeNameIsSet {
		body.SetOpenStackVolumeType(volumeTypeNameData.(string))
	}

	availabilityZoneData, availabilityZoneIsSet := d.GetOk("availability_zone")
	if availabilityZoneIsSet {
		body.SetOpenStackAvailabilityZone(availabilityZoneData.(string))
	}

	continentData, continentIsSet := d.GetOk("continent")
	if continentIsSet {
		body.SetOpenStackContinent(utils.ContinentShorthand(continentData.(string)))
	}

	createResult, res, err := apiClient.Client.OpenstackCloudCredentialAPI.OpenstackCreate(context.TODO()).CreateOpenstackCloudCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := utils.Atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCloudCredentialOpenStackReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialOpenStackReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialOpenStackRead(true)
}
func generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialOpenStackRead(false)
}
func generateResourceTaikunCloudCredentialOpenStackRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.OpenstackCloudCredentialAPI.OpenstackList(context.TODO()).Id(id).Execute()
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

		rawCloudCredentialOpenStack := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunCloudCredentialOpenStack(&rawCloudCredentialOpenStack))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialOpenStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("user", "password", "name") {
		updateBody := tkcore.UpdateOpenStackCommand{}
		updateBody.SetId(id)
		updateBody.SetName(d.Get("name").(string))

		username := d.Get("user").(string)
		password := d.Get("password").(string)
		appcredenabled := d.Get("using_application_credentials").(bool)

		if appcredenabled {
			updateBody.SetOpenStackUser(username)
			updateBody.SetOpenStackPassword(password)
		} else {
			updateBody.SetOpenStackUser(username)
			updateBody.SetOpenStackPassword(password)
		}

		res, err := apiClient.Client.OpenstackCloudCredentialAPI.OpenstackUpdate(context.TODO()).UpdateOpenStackCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCloudCredentialOpenStackReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialOpenStack(rawOpenStackCredential *tkcore.OpenstackCredentialsListDto) map[string]interface{} {
	return map[string]interface{}{
		"using_application_credentials": rawOpenStackCredential.GetApplicationCredEnabled(),
		"created_by":                    rawOpenStackCredential.GetCreatedBy(),
		"id":                            utils.I32toa(rawOpenStackCredential.GetId()),
		"lock":                          rawOpenStackCredential.GetIsLocked(),
		"is_default":                    rawOpenStackCredential.GetIsDefault(),
		"last_modified":                 rawOpenStackCredential.GetLastModified(),
		"last_modified_by":              rawOpenStackCredential.GetLastModifiedBy(),
		"name":                          rawOpenStackCredential.GetName(),
		"user":                          rawOpenStackCredential.GetUser(),
		"project_name":                  rawOpenStackCredential.GetProject(),
		"project_id":                    rawOpenStackCredential.GetTenantId(),
		"organization_id":               utils.I32toa(rawOpenStackCredential.GetOrganizationId()),
		"organization_name":             rawOpenStackCredential.GetOrganizationName(),
		"public_network_name":           rawOpenStackCredential.GetPublicNetwork(),
		"availability_zone":             rawOpenStackCredential.GetAvailabilityZone(),
		"domain":                        rawOpenStackCredential.GetDomain(),
		"region":                        rawOpenStackCredential.GetRegion(),
		"continent":                     rawOpenStackCredential.GetContinentName(),
		"volume_type_name":              rawOpenStackCredential.GetVolumeType(),
		"imported_network_subnet_id":    rawOpenStackCredential.GetInternalSubnetId(),
	}
}

func resourceTaikunCloudCredentialOpenStackLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
