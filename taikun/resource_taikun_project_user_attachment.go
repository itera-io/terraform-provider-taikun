package taikun

import (
	"context"
	"fmt"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/user_projects"
	"github.com/itera-io/taikungoclient/client/users"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunProjectUserAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_id": {
			Description: "ID of the user.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"project_id": {
			Description:      "ID of the project.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"project_name": {
			Description: "Name of the project.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunProjectUserAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project - User Attachment",
		CreateContext: resourceTaikunProjectUserAttachmentCreate,
		ReadContext:   resourceTaikunProjectUserAttachmentRead,
		DeleteContext: resourceTaikunProjectUserAttachmentDelete,
		Schema:        resourceTaikunProjectUserAttachmentSchema(),
	}
}

func resourceTaikunProjectUserAttachmentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	userId := data.Get("user_id").(string)

	projectId, err := atoi32(data.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("project_id isn't valid: %s", data.Get("project_id").(string))
	}

	body := &models.BindUsersCommand{
		Users: []*models.UpdateProjectUserDto{
			{
				IsBound: true,
				UserID:  userId,
			},
		},
		ProjectID: projectId,
	}
	params := user_projects.NewUserProjectsBindUsersParams().WithV(ApiVersion).WithBody(body)
	_, err = client.client.UserProjects.UserProjectsBindUsers(params, client)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d/%s", projectId, userId)
	data.SetId(id)

	return resourceTaikunProjectUserAttachmentRead(ctx, data, meta)
}

func resourceTaikunProjectUserAttachmentRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id := data.Id()
	data.SetId("")
	projectId, userId, err := parseProjectUserAttachmentId(id)
	if err != nil {
		return diag.Errorf("Error while reading taikun_project_user_attachment : %s", err)
	}

	params := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
	response, err := apiClient.client.Users.UsersList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(response.Payload.Data) != 1 {
		return nil
	}

	rawUser := response.GetPayload().Data[0]

	for _, e := range rawUser.BoundProjects {
		if e.ProjectID == projectId {
			if err := data.Set("project_id", i32toa(e.ProjectID)); err != nil {
				return diag.FromErr(err)
			}
			if err := data.Set("project_name", e.ProjectName); err != nil {
				return diag.FromErr(err)
			}
			if err := data.Set("user_id", rawUser.ID); err != nil {
				return diag.FromErr(err)
			}
			data.SetId(id)
			break
		}
	}

	return nil
}

func resourceTaikunProjectUserAttachmentDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	projectId, userId, err := parseProjectUserAttachmentId(data.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_project_user_attachment : %s", err)
	}

	usersListParams := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
	usersListResponse, err := apiClient.client.Users.UsersList(usersListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(usersListResponse.Payload.Data) != 1 {
		data.SetId("")
		return nil
	}

	projectsListParams := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&projectId)
	projectsListResponse, err := apiClient.client.Projects.ProjectsList(projectsListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(projectsListResponse.Payload.Data) != 1 {
		data.SetId("")
		return nil
	}

	body := &models.BindUsersCommand{
		Users: []*models.UpdateProjectUserDto{
			{
				IsBound: false,
				UserID:  userId,
			},
		},
		ProjectID: projectId,
	}
	params := user_projects.NewUserProjectsBindUsersParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.client.UserProjects.UserProjectsBindUsers(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func parseProjectUserAttachmentId(id string) (int32, string, error) {
	list := strings.Split(id, "/")
	if len(list) != 2 {
		return 0, "", fmt.Errorf("unable to determine taikun_project_user_attachment ID")
	}

	projectId, err := atoi32(list[0])
	if err != nil {
		return 0, "", fmt.Errorf("unable to determine taikun_project_user_attachment ID")
	}

	userId := list[1]

	return projectId, userId, nil
}
