package policy_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			ValidateDiagFunc: utils.StringIsInt,
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

func ResourceTaikunPolicyProfile() *schema.Resource {
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
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateOpaProfileCommand{}
	body.SetAllowedRepo(utils.ResourceGetStringList(d.Get("allowed_repos").(*schema.Set).List()))
	body.SetForbidHttpIngress(d.Get("forbid_http_ingress").(bool))
	body.SetForbidNodePort(d.Get("forbid_node_port").(bool))
	body.SetForbidSpecificTags(utils.ResourceGetStringList(d.Get("forbidden_tags").(*schema.Set).List()))
	body.SetIngressWhitelist(utils.ResourceGetStringList(d.Get("ingress_whitelist").(*schema.Set).List()))
	body.SetName(d.Get("name").(string))
	body.SetRequireProbe(d.Get("require_probe").(bool))
	body.SetUniqueIngresses(d.Get("unique_ingress").(bool))
	body.SetUniqueServiceSelector(d.Get("unique_service_selector").(bool))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.OpaProfilesAPI.OpaprofilesCreate(context.TODO()).CreateOpaProfileCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(createResult.GetId())

	locked := d.Get("lock").(bool)
	if locked {
		id, err := utils.Atoi32(createResult.GetId())
		if err != nil {
			return diag.FromErr(err)
		}
		err = resourceTaikunPolicyProfileLock(id, true, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunPolicyProfileReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunPolicyProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunPolicyProfileRead(true)
}
func generateResourceTaikunPolicyProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunPolicyProfileRead(false)
}
func generateResourceTaikunPolicyProfileRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		rawPolicyProfile, err := ResourceTaikunPolicyProfileFind(id, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if rawPolicyProfile == nil {
			if isAfterUpdateOrCreate {
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		err = utils.SetResourceDataFromMap(d, flattenTaikunPolicyProfile(rawPolicyProfile))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunPolicyProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
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

		body := tkcore.OpaProfileUpdateCommand{}
		body.SetAllowedRepo(utils.ResourceGetStringList(d.Get("allowed_repos").(*schema.Set).List()))
		body.SetForbidHttpIngress(d.Get("forbid_http_ingress").(bool))
		body.SetForbidNodePort(d.Get("forbid_node_port").(bool))
		body.SetForbidSpecificTags(utils.ResourceGetStringList(d.Get("forbidden_tags").(*schema.Set).List()))
		body.SetIngressWhitelist(utils.ResourceGetStringList(d.Get("ingress_whitelist").(*schema.Set).List()))
		body.SetName(d.Get("name").(string))
		body.SetRequireProbe(d.Get("require_probe").(bool))
		body.SetUniqueIngresses(d.Get("unique_ingress").(bool))
		body.SetUniqueServiceSelector(d.Get("unique_service_selector").(bool))
		body.SetId(id)

		_, res, err := apiClient.Client.OpaProfilesAPI.OpaprofilesUpdate(context.TODO()).OpaProfileUpdateCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

	}

	if d.Get("lock").(bool) {
		err := resourceTaikunPolicyProfileLock(id, true, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunPolicyProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunPolicyProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.OpaProfilesAPI.OpaprofilesDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func resourceTaikunPolicyProfileLock(id int32, lock bool, apiClient *tk.Client) error {
	lockBody := tkcore.OpaProfileLockManagerCommand{}
	lockBody.SetId(id)
	lockBody.SetMode(utils.GetLockMode(lock))

	_, res, err := apiClient.Client.OpaProfilesAPI.OpaprofilesLockManager(context.TODO()).OpaProfileLockManagerCommand(lockBody).Execute()

	return tk.CreateError(res, err)
}

func flattenTaikunPolicyProfile(rawPolicyProfile *tkcore.OpaProfileListDto) map[string]interface{} {

	return map[string]interface{}{
		"allowed_repos":           rawPolicyProfile.GetAllowedRepo(),
		"forbid_node_port":        rawPolicyProfile.GetForbidNodePort(),
		"forbid_http_ingress":     rawPolicyProfile.GetForbidHttpIngress(),
		"forbidden_tags":          rawPolicyProfile.GetForbidSpecificTags(),
		"id":                      utils.I32toa(rawPolicyProfile.GetId()),
		"ingress_whitelist":       rawPolicyProfile.GetIngressWhitelist(),
		"is_default":              rawPolicyProfile.GetIsDefault(),
		"lock":                    rawPolicyProfile.GetIsLocked(),
		"name":                    rawPolicyProfile.GetName(),
		"organization_id":         utils.I32toa(rawPolicyProfile.GetOrganizationId()),
		"organization_name":       rawPolicyProfile.GetOrganizationName(),
		"require_probe":           rawPolicyProfile.GetRequireProbe(),
		"unique_ingress":          rawPolicyProfile.GetUniqueIngresses(),
		"unique_service_selector": rawPolicyProfile.GetUniqueServiceSelector(),
	}
}

func ResourceTaikunPolicyProfileFind(id int32, apiClient *tk.Client) (*tkcore.OpaProfileListDto, error) {
	params := apiClient.Client.OpaProfilesAPI.OpaprofilesList(context.TODO())
	var offset int32 = 0

	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return nil, tk.CreateError(res, err)
		}

		for _, policyProfile := range response.GetData() {
			if policyProfile.GetId() == id {
				return &policyProfile, nil
			}
		}

		offset += int32(len(response.GetData()))
		if offset == response.GetTotalCount() {
			break
		}
	}

	return nil, nil
}
