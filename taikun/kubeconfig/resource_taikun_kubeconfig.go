package kubeconfig

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			ValidateDiagFunc: utils.StringIsInt,
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
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      -1,
			ValidateFunc: validation.IntAtLeast(-1),
			ForceNew:     true,
		},
	}
}

func ResourceTaikunKubeconfig() *schema.Resource {
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
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateKubeConfigCommand{}
	body.SetIsAccessibleForAll(d.Get("access_scope").(string) == "all")
	body.SetIsAccessibleForManager(d.Get("access_scope").(string) == "managers")
	body.SetKubeConfigRoleId(utils.GetKubeconfigRoleID(d.Get("role").(string)))
	body.SetName(d.Get("name").(string))
	body.SetTtl(int32(d.Get("validity_period").(int)))

	if userId, userIdIsSet := d.GetOk("user_id"); userIdIsSet {
		body.SetUserId(userId.(string))
	}

	if namespace, namespaceIsSet := d.GetOk("namespace"); namespaceIsSet {
		body.SetNamespace(namespace.(string))
	}

	projectID, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	body.SetProjectId(projectID)

	response, res, err := apiClient.Client.KubeConfigAPI.KubeconfigCreate(ctx).CreateKubeConfigCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	d.SetId(response.GetId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunKubeconfigReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunKubeconfigReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubeconfigRead(true)
}
func generateResourceTaikunKubeconfigReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunKubeconfigRead(false)
}
func generateResourceTaikunKubeconfigRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id := d.Id()
		id32, err := utils.Atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		projectID, err := utils.Atoi32(d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.KubeConfigAPI.KubeconfigList(context.TODO()).Id(id32).ProjectId(projectID).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		kubeconfigDTO := response.Data[0]
		kubeconfigContent := resourceTaikunKubeconfigGetContent(
			kubeconfigDTO.GetProjectId(),
			kubeconfigDTO.GetId(),
			apiClient,
		)
		if err := utils.SetResourceDataFromMap(d, flattenTaikunKubeconfig(&kubeconfigDTO, kubeconfigContent)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(kubeconfigDTO.GetId()))

		return nil
	}
}

func resourceTaikunKubeconfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.DeleteKubeConfigCommand{}
	body.SetId(id)

	_, res, err := apiClient.Client.KubeConfigAPI.KubeconfigDelete(context.TODO()).DeleteKubeConfigCommand(body).Execute()

	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunKubeconfig(kubeconfigDTO *tkcore.KubeConfigForUserDto, kubeconfigContent string) map[string]interface{} {
	kubeconfigMap := map[string]interface{}{
		"content":      kubeconfigContent,
		"id":           utils.I32toa(kubeconfigDTO.GetId()),
		"name":         kubeconfigDTO.GetDisplayName(),
		"project_id":   utils.I32toa(kubeconfigDTO.GetProjectId()),
		"project_name": kubeconfigDTO.GetProjectName(),
		"user_id":      kubeconfigDTO.GetUserId(),
		"namespace":    kubeconfigDTO.GetNamespace(),
	}

	if kubeconfigDTO.GetIsAccessibleForAll() {
		kubeconfigMap["access_scope"] = "all"
	} else if kubeconfigDTO.GetIsAccessibleForManager() {
		kubeconfigMap["access_scope"] = "managers"
	} else {
		kubeconfigMap["access_scope"] = "personal"
	}
	return kubeconfigMap
}

func resourceTaikunKubeconfigGetContent(projectID int32, kubeconfigID int32, apiClient *tk.Client) string {

	body := tkcore.DownloadKubeConfigCommand{}
	body.SetProjectId(projectID)
	body.SetId(kubeconfigID)
	response, _, err := apiClient.Client.KubeConfigAPI.KubeconfigDownload(context.TODO()).DownloadKubeConfigCommand(body).Execute()

	if err != nil {
		return "Failed to retrieve content of kubeconfig"
	}

	return response
}
