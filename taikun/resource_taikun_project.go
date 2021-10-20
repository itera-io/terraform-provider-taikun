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
			Description:      "ID of the project's access profile",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"alerting_profile_id": {
			Description:      "ID of the project's alerting profile",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true, // TODO alerting profile can be detached, maybe handle in Update?
		},
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to store the project",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"id": {
			Description: "Project ID",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"kubernetes_profile_id": {
			Description:      "ID of the project's kubernetes profile",
			Type:             schema.TypeString,
			Optional:         true,
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
		// UpdateContext: resourceTaikunProjectUpdate,
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

	params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id32)
	response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	projectDetailsDTO := response.Payload.Project
	if err := data.Set("access_profile_id", i32toa(projectDetailsDTO.AccessProfileID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("alerting_profile_id", i32toa(projectDetailsDTO.AlertingProfileID)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("cloud_credential_id", projectDetailsDTO.CloudID); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("id", id); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("kubernetes_profile_id", projectDetailsDTO.KubernetesProfileID); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("name", projectDetailsDTO.ProjectName); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("organization_id", projectDetailsDTO.OrganizationID); err != nil {
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
	// FIXME
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
