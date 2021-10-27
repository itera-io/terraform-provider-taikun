package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/ops_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBillingCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the billing credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"prometheus_username": {
			Description:  "The prometheus username.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"prometheus_password": {
			Description:  "The prometheus password.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"prometheus_url": {
			Description:  "The prometheus URL.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the billing credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the billing credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the billing credential is the organization's default.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Billing Credential",
		CreateContext: resourceTaikunBillingCredentialCreate,
		ReadContext:   generateResourceTaikunBillingCredentialRead(false),
		UpdateContext: resourceTaikunBillingCredentialUpdate,
		DeleteContext: resourceTaikunBillingCredentialDelete,
		Schema:        resourceTaikunBillingCredentialSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunBillingCredentialCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.OperationCredentialsCreateCommand{
		Name:               data.Get("name").(string),
		PrometheusPassword: data.Get("prometheus_password").(string),
		PrometheusURL:      data.Get("prometheus_url").(string),
		PrometheusUsername: data.Get("prometheus_username").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := ops_credentials.NewOpsCredentialsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.OpsCredentials.OpsCredentialsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	locked := data.Get("lock").(bool)
	if locked {
		id, err := atoi32(createResult.Payload.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		lockBody := models.OperationCredentialLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := ops_credentials.NewOpsCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.OpsCredentials.OpsCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunBillingCredentialRead(true), ctx, data, meta)
}

func generateResourceTaikunBillingCredentialRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(data.Id())
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.OpsCredentials.OpsCredentialsList(ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if isAfterUpdateOrCreate {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingCredential := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunBillingCredential(rawBillingCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBillingCredentialUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("lock") {
		lockBody := models.OperationCredentialLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("lock").(bool)),
		}
		lockParams := ops_credentials.NewOpsCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.OpsCredentials.OpsCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunBillingCredentialRead(true), ctx, data, meta)
}

func resourceTaikunBillingCredentialDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := ops_credentials.NewOpsCredentialsDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.OpsCredentials.OpsCredentialsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func flattenTaikunBillingCredential(rawOperationCredential *models.OperationCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":          rawOperationCredential.CreatedBy,
		"id":                  i32toa(rawOperationCredential.ID),
		"lock":                rawOperationCredential.IsLocked,
		"is_default":          rawOperationCredential.IsDefault,
		"last_modified":       rawOperationCredential.LastModified,
		"last_modified_by":    rawOperationCredential.LastModifiedBy,
		"name":                rawOperationCredential.Name,
		"organization_id":     i32toa(rawOperationCredential.OrganizationID),
		"organization_name":   rawOperationCredential.OrganizationName,
		"prometheus_password": rawOperationCredential.PrometheusPassword,
		"prometheus_url":      rawOperationCredential.PrometheusURL,
		"prometheus_username": rawOperationCredential.PrometheusUsername,
	}
}
