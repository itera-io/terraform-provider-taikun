package cc_azure

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunCloudCredentialAzureSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"availability_zones": {
			Description: "The given Azure availability zones for the location.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"az_count": {
			Description:  "The number of Azure availability zone expected for the region.",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
			Default:      1,
		},
		"client_id": {
			Description:  "The Azure client ID. (Can be set with env AZURE_CLIENT_ID)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description:  "The Azure client secret. (Can be set with env AZURE_SECRET)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_SECRET", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"created_by": {
			Description: "The creator of the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the Azure cloud credential is the default one.",
			Type:        schema.TypeBool,
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
		"location": {
			Description:  "The Azure location.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"lock": {
			Description: "Indicates whether to lock the Azure cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
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
		"organization_id": {
			Description:      "The ID of the organization which owns the Azure cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Azure cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"subscription_id": {
			Description:  "The Azure subscription ID. (Can be set with env AZURE_SUBSCRIPTION)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_SUBSCRIPTION", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tenant_id": {
			Description:  "The Azure tenant ID. (Can be set with env AZURE_TENANT)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_TENANT", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func ResourceTaikunCloudCredentialAzure() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Azure Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialAzureCreate,
		ReadContext:   generateResourceTaikunCloudCredentialAzureReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialAzureUpdate,
		DeleteContext: utils.ResourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialAzureSchema(),
	}
}

func resourceTaikunCloudCredentialAzureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateAzureCloudCommand{}
	body.SetName(d.Get("name").(string))
	body.SetAzureTenantId(d.Get("tenant_id").(string))
	body.SetAzureClientId(d.Get("client_id").(string))
	body.SetAzureClientSecret(d.Get("client_secret").(string))
	body.SetAzureSubscriptionId(d.Get("subscription_id").(string))
	body.SetAzureLocation(d.Get("location").(string))
	azCount := int32(d.Get("az_count").(int))
	/*
		if err != nil {
			return diag.FromErr(err)
		} else if azCount < 1 || azCount > 3 {
			return diag.Errorf("The az_count value must be between 1 and 3 inclusive.")
		}
	*/
	body.SetAzCount(azCount)

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.AzureCloudCredentialAPI.AzureCreate(context.TODO()).CreateAzureCloudCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := utils.Atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCloudCredentialAzureReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialAzureReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAzureRead(true)
}
func generateResourceTaikunCloudCredentialAzureReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAzureRead(false)
}
func generateResourceTaikunCloudCredentialAzureRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.AzureCloudCredentialAPI.AzureList(context.TODO()).Id(id).Execute()
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

		rawCloudCredentialAzure := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunCloudCredentialAzure(&rawCloudCredentialAzure))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialAzureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("client_id", "client_secret", "name") {
		updateBody := tkcore.UpdateAzureCommand{}
		updateBody.SetId(id)
		updateBody.SetName(d.Get("name").(string))
		updateBody.SetAzureClientId(d.Get("client_id").(string))
		updateBody.SetAzureClientSecret(d.Get("client_secret").(string))

		res, err := apiClient.Client.AzureCloudCredentialAPI.AzureUpdate(context.TODO()).UpdateAzureCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAzureLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCloudCredentialAzureReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialAzure(rawAzureCredential *tkcore.AzureCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":         rawAzureCredential.GetCreatedBy(),
		"id":                 utils.I32toa(rawAzureCredential.GetId()),
		"lock":               rawAzureCredential.GetIsLocked(),
		"is_default":         rawAzureCredential.GetIsDefault(),
		"last_modified":      rawAzureCredential.GetLastModified(),
		"last_modified_by":   rawAzureCredential.GetLastModifiedBy(),
		"name":               rawAzureCredential.GetName(),
		"organization_id":    utils.I32toa(rawAzureCredential.GetOrganizationId()),
		"organization_name":  rawAzureCredential.GetOrganizationName(),
		"availability_zones": rawAzureCredential.GetAvailabilityZones(),
		"location":           rawAzureCredential.GetLocation(),
		"tenant_id":          rawAzureCredential.GetTenantId(),
		"az_count":           rawAzureCredential.GetAvailabilityZonesCount(),
	}
}

func resourceTaikunCloudCredentialAzureLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
