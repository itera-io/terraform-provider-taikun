package taikun

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/user_projects"
	"github.com/itera-io/taikungoclient/client/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunProjectUserAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"user_id": {
			Description:  "ID of the user.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceTaikunProjectUserAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project-User Attachment",
		CreateContext: resourceTaikunProjectUserAttachmentCreate,
		ReadContext:   generateResourceTaikunProjectUserAttachmentReadWithoutRetries(),
		DeleteContext: resourceTaikunProjectUserAttachmentDelete,
		Schema:        resourceTaikunProjectUserAttachmentSchema(),
	}
}

func resourceTaikunProjectUserAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*taikungoclient.Client)

	userId := d.Get("user_id").(string)

	projectId, err := atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("project_id isn't valid: %s", d.Get("project_id").(string))
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
	d.SetId(id)

	return readAfterCreateWithRetries(generateResourceTaikunProjectUserAttachmentReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunProjectUserAttachmentReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectUserAttachmentRead(true)
}
func generateResourceTaikunProjectUserAttachmentReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectUserAttachmentRead(false)
}
func generateResourceTaikunProjectUserAttachmentRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)

		id := d.Id()
		d.SetId("")
		projectId, userId, err := parseProjectUserAttachmentId(id)
		if err != nil {
			return diag.Errorf("Error while reading taikun_project_user_attachment : %s", err)
		}

		params := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
		response, err := apiClient.Client.Users.UsersList(params, apiClient)
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

		rawUser := response.GetPayload().Data[0]

		for _, e := range rawUser.BoundProjects {
			if e.ProjectID == projectId {
				if err := d.Set("project_id", i32toa(e.ProjectID)); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("project_name", e.ProjectName); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("user_id", rawUser.ID); err != nil {
					return diag.FromErr(err)
				}
				d.SetId(id)
				return nil
			}
		}

		if withRetries {
			d.SetId(id)
			return diag.Errorf(notFoundAfterCreateOrUpdateError)
		}
		return nil
	}
}

func resourceTaikunProjectUserAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	projectId, userId, err := parseProjectUserAttachmentId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_project_user_attachment : %s", err)
	}

	usersListParams := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
	usersListResponse, err := apiClient.Client.Users.UsersList(usersListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(usersListResponse.Payload.Data) != 1 {
		d.SetId("")
		return nil
	}

	projectsListParams := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&projectId)
	projectsListResponse, err := apiClient.Client.Projects.ProjectsList(projectsListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(projectsListResponse.Payload.Data) != 1 {
		d.SetId("")
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
	_, err = apiClient.Client.UserProjects.UserProjectsBindUsers(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
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
