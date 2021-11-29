package taikun

import (
	"context"

	"regexp"

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
			Description: "The ID of the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the Azure cloud credential.",
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
		"subscription_id": {
			Description:  "The Azure subscription ID.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_id": {
			Description:  "The Azure client ID.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tenant_id": {
			Description:  "The Azure tenant ID.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description:  "The Azure client secret.",
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
			Description:      "The ID of the organization which owns the Azure cloud credential.",
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
		"lock": {
			Description: "Indicates whether to lock the Azure cloud credential.",
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
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunCloudCredentialAzure() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Azure Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialAzureCreate,
		ReadContext:   generateResourceTaikunCloudCredentialAzureReadWithoutRetries(),
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
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	if data.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunCloudCredentialAzureReadWithRetries(), ctx, data, meta)
}
func generateResourceTaikunCloudCredentialAzureReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAzureRead(true)
}
func generateResourceTaikunCloudCredentialAzureReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAzureRead(false)
}
func generateResourceTaikunCloudCredentialAzureRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		if len(response.Payload.Azure) != 1 {
			if withRetries {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialAzure := response.GetPayload().Azure[0]

		err = setResourceDataFromMap(data, flattenTaikunCloudCredentialAzure(rawCloudCredentialAzure))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialAzureUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := data.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChanges("client_id", "client_secret", "name") {
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

	if data.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunCloudCredentialAzureReadWithRetries(), ctx, data, meta)
}

func flattenTaikunCloudCredentialAzure(rawAzureCredential *models.AzureCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawAzureCredential.CreatedBy,
		"id":                i32toa(rawAzureCredential.ID),
		"lock":              rawAzureCredential.IsLocked,
		"is_default":        rawAzureCredential.IsDefault,
		"last_modified":     rawAzureCredential.LastModified,
		"last_modified_by":  rawAzureCredential.LastModifiedBy,
		"name":              rawAzureCredential.Name,
		"organization_id":   i32toa(rawAzureCredential.OrganizationID),
		"organization_name": rawAzureCredential.OrganizationName,
		"availability_zone": rawAzureCredential.AvailabilityZone,
		"location":          rawAzureCredential.Location,
		"tenant_id":         rawAzureCredential.TenantID,
	}
}

func resourceTaikunCloudCredentialAzureLock(id int32, lock bool, apiClient *apiClient) error {
	body := models.CloudLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.CloudCredentials.CloudCredentialsLockManager(params, apiClient)
	return err
}
