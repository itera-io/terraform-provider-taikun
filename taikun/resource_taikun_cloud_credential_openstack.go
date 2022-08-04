package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/openstack"
	"github.com/itera-io/taikungoclient/models"
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
		"domain": {
			Description:  "The OpenStack domain.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
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
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"password": {
			Description:  "The OpenStack password.",
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
			Description:  "The OpenStack project name.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PROJECT_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"public_network_name": {
			Description:  "The name of the public OpenStack network to use.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_INTERFACE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"region": {
			Description:  "The OpenStack region.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_REGION_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"url": {
			Description:  "The OpenStack authentication URL.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user": {
			Description:  "The OpenStack user.",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"volume_type_name": {
			Description: "The OpenStack type of volume.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func resourceTaikunCloudCredentialOpenStack() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun OpenStack Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialOpenStackCreate,
		ReadContext:   generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialOpenStackUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialOpenStackSchema(),
	}
}

func resourceTaikunCloudCredentialOpenStackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	body := &models.CreateOpenstackCloudCommand{
		Name:                   d.Get("name").(string),
		OpenStackUser:          d.Get("user").(string),
		OpenStackPassword:      d.Get("password").(string),
		OpenStackURL:           d.Get("url").(string),
		OpenStackProject:       d.Get("project_name").(string),
		OpenStackPublicNetwork: d.Get("public_network_name").(string),
		OpenStackDomain:        d.Get("domain").(string),
		OpenStackRegion:        d.Get("region").(string),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	importedNetworkSubnetIDData, importedNetworkSubnetIDDataIsSet := d.GetOk("imported_network_subnet_id")
	if importedNetworkSubnetIDDataIsSet {
		body.OpenStackImportNetwork = true
		body.OpenStackInternalSubnetID = importedNetworkSubnetIDData.(string)
	}

	volumeTypeNameData, volumeTypeNameIsSet := d.GetOk("volume_type_name")
	if volumeTypeNameIsSet {
		body.OpenStackVolumeType = volumeTypeNameData.(string)
	}

	availabilityZoneData, availabilityZoneIsSet := d.GetOk("availability_zone")
	if availabilityZoneIsSet {
		body.OpenStackAvailabilityZone = availabilityZoneData.(string)
	}

	params := openstack.NewOpenstackCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.Client.Openstack.OpenstackCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunCloudCredentialOpenStackReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialOpenStackReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialOpenStackRead(true)
}
func generateResourceTaikunCloudCredentialOpenStackReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialOpenStackRead(false)
}
func generateResourceTaikunCloudCredentialOpenStackRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Openstack) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialOpenStack := response.GetPayload().Openstack[0]

		err = setResourceDataFromMap(d, flattenTaikunCloudCredentialOpenStack(rawCloudCredentialOpenStack))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialOpenStackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("user", "password", "name") {
		updateBody := &models.UpdateOpenStackCommand{
			ID:                id,
			Name:              d.Get("name").(string),
			OpenStackPassword: d.Get("password").(string),
			OpenStackUser:     d.Get("user").(string),
		}
		updateParams := openstack.NewOpenstackUpdateParams().WithV(ApiVersion).WithBody(updateBody)
		_, err := apiClient.Client.Openstack.OpenstackUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialOpenStackLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunCloudCredentialOpenStackReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialOpenStack(rawOpenStackCredential *models.OpenstackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":                 rawOpenStackCredential.CreatedBy,
		"id":                         i32toa(rawOpenStackCredential.ID),
		"lock":                       rawOpenStackCredential.IsLocked,
		"is_default":                 rawOpenStackCredential.IsDefault,
		"last_modified":              rawOpenStackCredential.LastModified,
		"last_modified_by":           rawOpenStackCredential.LastModifiedBy,
		"name":                       rawOpenStackCredential.Name,
		"user":                       rawOpenStackCredential.User,
		"project_name":               rawOpenStackCredential.Project,
		"project_id":                 rawOpenStackCredential.TenantID,
		"organization_id":            i32toa(rawOpenStackCredential.OrganizationID),
		"organization_name":          rawOpenStackCredential.OrganizationName,
		"public_network_name":        rawOpenStackCredential.PublicNetwork,
		"availability_zone":          rawOpenStackCredential.AvailabilityZone,
		"domain":                     rawOpenStackCredential.Domain,
		"region":                     rawOpenStackCredential.Region,
		"volume_type_name":           rawOpenStackCredential.VolumeType,
		"imported_network_subnet_id": rawOpenStackCredential.InternalSubnetID,
	}
}

func resourceTaikunCloudCredentialOpenStackLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.CloudLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.Client.CloudCredentials.CloudCredentialsLockManager(params, apiClient)
	return err
}
