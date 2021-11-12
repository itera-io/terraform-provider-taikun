package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/opa_profiles"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOPAProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_repos": {
			Description: "The list of allowed repo.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"forbid_node_port": {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"forbid_http_ingress": {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"forbidden_tags": {
			Description: "The list of forbidden tags.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"id": {
			Description: "The ID of the OPA profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"ingress_whitelist": {
			Description: "Whitelist of the k8s ingress.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"is_default": {
			Description: "Indicates whether the OPA Profile is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Kubernetes profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the OPA profile.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 50),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the OPA profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the OPA profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"require_probe": {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"unique_ingress": {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"unique_service_selector": {
			Description: "",
			Type:        schema.TypeBool,
			Required:    true,
		},
	}
}

func resourceTaikunOPAProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun OPA Profile",
		CreateContext: resourceTaikunOPAProfileCreate,
		ReadContext:   generateResourceTaikunOPAProfileRead(false),
		UpdateContext: resourceTaikunOPAProfileUpdate,
		DeleteContext: resourceTaikunOPAProfileDelete,
		Schema:        resourceTaikunOPAProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunOPAProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateOpaProfileCommand{
		AllowedRepo:           resourceGetStringList(data.Get("allowed_repos")),
		ForbidHTTPIngress:     data.Get("forbid_http_ingress").(bool),
		ForbidNodePort:        data.Get("forbid_node_port").(bool),
		ForbidSpecificTags:    resourceGetStringList(data.Get("forbidden_tags")),
		IngressWhitelist:      resourceGetStringList(data.Get("ingress_whitelist")),
		Name:                  data.Get("name").(string),
		RequireProbe:          data.Get("require_probe").(bool),
		UniqueIngresses:       data.Get("unique_ingress").(bool),
		UniqueServiceSelector: data.Get("unique_service_selector").(bool),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := opa_profiles.NewOpaProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.OpaProfiles.OpaProfilesCreate(params, apiClient)
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
		lockBody := models.OpaProfileLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := opa_profiles.NewOpaProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.OpaProfiles.OpaProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunOPAProfileRead(true), ctx, data, meta)
}

func generateResourceTaikunOPAProfileRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(data.Id())
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.OpaProfiles.OpaProfilesList(opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion).WithID(&id), apiClient)
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

		rawOPAProfile := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunOPAProfile(rawOPAProfile))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunOPAProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("lock") {

		lockBody := models.OpaProfileLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("lock").(bool)),
		}
		lockParams := opa_profiles.NewOpaProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.OpaProfiles.OpaProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChangeExcept("lock") {

		body := &models.OpaProfileUpdateCommand{
			AllowedRepo:           resourceGetStringList(data.Get("allowed_repos")),
			ForbidHTTPIngress:     data.Get("forbid_http_ingress").(bool),
			ForbidNodePort:        data.Get("forbid_node_port").(bool),
			ForbidSpecificTags:    resourceGetStringList(data.Get("forbidden_tags")),
			IngressWhitelist:      resourceGetStringList(data.Get("ingress_whitelist")),
			Name:                  data.Get("name").(string),
			RequireProbe:          data.Get("require_probe").(bool),
			UniqueIngresses:       data.Get("unique_ingress").(bool),
			UniqueServiceSelector: data.Get("unique_service_selector").(bool),
			ID:                    id,
		}
		params := opa_profiles.NewOpaProfilesUpdateParams().WithV(ApiVersion).WithBody(body)
		_, err = apiClient.client.OpaProfiles.OpaProfilesUpdate(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return readAfterUpdateWithRetries(generateResourceTaikunOPAProfileRead(true), ctx, data, meta)
}

func resourceTaikunOPAProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := opa_profiles.NewOpaProfilesDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteOpaProfileCommand{ID: id})
	_, err = apiClient.client.OpaProfiles.OpaProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func flattenTaikunOPAProfile(rawOPAProfile *models.OpaProfileListDto) map[string]interface{} {

	return map[string]interface{}{
		"allowed_repos":           rawOPAProfile.AllowedRepo,
		"forbid_node_port":        rawOPAProfile.ForbidNodePort,
		"forbid_http_ingress":     rawOPAProfile.ForbidHTTPIngress,
		"forbidden_tags":          rawOPAProfile.ForbidSpecificTags,
		"id":                      i32toa(rawOPAProfile.ID),
		"ingress_whitelist":       rawOPAProfile.IngressWhitelist,
		"is_default":              rawOPAProfile.IsDefault,
		"lock":                    rawOPAProfile.IsLocked,
		"name":                    rawOPAProfile.Name,
		"organization_id":         i32toa(rawOPAProfile.OrganizationID),
		"organization_name":       rawOPAProfile.OrganizationName,
		"require_probe":           rawOPAProfile.RequireProbe,
		"unique_ingress":          rawOPAProfile.UniqueIngresses,
		"unique_service_selector": rawOPAProfile.UniqueServiceSelector,
	}
}
