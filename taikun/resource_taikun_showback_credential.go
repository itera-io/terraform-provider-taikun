package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Showback Credential",
		CreateContext: resourceTaikunShowbackCredentialCreate,
		ReadContext:   resourceTaikunShowbackCredentialRead,
		UpdateContext: resourceTaikunShowbackCredentialUpdate,
		DeleteContext: resourceTaikunShowbackCredentialDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The id of the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the billing credential.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"prometheus_username": {
				Description: "The prometheus username.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"prometheus_password": {
				Description: "The prometheus password.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"prometheus_url": {
				Description: "The prometheus url.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"organization_id": {
				Description:  "The id of the organization which owns the showback credential.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
				ForceNew:     true,
			},
			"organization_name": {
				Description: "The name of the organization which owns the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_locked": {
				Description: "Indicates whether the showback credential is locked or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"created_by": {
				Description: "The creator of the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified": {
				Description: "Time of last modification.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The last user who modified the showback credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunShowbackCredentialCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateShowbackCredentialCommand{
		Name:     data.Get("name").(string),
		Password: data.Get("prometheus_password").(string),
		URL:      data.Get("prometheus_url").(string),
		Username: data.Get("prometheus_username").(string),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := showback.NewShowbackCreateCredentialParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.client.Showback.ShowbackCreateCredential(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	//TODO Wait for proper response
	//locked := data.Get("is_locked").(bool)
	//if locked {
	//	id, err := atoi32(createResult.Payload.ID)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//	lockBody := models.OperationCredentialLockManagerCommand{
	//		ID:   id,
	//		Mode: getLockMode(locked),
	//	}
	//	lockParams := ops_credentials.NewOpsCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
	//	_, err = apiClient.client.OpsCredentials.OpsCredentialsLockManager(lockParams, apiClient)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//}

	//data.SetId(createResult.Payload.ID)

	return resourceTaikunShowbackCredentialRead(ctx, data, meta)
}

func resourceTaikunShowbackCredentialRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.Showback.ShowbackCredentialsList(showback.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawShowbackCredential := response.GetPayload().Data[0]

		if err := data.Set("created_by", rawShowbackCredential.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawShowbackCredential.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawShowbackCredential.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawShowbackCredential.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawShowbackCredential.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawShowbackCredential.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawShowbackCredential.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawShowbackCredential.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_password", rawShowbackCredential.Password); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("prometheus_url", rawShowbackCredential.URL); err != nil {
			return diag.FromErr(err)
		}
		//TODO
		if err := data.Set("prometheus_username", ""); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunShowbackCredentialUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("is_locked") {
		lockBody := models.ShowbackCredentialLockCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := showback.NewShowbackLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.Showback.ShowbackLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunBillingCredentialRead(ctx, data, meta)
}

func resourceTaikunShowbackCredentialDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := showback.NewShowbackDeleteShowbackCredentialParams().WithV(ApiVersion).WithBody(&models.DeleteShowbackCredentialCommand{ID: id})
	_, err = apiClient.client.Showback.ShowbackDeleteShowbackCredential(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
