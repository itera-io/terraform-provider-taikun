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
			Description: "The id of the billing credential.",
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
			Description:  "The prometheus url.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The id of the organization which owns the billing credential.",
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
		"is_locked": {
			Description: "Indicates whether the billing credential is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the billing credential is the organization's default or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the billing credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunBillingCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Billing Credential",
		CreateContext: resourceTaikunBillingCredentialCreate,
		ReadContext:   resourceTaikunBillingCredentialRead,
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

	locked := data.Get("is_locked").(bool)
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

	return resourceTaikunBillingCredentialRead(ctx, data, meta)
}

func resourceTaikunBillingCredentialRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if len(response.Payload.Data) == 1 {
		rawBillingCredential := response.GetPayload().Data[0]

		if err := data.Set("created_by", rawBillingCredential.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawBillingCredential.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawBillingCredential.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_default", rawBillingCredential.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawBillingCredential.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawBillingCredential.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawBillingCredential.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawBillingCredential.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawBillingCredential.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_password", rawBillingCredential.PrometheusPassword); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_url", rawBillingCredential.PrometheusURL); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_username", rawBillingCredential.PrometheusUsername); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunBillingCredentialUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("is_locked") {
		lockBody := models.OperationCredentialLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := ops_credentials.NewOpsCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.OpsCredentials.OpsCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunBillingCredentialRead(ctx, data, meta)
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
