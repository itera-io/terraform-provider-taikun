package cc_zadara

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

func resourceTaikunCloudCredentialZadaraSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key_id": {
			Description:  "The Zadara access key ID. (Can be set with env ZADARA_ACCESS_KEY_ID)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ZADARA_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zones": {
			Description: "The given Zadara availability zones for the region.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"az_count": {
			Description:  "The number of Zadara availability zone expected for the region. (Can be set with env ZADARA_AZ_COUNT)",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
			DefaultFunc:  schema.EnvDefaultFunc("ZADARA_AZ_COUNT", nil),
			Default:      1,
		},
		"created_by": {
			Description: "The creator of the Zadara cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the Zadara cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the Zadara cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the Zadara cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Zadara cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"url": {
			Description:  "The Zadara authentication URL. (Can be set with env ZADARA_AUTH_URL)",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ZADARA_AUTH_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"name": {
			Description: "The name of the Zadara cloud credential.",
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
			Description:      "The ID of the organization which owns the Zadara cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Zadara cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"region": {
			Description: "The Zadara region. (Can be set with env ZADARA_DEFAULT_REGION)",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("ZADARA_DEFAULT_REGION", nil),
			ValidateFunc: validation.StringInSlice(
				[]string{
					"symphony",
				},
				false,
			),
		},
		"secret_access_key": {
			Description:  "The Zadara secret access key. (Can be set with env ZADARA_SECRET_ACCESS_KEY)",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ZADARA_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"volume_type": {
			Description: "The volume type for Zadara. (Can be set with env ZADARA_VOLUME_TYPE)",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("ZADARA_VOLUME_TYPE", nil),
			ValidateFunc: validation.StringInSlice(
				[]string{
					"io1",
					"io2",
					"gp2",
					"gp3",
					"sc1",
					"st1",
					"standard",
					"sbp1",
					"sbg1",
				},
				false,
			),
		},
	}
}

func ResourceTaikunCloudCredentialZadara() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Zadara Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialZadaraCreate,
		ReadContext:   generateResourceTaikunCloudCredentialZadaraReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialZadaraUpdate,
		DeleteContext: utils.ResourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialZadaraSchema(),
	}
}

func resourceTaikunCloudCredentialZadaraCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateZadaraCloudCommand{}
	body.SetName(d.Get("name").(string))
	body.SetZadaraAccessKeyId(d.Get("access_key_id").(string))
	body.SetZadaraSecretAccessKey(d.Get("secret_access_key").(string))
	body.SetZadaraRegion(d.Get("region").(string))
	body.SetZadaraUrl(d.Get("url").(string))
	body.SetZadaraVolumeType(d.Get("volume_type").(string))
	body.SetAzCount(int32(d.Get("az_count").(int)))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.ZadaraCloudCredentialAPI.ZadaraCreate(context.TODO()).CreateZadaraCloudCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := utils.Atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialZadaraLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCloudCredentialZadaraReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialZadaraReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialZadaraRead(true)
}
func generateResourceTaikunCloudCredentialZadaraReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialZadaraRead(false)
}
func generateResourceTaikunCloudCredentialZadaraRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.ZadaraCloudCredentialAPI.ZadaraList(context.TODO()).Id(id).Execute()
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

		rawCloudCredentialZadara := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunCloudCredentialZadara(&rawCloudCredentialZadara))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialZadaraUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialZadaraLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("access_key_id", "secret_access_key", "name") {
		updateBody := tkcore.UpdateZadaraCommand{}
		updateBody.SetId(id)
		updateBody.SetName(d.Get("name").(string))
		updateBody.SetZadaraAccessKeyId(d.Get("access_key_id").(string))
		updateBody.SetZadaraSecretAccessKey(d.Get("secret_access_key").(string))

		res, err := apiClient.Client.ZadaraCloudCredentialAPI.ZadaraUpdate(context.TODO()).UpdateZadaraCommand(updateBody).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialZadaraLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCloudCredentialZadaraReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialZadara(rawZadaraCredential *tkcore.ZadaraCredentialsListDto) map[string]interface{} {
	return map[string]interface{}{
		"created_by":         rawZadaraCredential.GetCreatedBy(),
		"id":                 utils.I32toa(rawZadaraCredential.GetId()),
		"lock":               rawZadaraCredential.GetIsLocked(),
		"is_default":         rawZadaraCredential.GetIsDefault(),
		"last_modified":      rawZadaraCredential.GetLastModified(),
		"last_modified_by":   rawZadaraCredential.GetLastModifiedBy(),
		"name":               rawZadaraCredential.GetName(),
		"organization_id":    utils.I32toa(rawZadaraCredential.GetOrganizationId()),
		"organization_name":  rawZadaraCredential.GetOrganizationName(),
		"availability_zones": rawZadaraCredential.GetAvailabilityZones(),
		"region":             rawZadaraCredential.GetRegion(),
		"az_count":           rawZadaraCredential.GetAvailabilityZonesCount(),
		"url":                rawZadaraCredential.GetZadaraApiUrl(),
		"volume_type":        rawZadaraCredential.GetZadaraVolumeType(),
	}
}

func resourceTaikunCloudCredentialZadaraLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.CloudLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsLockManager(context.TODO()).CloudLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
