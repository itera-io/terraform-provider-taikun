package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkshowback "github.com/itera-io/taikungoclient/showbackclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	apiClient := meta.(*tk.Client)

	body := tkshowback.CreateShowbackCredentialCommand{}
	body.SetName(d.Get("name").(string))
	body.SetPassword(d.Get("password").(string))
	body.SetUrl(d.Get("url").(string))
	body.SetUsername(d.Get("username").(string))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, resp, err := apiClient.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsCreate(ctx).CreateShowbackCredentialCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(resp, err))
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

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
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, resp, err := apiClient.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(resp, err))
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawShowbackCredential := response.GetData()[0]

		err = setResourceDataFromMap(d, flattenTaikunShowbackCredential(&rawShowbackCredential))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunShowbackCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
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
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(resp, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunShowbackCredential(rawShowbackCredential *tkshowback.ShowbackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":        rawShowbackCredential.GetCreatedBy(),
		"id":                i32toa(rawShowbackCredential.GetId()),
		"lock":              rawShowbackCredential.GetIsLocked(),
		"last_modified":     rawShowbackCredential.GetLastModified(),
		"last_modified_by":  rawShowbackCredential.GetLastModifiedBy(),
		"name":              rawShowbackCredential.GetName(),
		"organization_id":   i32toa(rawShowbackCredential.GetOrganizationId()),
		"organization_name": rawShowbackCredential.GetOrganizationName(),
		"password":          rawShowbackCredential.GetPassword(),
		"url":               rawShowbackCredential.GetUrl(),
		"username":          rawShowbackCredential.GetUsername(),
	}
}

func resourceTaikunShowbackCredentialLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkshowback.ShowbackCredentialLockCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	resp, err := apiClient.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsLockManagement(context.TODO()).ShowbackCredentialLockCommand(body).Execute()
	return tk.CreateError(resp, err)
}
