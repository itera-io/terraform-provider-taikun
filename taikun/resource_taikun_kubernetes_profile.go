package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Kubernetes Profile",
		CreateContext: resourceTaikunKubernetesProfileCreate,
		ReadContext:   resourceTaikunKubernetesProfileRead,
		UpdateContext: resourceTaikunKubernetesProfileUpdate,
		DeleteContext: resourceTaikunKubernetesProfileDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The id of the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"bastion_proxy_enabled": {
				Description: "Exposes the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`.",
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
			"is_locked": {
				Description: "Indicates whether the Kubernetes profile is locked or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"last_modified": {
				Description: "Time of last modification.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The last user who modified the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"load_balancing_solution": {
				Description:  "Load-balancing solution.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Octavia",
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"None", "Octavia", "Taikun"}, false),
			},
			"name": {
				Description: "The name of the Kubernetes profile.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"organization_id": {
				Description:  "The id of the organization which owns the Kubernetes profile.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: stringIsInt,
			},
			"organization_name": {
				Description: "The name of the organization which owns the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
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
		TaikunLBEnabled:         taikunLBEnabled,
		OctaviaEnabled:          octaviaEnabled,
		ExposeNodePortOnBastion: data.Get("bastion_proxy_enabled").(bool),
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

	locked := data.Get("is_locked").(bool)
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

	data.SetId(createResult.Payload.ID)

	return resourceTaikunKubernetesProfileRead(ctx, data, meta)
}

func resourceTaikunKubernetesProfileRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if response.Payload.TotalCount == 1 {
		rawKubernetesProfile := response.GetPayload().Data[0]

		if err := data.Set("bastion_proxy_enabled", rawKubernetesProfile.ExposeNodePortOnBastion); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("created_by", rawKubernetesProfile.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("cni", rawKubernetesProfile.Cni); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawKubernetesProfile.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawKubernetesProfile.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawKubernetesProfile.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("load_balancing_solution", getLoadBalancingSolution(rawKubernetesProfile.OctaviaEnabled, rawKubernetesProfile.TaikunLBEnabled)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawKubernetesProfile.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawKubernetesProfile.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawKubernetesProfile.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawKubernetesProfile.OrganizationName); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunKubernetesProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("is_locked") {
		lockBody := models.KubernetesProfilesLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := kubernetes_profiles.NewKubernetesProfilesLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.KubernetesProfiles.KubernetesProfilesLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunKubernetesProfileRead(ctx, data, meta)
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
