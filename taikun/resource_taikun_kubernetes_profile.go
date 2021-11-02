package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunKubernetesProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"bastion_proxy": {
			Description: "Whether to expose the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"schedule_on_master": {
			Description: "When enabled, the workload will also run on master nodes (not recommended).",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"created_by": {
			Description: "The creator of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"cni": {
			Description: "Container Network Interface(CNI) of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Kubernetes profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"load_balancing_solution": {
			Description:  "The load-balancing solution: `None`, `Octavia` or `Taikun`. `Octavia` and `Taikun` are only available for OpenStack cloud.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "Octavia",
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"None", "Octavia", "Taikun"}, false),
		},
		"name": {
			Description:  "The name of the Kubernetes profile.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the Kubernetes profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Kubernetes Profile",
		CreateContext: resourceTaikunKubernetesProfileCreate,
		ReadContext:   generateResourceTaikunKubernetesProfileRead(false),
		UpdateContext: resourceTaikunKubernetesProfileUpdate,
		DeleteContext: resourceTaikunKubernetesProfileDelete,
		Schema:        resourceTaikunKubernetesProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunKubernetesProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	octaviaEnabled, taikunLBEnabled := parseLoadBalancingSolution(data.Get("load_balancing_solution").(string))
	body := &models.CreateKubernetesProfileCommand{
		Name:                    data.Get("name").(string),
		AllowSchedulingOnMaster: data.Get("schedule_on_master").(bool),
		TaikunLBEnabled:         taikunLBEnabled,
		OctaviaEnabled:          octaviaEnabled,
		ExposeNodePortOnBastion: data.Get("bastion_proxy").(bool),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := kubernetes_profiles.NewKubernetesProfilesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.KubernetesProfiles.KubernetesProfilesCreate(params, apiClient)
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
		lockBody := models.KubernetesProfilesLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := kubernetes_profiles.NewKubernetesProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.KubernetesProfiles.KubernetesProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunKubernetesProfileRead(true), ctx, data, meta)
}

func generateResourceTaikunKubernetesProfileRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(data.Id())
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.KubernetesProfiles.KubernetesProfilesList(kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&id), apiClient)
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

		rawKubernetesProfile := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunKubernetesProfile(rawKubernetesProfile))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunKubernetesProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("lock") {
		lockBody := models.KubernetesProfilesLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("lock").(bool)),
		}
		lockParams := kubernetes_profiles.NewKubernetesProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.KubernetesProfiles.KubernetesProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunKubernetesProfileRead(true), ctx, data, meta)
}

func resourceTaikunKubernetesProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := kubernetes_profiles.NewKubernetesProfilesDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.client.KubernetesProfiles.KubernetesProfilesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func flattenTaikunKubernetesProfile(rawKubernetesProfile *models.KubernetesProfilesListDto) map[string]interface{} {

	return map[string]interface{}{
		"bastion_proxy":           rawKubernetesProfile.ExposeNodePortOnBastion,
		"created_by":              rawKubernetesProfile.CreatedBy,
		"cni":                     rawKubernetesProfile.Cni,
		"id":                      i32toa(rawKubernetesProfile.ID),
		"lock":                    rawKubernetesProfile.IsLocked,
		"last_modified":           rawKubernetesProfile.LastModified,
		"last_modified_by":        rawKubernetesProfile.LastModifiedBy,
		"load_balancing_solution": getLoadBalancingSolution(rawKubernetesProfile.OctaviaEnabled, rawKubernetesProfile.TaikunLBEnabled),
		"name":                    rawKubernetesProfile.Name,
		"organization_id":         i32toa(rawKubernetesProfile.OrganizationID),
		"organization_name":       rawKubernetesProfile.OrganizationName,
		"schedule_on_master":      rawKubernetesProfile.AllowSchedulingOnMaster,
	}
}
