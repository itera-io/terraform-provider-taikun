package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/openstack"
	"github.com/itera-io/taikungoclient/models"
	"regexp"
)

func resourceTaikunCloudCredentialOpenStackSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The id of the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
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
		"user": {
			Description:  "The OpenStack user.",
			Type:         schema.TypeString,
			Required:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Description:  "The OpenStack password.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PASSWORD", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"url": {
			Description:  "The OpenStack authentication url.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_name": {
			Description:  "The OpenStack project name.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PROJECT_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Description: "The OpenStack project id.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"public_network_name": {
			Description:  "The name of the public OpenStack network which will be used.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_INTERFACE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zone": {
			Description: "The OpenStack availability zone.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"domain": {
			Description:  "The OpenStack domain.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", nil),
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
		"volume_type_name": {
			Description: "The OpenStack type of volume.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"imported_network_subnet_id": {
			Description: "The OpenStack subnet id in case you want to import network.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"organization_id": {
			Description:      "The id of the organization which owns the OpenStack cloud credential.",
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
		"is_locked": {
			Description: "Indicates whether the OpenStack cloud credential is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the OpenStack cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the OpenStack cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunCloudCredentialOpenStack() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun OpenStack Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialOpenStackCreate,
		ReadContext:   resourceTaikunCloudCredentialOpenStackRead,
		UpdateContext: resourceTaikunCloudCredentialOpenStackUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialOpenStackSchema(),
	}
}

func resourceTaikunCloudCredentialOpenStackCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateOpenstackCloudCommand{
		Name:                   data.Get("name").(string),
		OpenStackUser:          data.Get("user").(string),
		OpenStackPassword:      data.Get("password").(string),
		OpenStackURL:           data.Get("url").(string),
		OpenStackProject:       data.Get("project_name").(string),
		OpenStackPublicNetwork: data.Get("public_network_name").(string),
		OpenStackDomain:        data.Get("domain").(string),
		OpenStackRegion:        data.Get("region").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	importedNetworkSubnetIDData, importedNetworkSubnetIDDataIsSet := data.GetOk("imported_network_subnet_id")
	if importedNetworkSubnetIDDataIsSet {
		body.OpenStackImportNetwork = true
		body.OpenStackInternalSubnetID = importedNetworkSubnetIDData.(string)
	}

	volumeTypeNameData, volumeTypeNameIsSet := data.GetOk("volume_type_name")
	if volumeTypeNameIsSet {
		body.OpenStackVolumeType = volumeTypeNameData.(string)
	}

	availabilityZoneData, availabilityZoneIsSet := data.GetOk("availability_zone")
	if availabilityZoneIsSet {
		body.OpenStackAvailabilityZone = availabilityZoneData.(string)
	}

	params := openstack.NewOpenstackCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Openstack.OpenstackCreate(params, apiClient)
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
		lockBody := models.CloudLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.CloudCredentials.CloudCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunCloudCredentialOpenStackRead(ctx, data, meta)
}

func resourceTaikunCloudCredentialOpenStackRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(response.Payload.Openstack) == 1 {
		rawCloudCredentialOpenStack := response.GetPayload().Openstack[0]

		if err := data.Set("created_by", rawCloudCredentialOpenStack.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawCloudCredentialOpenStack.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawCloudCredentialOpenStack.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_default", rawCloudCredentialOpenStack.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawCloudCredentialOpenStack.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawCloudCredentialOpenStack.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawCloudCredentialOpenStack.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("user", rawCloudCredentialOpenStack.User); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("project_name", rawCloudCredentialOpenStack.Project); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("project_id", rawCloudCredentialOpenStack.TenantID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("public_network_name", rawCloudCredentialOpenStack.PublicNetwork); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("availability_zone", rawCloudCredentialOpenStack.AvailabilityZone); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("domain", rawCloudCredentialOpenStack.Domain); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("region", rawCloudCredentialOpenStack.Region); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("volume_type_name", rawCloudCredentialOpenStack.VolumeType); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("imported_network_subnet_id", rawCloudCredentialOpenStack.InternalSubnetID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawCloudCredentialOpenStack.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawCloudCredentialOpenStack.OrganizationName); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunCloudCredentialOpenStackUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("username") || data.HasChange("password") || data.HasChange("name") {
		updateBody := &models.UpdateOpenStackCommand{
			ID:                id,
			Name:              data.Get("name").(string),
			OpenStackPassword: data.Get("password").(string),
			OpenStackUser:     data.Get("user").(string),
		}
		updateParams := openstack.NewOpenstackUpdateParams().WithV(ApiVersion).WithBody(updateBody)
		_, err := apiClient.client.Openstack.OpenstackUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("is_locked") {
		lockBody := models.CloudLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.CloudCredentials.CloudCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunCloudCredentialOpenStackRead(ctx, data, meta)
}
