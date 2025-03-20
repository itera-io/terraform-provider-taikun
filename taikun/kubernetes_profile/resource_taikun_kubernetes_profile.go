package kubernetes_profile

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTaikunKubernetesProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bastion_proxy": {
			Description: "Whether to expose the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"cni": {
			Description: "Container Network Interface (CNI) of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
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
		"lock": {
			Description: "Indicates whether to lock the Kubernetes profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
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
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Kubernetes profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"schedule_on_master": {
			Description: "When enabled, the workload will also run on master nodes (not recommended).",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"unique_cluster_name": {
			Description: "If not enabled, the cluster name will be cluster.local.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"nvidia_gpu_operator": {
			Description: "When enabled, the Kubernetes will have NVIDIA GPU support.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"wasm": {
			Description: "When enabled, the Kubernetes will have WASM support.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"proxmox_storage": {
			Description:      "Choose the default Proxmox storage.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: utils.IgnoreChangeFromEmpty,
		},
	}
}

func ResourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Kubernetes Profile",
		CreateContext: resourceTaikunKubernetesProfileCreate,
		ReadContext:   GenerateResourceTaikunKubernetesProfileReadWithoutRetries(),
		UpdateContext: resourceTaikunKubernetesProfileUpdate,
		DeleteContext: resourceTaikunKubernetesProfileDelete,
		Schema:        ResourceTaikunKubernetesProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunKubernetesProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	octaviaEnabled, taikunLBEnabled := utils.ParseLoadBalancingSolution(d.Get("load_balancing_solution").(string))
	body := tkcore.CreateKubernetesProfileCommand{}
	body.SetName(d.Get("name").(string))
	body.SetAllowSchedulingOnMaster(d.Get("schedule_on_master").(bool))
	body.SetTaikunLBEnabled(taikunLBEnabled)
	body.SetOctaviaEnabled(octaviaEnabled)
	body.SetExposeNodePortOnBastion(d.Get("bastion_proxy").(bool))
	body.SetNvidiaGpuOperatorEnabled(d.Get("nvidia_gpu_operator").(bool))
	body.SetUniqueClusterName(d.Get("unique_cluster_name").(bool))
	body.SetWasmEnabled(d.Get("wasm").(bool))

	proxmoxStorage, err := tkcore.NewProxmoxStorageFromValue(d.Get("proxmox_storage").(string))
	if err == nil {
		body.SetProxmoxStorage(*proxmoxStorage) // User input is not empty and valid.
	} else if d.Get("proxmox_storage").(string) != "" {
		return diag.FromErr(err) // User input is not empty and not valid.
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, res, err := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesCreate(context.TODO()).CreateKubernetesProfileCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	id, err := utils.Atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunKubernetesProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunKubernetesProfileReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunKubernetesProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubernetesProfileRead(true)
}
func GenerateResourceTaikunKubernetesProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubernetesProfileRead(false)
}
func generateResourceTaikunKubernetesProfileRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawKubernetesProfile := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, FlattenTaikunKubernetesProfile(&rawKubernetesProfile))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunKubernetesProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("lock") {
		if err := resourceTaikunKubernetesProfileLock(id, d.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunKubernetesProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunKubernetesProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesDelete(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func FlattenTaikunKubernetesProfile(rawKubernetesProfile *tkcore.KubernetesProfilesListDto) map[string]interface{} {

	return map[string]interface{}{
		"bastion_proxy":           rawKubernetesProfile.GetExposeNodePortOnBastion(),
		"created_by":              rawKubernetesProfile.GetCreatedBy(),
		"cni":                     rawKubernetesProfile.GetCni(),
		"id":                      utils.I32toa(rawKubernetesProfile.GetId()),
		"lock":                    rawKubernetesProfile.GetIsLocked(),
		"last_modified":           rawKubernetesProfile.GetLastModified(),
		"last_modified_by":        rawKubernetesProfile.GetLastModifiedBy(),
		"load_balancing_solution": utils.GetLoadBalancingSolution(rawKubernetesProfile.GetOctaviaEnabled(), rawKubernetesProfile.GetTaikunLBEnabled()),
		"name":                    rawKubernetesProfile.GetName(),
		"organization_id":         utils.I32toa(rawKubernetesProfile.GetOrganizationId()),
		"organization_name":       rawKubernetesProfile.GetOrganizationName(),
		"schedule_on_master":      rawKubernetesProfile.GetAllowSchedulingOnMaster(),
		"unique_cluster_name":     rawKubernetesProfile.GetUniqueClusterName(),
		"nvidia_gpu_operator":     rawKubernetesProfile.GetNvidiaGpuOperatorEnabled(),
		"wasm":                    rawKubernetesProfile.GetWasmEnabled(),
		"proxmox_storage":         rawKubernetesProfile.GetProxmoxStorage(),
	}
}

func resourceTaikunKubernetesProfileLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.KubernetesProfilesLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))

	_, res, err := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesLockManager(context.TODO()).KubernetesProfilesLockManagerCommand(body).Execute()
	return tk.CreateError(res, err)
}
