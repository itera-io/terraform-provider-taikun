package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunShowbackCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the showback credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the showback credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"username": {
			Description:  "The prometheus username or other credential.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Description:  "The prometheus password or other credential.",
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
		"lock": {
			Description: "Indicates whether to lock the showback credential.",
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
	}
}

func resourceTaikunShowbackCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Showback Credential",
		CreateContext: resourceTaikunShowbackCredentialCreate,
		ReadContext:   generateResourceTaikunShowbackCredentialRead(false),
		UpdateContext: resourceTaikunShowbackCredentialUpdate,
		DeleteContext: resourceTaikunShowbackCredentialDelete,
		Schema:        resourceTaikunShowbackCredentialSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunShowbackCredentialCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateShowbackCredentialCommand{
		Name:     data.Get("name").(string),
		Password: data.Get("password").(string),
		URL:      data.Get("url").(string),
		Username: data.Get("username").(string),
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
	createResult, err := apiClient.client.Showback.ShowbackCreateCredential(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	if data.Get("lock").(bool) {
		id, _ := atoi32(createResult.Payload.ID)
		if err := resourceTaikunShowbackCredentialLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunShowbackCredentialRead(true), ctx, data, meta)
}

func generateResourceTaikunShowbackCredentialRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		if len(response.Payload.Data) != 1 {
			if isAfterUpdateOrCreate {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawShowbackCredential := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunShowbackCredential(rawShowbackCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunShowbackCredentialUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("lock") {
		if err := resourceTaikunShowbackCredentialLock(id, data.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunShowbackCredentialRead(true), ctx, data, meta)
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

func resourceTaikunShowbackCredentialLock(id int32, lock bool, apiClient *apiClient) error {
	body := models.ShowbackCredentialLockCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := showback.NewShowbackLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.Showback.ShowbackLockManager(params, apiClient)
	return err
}
