package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/ops_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBillingCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the billing credential is the organization's default.",
			Type:        schema.TypeBool,
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
		"lock": {
			Description: "Indicates whether to lock the billing credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the billing credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
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
		"prometheus_password": {
			Description:  "The Prometheus password.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"prometheus_url": {
			Description:  "The Prometheus URL.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"prometheus_username": {
			Description:  "The Prometheus username.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Billing Credential",
		CreateContext: resourceTaikunBillingCredentialCreate,
		ReadContext:   generateResourceTaikunBillingCredentialReadWithoutRetries(),
		UpdateContext: resourceTaikunBillingCredentialUpdate,
		DeleteContext: resourceTaikunBillingCredentialDelete,
		Schema:        resourceTaikunBillingCredentialSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunBillingCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.OperationCredentialsCreateCommand{
		Name:               d.Get("name").(string),
		PrometheusPassword: d.Get("prometheus_password").(string),
		PrometheusURL:      d.Get("prometheus_url").(string),
		PrometheusUsername: d.Get("prometheus_username").(string),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := ops_credentials.NewOpsCredentialsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.OpsCredentials.OpsCredentialsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunBillingCredentialLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunBillingCredentialReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunBillingCredentialReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingCredentialRead(true)
}
func generateResourceTaikunBillingCredentialReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingCredentialRead(false)
}
func generateResourceTaikunBillingCredentialRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.OpsCredentials.OpsCredentialsList(ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingCredential := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunBillingCredential(rawBillingCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBillingCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("lock") {
		if err := resourceTaikunBillingCredentialLock(id, d.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunBillingCredentialReadWithRetries(), ctx, d, meta)
}

func resourceTaikunBillingCredentialDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := ops_credentials.NewOpsCredentialsDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.OpsCredentials.OpsCredentialsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
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

func resourceTaikunBillingCredentialLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.OperationCredentialLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := ops_credentials.NewOpsCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.OpsCredentials.OpsCredentialsLockManager(params, apiClient)
	return err
}
