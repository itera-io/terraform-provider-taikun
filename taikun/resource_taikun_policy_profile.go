package taikun

import (
	"context"
	"regexp"

	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/opa_profiles"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunPolicyProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_repos": {
			Description: "Requires container images to begin with a string from the specified list.",
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z]([-\da-zA-Z]*[\da-zA-Z])?(\.[a-zA-Z]([-\da-zA-Z]*[\da-zA-Z])?)*(:\d{1,5})?/)?([a-z\d]+((\.|(_{1,2}|-+))[a-z\d]+)*)(/([a-z\d]+((\.|(_{1,2}|-+))[a-z\d]+)*))*/?$`),
					"Please specify valid Docker image prefix",
				),
			},
		},
		"forbid_http_ingress": {
			Description: "Requires Ingress resources to be HTTPS only.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"forbid_node_port": {
			Description: "Disallows all Services with type NodePort.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"forbidden_tags": {
			Description: "Container images must have an image tag different from the ones in the list.",
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.All(
					validation.StringMatch(
						regexp.MustCompile("^[a-z0-9_][a-z0-9_.-]*$"),
						"Please specify valid Docker image tag",
					),
					validation.StringLenBetween(0, 128),
				),
			},
		},
		"id": {
			Description: "The ID of the Policy profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"ingress_whitelist": {
			Description: "List of allowed Ingress rule hosts.",
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: validation.All(
					validation.StringMatch(
						regexp.MustCompile(`^(\*|([a-zA-Z]([-\da-zA-Z]*[\da-zA-Z])?))(\.[a-zA-Z]([-\da-zA-Z]*[\da-zA-Z])?)+$`),
						"Please specify valid ingress domain",
					),
					validation.StringLenBetween(0, 253),
				),
			},
		},
		"is_default": {
			Description: "Indicates whether the Policy Profile is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Policy profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:  "The name of the Policy profile.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(2, 50),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the Policy profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Policy profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"require_probe": {
			Description: "Requires Pods to have readiness and liveness probes.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"unique_ingress": {
			Description: "Requires all Ingress rule hosts to be unique.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"unique_service_selector": {
			Description: "Whether services must have globally unique service selectors or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func resourceTaikunPolicyProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Policy Profile",
		CreateContext: resourceTaikunPolicyProfileCreate,
		ReadContext:   generateResourceTaikunPolicyProfileReadWithoutRetries(),
		UpdateContext: resourceTaikunPolicyProfileUpdate,
		DeleteContext: resourceTaikunPolicyProfileDelete,
		Schema:        resourceTaikunPolicyProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunPolicyProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateOpaProfileCommand{
		AllowedRepo:           resourceGetStringList(d.Get("allowed_repos").(*schema.Set).List()),
		ForbidHTTPIngress:     d.Get("forbid_http_ingress").(bool),
		ForbidNodePort:        d.Get("forbid_node_port").(bool),
		ForbidSpecificTags:    resourceGetStringList(d.Get("forbidden_tags").(*schema.Set).List()),
		IngressWhitelist:      resourceGetStringList(d.Get("ingress_whitelist").(*schema.Set).List()),
		Name:                  d.Get("name").(string),
		RequireProbe:          d.Get("require_probe").(bool),
		UniqueIngresses:       d.Get("unique_ingress").(bool),
		UniqueServiceSelector: d.Get("unique_service_selector").(bool),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := opa_profiles.NewOpaProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.OpaProfiles.OpaProfilesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	locked := d.Get("lock").(bool)
	if locked {
		id, err := atoi32(createResult.Payload.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = resourceTaikunPolicyProfileLock(id, true, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunPolicyProfileReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunPolicyProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunPolicyProfileRead(true)
}
func generateResourceTaikunPolicyProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunPolicyProfileRead(false)
}
func generateResourceTaikunPolicyProfileRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.OpaProfiles.OpaProfilesList(opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if isAfterUpdateOrCreate {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawPolicyProfile := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunPolicyProfile(rawPolicyProfile))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunPolicyProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		err := resourceTaikunPolicyProfileLock(id, false, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChangeExcept("lock") {

		body := &models.OpaProfileUpdateCommand{
			AllowedRepo:           resourceGetStringList(d.Get("allowed_repos").(*schema.Set).List()),
			ForbidHTTPIngress:     d.Get("forbid_http_ingress").(bool),
			ForbidNodePort:        d.Get("forbid_node_port").(bool),
			ForbidSpecificTags:    resourceGetStringList(d.Get("forbidden_tags").(*schema.Set).List()),
			IngressWhitelist:      resourceGetStringList(d.Get("ingress_whitelist").(*schema.Set).List()),
			Name:                  d.Get("name").(string),
			RequireProbe:          d.Get("require_probe").(bool),
			UniqueIngresses:       d.Get("unique_ingress").(bool),
			UniqueServiceSelector: d.Get("unique_service_selector").(bool),
			ID:                    id,
		}
		params := opa_profiles.NewOpaProfilesUpdateParams().WithV(ApiVersion).WithBody(body)
		_, err = apiClient.client.OpaProfiles.OpaProfilesUpdate(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if d.Get("lock").(bool) {
		err := resourceTaikunPolicyProfileLock(id, true, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunPolicyProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunPolicyProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := opa_profiles.NewOpaProfilesDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteOpaProfileCommand{ID: id})
	_, err = apiClient.client.OpaProfiles.OpaProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceTaikunPolicyProfileLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	lockBody := models.OpaProfileLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	lockParams := opa_profiles.NewOpaProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
	_, err := apiClient.client.OpaProfiles.OpaProfilesLockManager(lockParams, apiClient)

	return err
}

func flattenTaikunPolicyProfile(rawPolicyProfile *models.OpaProfileListDto) map[string]interface{} {

	return map[string]interface{}{
		"allowed_repos":           rawPolicyProfile.AllowedRepo,
		"forbid_node_port":        rawPolicyProfile.ForbidNodePort,
		"forbid_http_ingress":     rawPolicyProfile.ForbidHTTPIngress,
		"forbidden_tags":          rawPolicyProfile.ForbidSpecificTags,
		"id":                      i32toa(rawPolicyProfile.ID),
		"ingress_whitelist":       rawPolicyProfile.IngressWhitelist,
		"is_default":              rawPolicyProfile.IsDefault,
		"lock":                    rawPolicyProfile.IsLocked,
		"name":                    rawPolicyProfile.Name,
		"organization_id":         i32toa(rawPolicyProfile.OrganizationID),
		"organization_name":       rawPolicyProfile.OrganizationName,
		"require_probe":           rawPolicyProfile.RequireProbe,
		"unique_ingress":          rawPolicyProfile.UniqueIngresses,
		"unique_service_selector": rawPolicyProfile.UniqueServiceSelector,
	}
}
