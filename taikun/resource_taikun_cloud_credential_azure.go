package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/azure"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunCloudCredentialAzureSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The id of the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the Azure cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"subscription_id": {
			Description:  "The Azure Subscription ID.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_id": {
			Description:  "The Azure Client ID.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tenant_id": {
			Description:  "The Azure Tenant ID.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description:  "The Azure Client Secret.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zone": {
			Description:  "The Azure availability zone for the location.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"location": {
			Description:  "The Azure location.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The id of the organization which owns the Azure cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "Indicates whether the Azure cloud credential is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the Azure cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunCloudCredentialAzure() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Azure Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialAzureCreate,
		ReadContext:   resourceTaikunCloudCredentialAzureRead,
		UpdateContext: resourceTaikunCloudCredentialAzureUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialAzureSchema(),
	}
}

func resourceTaikunCloudCredentialAzureCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateAzureCloudCommand{
		Name:                  data.Get("name").(string),
		AzureTenantID:         data.Get("tenant_id").(string),
		AzureClientID:         data.Get("client_id").(string),
		AzureClientSecret:     data.Get("client_secret").(string),
		AzureSubscriptionID:   data.Get("subscription_id").(string),
		AzureLocation:         data.Get("location").(string),
		AzureAvailabilityZone: data.Get("availability_zone").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := azure.NewAzureCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Azure.AzureCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

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

	data.SetId(createResult.Payload.ID)

	return resourceTaikunCloudCredentialAzureRead(ctx, data, meta)
}

func resourceTaikunCloudCredentialAzureRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if len(response.Payload.Azure) == 1 {
		rawCloudCredentialAzure := response.GetPayload().Azure[0]

		if err := data.Set("created_by", rawCloudCredentialAzure.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawCloudCredentialAzure.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawCloudCredentialAzure.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_default", rawCloudCredentialAzure.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawCloudCredentialAzure.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawCloudCredentialAzure.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawCloudCredentialAzure.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("availability_zone", rawCloudCredentialAzure.AvailabilityZone); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("location", rawCloudCredentialAzure.Location); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("tenant_id", rawCloudCredentialAzure.TenantID); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("availability_zone", rawCloudCredentialAzure.AvailabilityZone); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawCloudCredentialAzure.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawCloudCredentialAzure.OrganizationName); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunCloudCredentialAzureUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("access_key_id") || data.HasChange("secret_access_key") || data.HasChange("name") {
		updateBody := &models.UpdateAzureCommand{
			ID:                id,
			Name:              data.Get("name").(string),
			AzureClientID:     data.Get("client_id").(string),
			AzureClientSecret: data.Get("client_secret").(string),
		}
		updateParams := azure.NewAzureUpdateParams().WithV(ApiVersion).WithBody(updateBody)
		_, err := apiClient.client.Azure.AzureUpdate(updateParams, apiClient)
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

	return resourceTaikunCloudCredentialAzureRead(ctx, data, meta)
}
