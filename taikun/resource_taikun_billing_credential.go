package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	apiClient := meta.(*tk.Client)

	body := tkcore.OperationCredentialsCreateCommand{}
	body.SetName(d.Get("name").(string))
	body.SetPrometheusPassword(d.Get("prometheus_password").(string))
	body.SetPrometheusUrl(d.Get("prometheus_url").(string))
	body.SetPrometheusUsername(d.Get("prometheus_username").(string))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.OperationCredentialsAPI.OpscredentialsCreate(ctx).OperationCredentialsCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

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
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		rawBillingCredential, err := resourceTaikunBillingCredentialFind(id, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if rawBillingCredential == nil {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		err = setResourceDataFromMap(d, flattenTaikunBillingCredential(rawBillingCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBillingCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
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
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.OperationCredentialsAPI.OpscredentialsDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunBillingCredential(rawOperationCredential *tkcore.OperationCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":          rawOperationCredential.GetCreatedBy(),
		"id":                  i32toa(rawOperationCredential.GetId()),
		"lock":                rawOperationCredential.GetIsLocked(),
		"is_default":          rawOperationCredential.GetIsDefault(),
		"last_modified":       rawOperationCredential.GetLastModified(),
		"last_modified_by":    rawOperationCredential.GetLastModifiedBy(),
		"name":                rawOperationCredential.GetName(),
		"organization_id":     i32toa(rawOperationCredential.GetOrganizationId()),
		"organization_name":   rawOperationCredential.GetOrganizationName(),
		"prometheus_password": rawOperationCredential.GetPrometheusPassword(),
		"prometheus_url":      rawOperationCredential.GetPrometheusUrl(),
		"prometheus_username": rawOperationCredential.GetPrometheusUsername(),
	}
}

func resourceTaikunBillingCredentialLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.OperationCredentialLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	res, err := apiClient.Client.OperationCredentialsAPI.OpscredentialsLockManager(context.TODO()).OperationCredentialLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}

// Returns the Billing Credential with the given ID or nil if it wasn't found
func resourceTaikunBillingCredentialFind(id int32, apiClient *tk.Client) (*tkcore.OperationCredentialsListDto, error) {
	params := apiClient.Client.OperationCredentialsAPI.OpscredentialsList(context.TODO())
	var offset int32 = 0

	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return nil, tk.CreateError(res, err)
		}

		for _, billingCredential := range response.Data {
			if billingCredential.GetId() == id {
				return &billingCredential, nil
			}
		}

		offset += int32(len(response.Data))
		if offset == response.GetTotalCount() {
			break
		}
	}

	return nil, nil
}
