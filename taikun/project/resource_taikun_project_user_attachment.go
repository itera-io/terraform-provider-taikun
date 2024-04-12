package project

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunProjectUserAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Description:      "ID of the project.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
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

func ResourceTaikunProjectUserAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project-User Attachment",
		CreateContext: resourceTaikunProjectUserAttachmentCreate,
		ReadContext:   generateResourceTaikunProjectUserAttachmentReadWithoutRetries(),
		DeleteContext: resourceTaikunProjectUserAttachmentDelete,
		Schema:        resourceTaikunProjectUserAttachmentSchema(),
	}
}

func resourceTaikunProjectUserAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tk.Client)

	userId := d.Get("user_id").(string)

	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.Errorf("project_id isn't valid: %s", d.Get("project_id").(string))
	}

	body := tkcore.BindUsersCommand{
		Users: []tkcore.UpdateProjectUserDto{
			{
				IsBound: utils.BoolPtr(true),
				Id:      *tkcore.NewNullableString(&userId),
			},
		},
		ProjectId: &projectId,
	}
	res, err := client.Client.UserProjectsAPI.UserprojectsBindUsers(ctx).BindUsersCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	id := fmt.Sprintf("%d/%s", projectId, userId)
	d.SetId(id)

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunProjectUserAttachmentReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunProjectUserAttachmentReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectUserAttachmentRead(true)
}
func generateResourceTaikunProjectUserAttachmentReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectUserAttachmentRead(false)
}
func generateResourceTaikunProjectUserAttachmentRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		id := d.Id()
		d.SetId("")
		projectId, userId, err := parseProjectUserAttachmentId(id)
		if err != nil {
			return diag.Errorf("Error while reading taikun_project_user_attachment : %s", err)
		}

		response, _, err := apiClient.Client.UsersAPI.UsersList(ctx).Id(userId).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawUser := response.GetData()[0]

		for _, e := range rawUser.BoundProjects {
			if e.GetProjectId() == projectId {
				if err := d.Set("project_id", utils.I32toa(e.GetProjectId())); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("project_name", e.GetProjectName()); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("user_id", rawUser.GetId()); err != nil {
					return diag.FromErr(err)
				}
				d.SetId(id)
				return nil
			}
		}

		if withRetries {
			d.SetId(id)
			return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
		}
		return nil
	}
}

func resourceTaikunProjectUserAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	projectId, userId, err := parseProjectUserAttachmentId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_project_user_attachment : %s", err)
	}

	usersListResponse, res, err := apiClient.Client.UsersAPI.UsersList(context.TODO()).Id(userId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(usersListResponse.GetData()) != 1 {
		d.SetId("")
		return nil
	}

	projectsListResponse, res, err := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(projectId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(projectsListResponse.GetData()) != 1 {
		d.SetId("")
		return nil
	}

	body := tkcore.BindUsersCommand{
		Users: []tkcore.UpdateProjectUserDto{
			{
				IsBound: utils.BoolPtr(false),
				Id:      *tkcore.NewNullableString(&userId),
			},
		},
		ProjectId: &projectId,
	}
	res, err = apiClient.Client.UserProjectsAPI.UserprojectsBindUsers(context.TODO()).BindUsersCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func parseProjectUserAttachmentId(id string) (int32, string, error) {
	list := strings.Split(id, "/")
	if len(list) != 2 {
		return 0, "", fmt.Errorf("unable to determine taikun_project_user_attachment ID")
	}

	projectId, err := utils.Atoi32(list[0])
	if err != nil {
		return 0, "", fmt.Errorf("unable to determine taikun_project_user_attachment ID")
	}

	userId := list[1]

	return projectId, userId, nil
}
