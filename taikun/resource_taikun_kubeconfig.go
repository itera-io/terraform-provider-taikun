package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/kube_config"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunKubeconfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_scope": {
			Description: "Who can use the kubeconfig: `personal` (only you), `managers` (managers only) or `all` (all users with access to this project).",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"all",
				"managers",
				"personal",
			}, false),
		},
		"id": {
			Description: "The kubeconfig's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The kubeconfig's name.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Description:      "ID of the kubeconfig's project.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"project_name": {
			Description: "Name of the kubeconfig's project.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"role": {
			Description: "The kubeconfig's role: `cluster-admin`, `admin`, `edit` or `view`.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"cluster-admin",
				"admin",
				"edit",
				"view",
			}, false),
		},
		"user_id": {
			Description: "ID of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"user_name": {
			Description: "Name of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"user_role": {
			Description: "Role of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunKubeconfig() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Kubeconfig",
		CreateContext: resourceTaikunKubeconfigCreate,
		ReadContext:   generateResourceTaikunKubeconfigReadWithoutRetries(),
		DeleteContext: resourceTaikunKubeconfigDelete,
		Schema:        resourceTaikunKubeconfigSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunKubeconfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateKubeConfigCommand{
		IsAccessibleForAll:     d.Get("access_scope").(string) == "all",
		IsAccessibleForManager: d.Get("access_scope").(string) == "managers",
		KubeConfigRoleID:       getKubeconfigRoleID(d.Get("role").(string)),
		Name:                   d.Get("name").(string),
	}
	projectID, err := atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	body.ProjectID = projectID

	params := kube_config.NewKubeConfigCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.KubeConfig.KubeConfigCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(response.Payload.ID)

	return readAfterCreateWithRetries(generateResourceTaikunKubeconfigReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunKubeconfigReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubeconfigRead(true)
}
func generateResourceTaikunKubeconfigReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubeconfigRead(false)
}
func generateResourceTaikunKubeconfigRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id := d.Id()
		id32, err := atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		params := kube_config.NewKubeConfigListParams().WithV(ApiVersion).WithID(&id32)
		response, err := apiClient.client.KubeConfig.KubeConfigList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		kubeconfigDTO := response.Payload.Data[0]
		if err := setResourceDataFromMap(d, flattenTaikunKubeconfig(kubeconfigDTO)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(kubeconfigDTO.ID))

		return nil
	}
}

func resourceTaikunKubeconfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteKubeConfigCommand{
		ID: id,
	}
	params := kube_config.NewKubeConfigDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, err := apiClient.client.KubeConfig.KubeConfigDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunKubeconfig(kubeconfigDTO *models.KubeConfigForUserDto) map[string]interface{} {
	kubeconfigMap := map[string]interface{}{
		"id":           i32toa(kubeconfigDTO.ID),
		"name":         kubeconfigDTO.ServiceAccountName,
		"project_id":   i32toa(kubeconfigDTO.ProjectID),
		"project_name": kubeconfigDTO.ProjectName,
		"user_id":      kubeconfigDTO.UserID,
		"user_name":    kubeconfigDTO.UserName,
		"user_role":    kubeconfigDTO.UserRole,
	}
	if kubeconfigDTO.IsAccessibleForAll {
		kubeconfigMap["access_scope"] = "all"
	} else if kubeconfigDTO.IsAccessibleForManager {
		kubeconfigMap["access_scope"] = "managers"
	} else {
		kubeconfigMap["access_scope"] = "personal"
	}
	return kubeconfigMap
}
