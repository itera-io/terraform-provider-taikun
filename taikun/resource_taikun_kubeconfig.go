package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
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
		"content": {
			Description: "Content of the kubeconfig's YAML file.",
			Type:        schema.TypeString,
			Sensitive:   true,
			Computed:    true,
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
		"namespace": {
			Description:  "The kubeconfig's namespace.",
			Type:         schema.TypeString,
			Optional:     true,
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
			Optional:    true,
			Computed:    true,
			ForceNew:    true,
		},
		"validity_period": {
			Description:  "The kubeconfig's validity period in minutes. Unlimited (-1) by default.",
			Type:         schema.TypeString,
			Optional:     true,
                        Default:      -1,
			ValidateFunc: validation.IntAtLeast(-1),
			ForceNew:     true,
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
	apiClient := meta.(*taikungoclient.Client)

	body := models.CreateKubeConfigCommand{
		IsAccessibleForAll:     d.Get("access_scope").(string) == "all",
		IsAccessibleForManager: d.Get("access_scope").(string) == "managers",
		KubeConfigRoleID:       getKubeconfigRoleID(d.Get("role").(string)),
		Name:                   d.Get("name").(string),
		TTL:                    int32(d.Get("validity_period").(int)),
	}

	if userId, userIdIsSet := d.GetOk("user_id"); userIdIsSet {
		body.UserID, _ = userId.(string)
	}

	if namespace, namespaceIsSet := d.GetOk("namespace"); namespaceIsSet {
		body.Namespace, _ = namespace.(string)
	}

	projectID, err := atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	body.ProjectID = projectID

	params := kube_config.NewKubeConfigCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.Client.KubeConfig.KubeConfigCreate(params, apiClient)
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
		apiClient := meta.(*taikungoclient.Client)
		id := d.Id()
		id32, err := atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		params := kube_config.NewKubeConfigListParams().WithV(ApiVersion).WithID(&id32)
		response, err := apiClient.Client.KubeConfig.KubeConfigList(params, apiClient)
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
		kubeconfigContent := resourceTaikunKubeconfigGetContent(
			kubeconfigDTO.ProjectID,
			kubeconfigDTO.ID,
			apiClient,
		)
		if err := setResourceDataFromMap(d, flattenTaikunKubeconfig(kubeconfigDTO, kubeconfigContent)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(kubeconfigDTO.ID))

		return nil
	}
}

func resourceTaikunKubeconfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteKubeConfigCommand{
		ID: id,
	}
	params := kube_config.NewKubeConfigDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, err := apiClient.Client.KubeConfig.KubeConfigDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunKubeconfig(kubeconfigDTO *models.KubeConfigForUserDto, kubeconfigContent string) map[string]interface{} {
	kubeconfigMap := map[string]interface{}{
		"content":      kubeconfigContent,
		"id":           i32toa(kubeconfigDTO.ID),
		"name":         kubeconfigDTO.DisplayName,
		"project_id":   i32toa(kubeconfigDTO.ProjectID),
		"project_name": kubeconfigDTO.ProjectName,
		"user_id":      kubeconfigDTO.UserID,
		"namespace":    kubeconfigDTO.Namespace,
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

func resourceTaikunKubeconfigGetContent(projectID int32, kubeconfigID int32, apiClient *taikungoclient.Client) string {

	body := models.DownloadKubeConfigCommand{
		ProjectID: projectID,
		ID:        kubeconfigID,
	}

	params := kube_config.NewKubeConfigDownloadParams().WithV(ApiVersion)
	params = params.WithBody(&body)

	response, err := apiClient.Client.KubeConfig.KubeConfigDownload(params, apiClient)
	if err != nil {
		return "Failed to retrieve content of kubeconfig"
	}

	return response.Payload.(string)
}
