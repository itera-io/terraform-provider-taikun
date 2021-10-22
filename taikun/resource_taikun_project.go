package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_profile_id": {
			Description:      "ID of the project's access profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"alerting_profile_id": {
			Description:      "ID of the project's alerting profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true, // TODO alerting profile can be detached, maybe handle in Update?
		},
		"auto_upgrades": {
			Description: "Kubespray version will be automatically upgraded if new version is available.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to store the project.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"expiration_date": {
			Description:      "Project's expiration date in the format: 'dd/mm/yyyy'.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsDate,
		},
		"id": {
			Description: "Project ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"kubernetes_profile_id": {
			Description:      "ID of the project's kubernetes profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"name": {
			Description:  "Project name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},
		"organization_id": {
			Description:      "ID of the organization which owns the project.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
	}
}

func resourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project",
		CreateContext: resourceTaikunProjectCreate,
		ReadContext:   resourceTaikunProjectRead,
		UpdateContext: resourceTaikunProjectUpdate,
		DeleteContext: resourceTaikunProjectDelete,
		Schema:        resourceTaikunProjectSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunProjectRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := data.Id()
	id32, err := atoi32(id)
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id32) // TODO use /api/v1/projects endpoint?
	response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
	if err != nil {
		return nil
	}

	projectDetailsDTO := response.Payload.Project
	err = setResourceDataFromMap(data, flattenTaikunProject(projectDetailsDTO))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(id)

	return nil
}

func resourceTaikunProjectCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateProjectCommand{
		Name: data.Get("name").(string),
	}
	body.CloudCredentialID, _ = atoi32(data.Get("cloud_credential_id").(string))
	if accessProfileID, accessProfileIDIsSet := data.GetOk("access_profile_id"); accessProfileIDIsSet {
		body.AccessProfileID, _ = atoi32(accessProfileID.(string))
	}
	if alertingProfileID, alertingProfileIDIsSet := data.GetOk("alerting_profile_id"); alertingProfileIDIsSet {
		body.AlertingProfileID, _ = atoi32(alertingProfileID.(string))
	}
	if autoUpgrades, autoUpgradesIsSet := data.GetOk("auto_upgrades"); autoUpgradesIsSet {
		body.IsAutoUpgrade = autoUpgrades.(bool)
	}
	if expirationDate, expirationDateIsSet := data.GetOk("expiration_date"); expirationDateIsSet {
		dateTime := dateToDateTime(expirationDate.(string))
		body.ExpiredAt = &dateTime
	} else {
		body.ExpiredAt = nil
	}
	if kubernetesProfileID, kubernetesProfileIDIsSet := data.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		body.KubernetesProfileID, _ = atoi32(kubernetesProfileID.(string))
	}
	if organizationID, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		body.OrganizationID, _ = atoi32(organizationID.(string))
	}

	params := projects.NewProjectsCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Projects.ProjectsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.Payload.ID)

	return resourceTaikunProjectRead(ctx, data, meta)
}

func resourceTaikunProjectUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("expiration_date") {
		body := models.ProjectExtendLifeTimeCommand{
			ProjectID: id,
		}
		if expirationDate, expirationDateIsSet := data.GetOk("expiration_date"); expirationDateIsSet {
			dateTime := dateToDateTime(expirationDate.(string))
			body.ExpireAt = &dateTime
		} else {
			body.ExpireAt = nil
		}

		params := projects.NewProjectsExtendLifeTimeParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.Projects.ProjectsExtendLifeTime(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	return resourceTaikunProjectRead(ctx, data, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteProjectCommand{ProjectID: id, IsForceDelete: false}
	params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.client.Projects.ProjectsDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

// TODO change type of DTO if read endpoint is modified
func flattenTaikunProject(projectDetailsDTO *models.ProjectDetailsForServersDto) map[string]interface{} {
	return map[string]interface{}{
		"access_profile_id":     i32toa(projectDetailsDTO.AccessProfileID),
		"alerting_profile_id":   i32toa(projectDetailsDTO.AlertingProfileID),
		"auto_upgrades":         projectDetailsDTO.IsAutoUpgrade,
		"cloud_credential_id":   i32toa(projectDetailsDTO.CloudID),
		"expiration_date":       rfc3339DateTimeToDate(projectDetailsDTO.ExpiredAt),
		"id":                    i32toa(projectDetailsDTO.ProjectID),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.KubernetesProfileID),
		"name":                  projectDetailsDTO.ProjectName,
		"organization_id":       i32toa(projectDetailsDTO.OrganizationID),
	}
}
