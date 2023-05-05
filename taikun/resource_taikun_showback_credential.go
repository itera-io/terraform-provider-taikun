package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/models"
	"github.com/itera-io/taikungoclient/showbackclient/showback_credentials"
)

func resourceTaikunShowbackCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the showback credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the showback credential.",
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
		"lock": {
			Description: "Indicates whether to lock the showback credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the showback credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the showback credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the showback credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"password": {
			Description:  "The Prometheus password or other credential.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"url": {
			Description:  "URL of the source.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"username": {
			Description:  "The Prometheus username or other credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Showback Credential",
		CreateContext: resourceTaikunShowbackCredentialCreate,
		ReadContext:   generateResourceTaikunShowbackCredentialReadWithoutRetries(),
		UpdateContext: resourceTaikunShowbackCredentialUpdate,
		DeleteContext: resourceTaikunShowbackCredentialDelete,
		Schema:        resourceTaikunShowbackCredentialSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunShowbackCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	// temporary hack to fix
	err := apiClient.Refresh()
	if err != nil {
		return diag.Errorf("showback_credential_id isn't valid: %s", d.Get("showback_credential_id").(string))
	}

	body := &models.CreateShowbackCredentialCommand{
		Name:     d.Get("name").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := showback_credentials.NewShowbackCredentialsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.ShowbackClient.ShowbackCredentials.ShowbackCredentialsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunShowbackCredentialLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunShowbackCredentialReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunShowbackCredentialReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackCredentialRead(true)
}
func generateResourceTaikunShowbackCredentialReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackCredentialRead(false)
}
func generateResourceTaikunShowbackCredentialRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.ShowbackClient.ShowbackCredentials.ShowbackCredentialsList(showback_credentials.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id), apiClient)
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

		rawShowbackCredential := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunShowbackCredential(rawShowbackCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunShowbackCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("lock") {
		if err := resourceTaikunShowbackCredentialLock(id, d.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunShowbackCredentialReadWithRetries(), ctx, d, meta)
}

func resourceTaikunShowbackCredentialDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := showback_credentials.NewShowbackCredentialsDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.ShowbackClient.ShowbackCredentials.ShowbackCredentialsDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunShowbackCredential(rawShowbackCredential *models.ShowbackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawShowbackCredential.CreatedBy,
		"id":                i32toa(rawShowbackCredential.ID),
		"lock":              rawShowbackCredential.IsLocked,
		"last_modified":     rawShowbackCredential.LastModified,
		"last_modified_by":  rawShowbackCredential.LastModifiedBy,
		"name":              rawShowbackCredential.Name,
		"organization_id":   i32toa(rawShowbackCredential.OrganizationID),
		"organization_name": rawShowbackCredential.OrganizationName,
		"password":          rawShowbackCredential.Password,
		"url":               rawShowbackCredential.URL,
		"username":          rawShowbackCredential.Username,
	}
}

func resourceTaikunShowbackCredentialLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.ShowbackCredentialLockCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := showback_credentials.NewShowbackCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.ShowbackClient.ShowbackCredentials.ShowbackCredentialsLockManager(params, apiClient)
	return err
}
