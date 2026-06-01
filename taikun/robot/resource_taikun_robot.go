package robot

import (
	"context"
	"time"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTaikunRobotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Robot user's name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"account_id": {
			Description: "Account ID the robot user belongs to.",
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
		},
		"organization_id": {
			Description: "Organization ID for the robot user.",
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Description: "Robot user's description.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"expires_at": {
			Description: "Expiration date in RFC3339 format.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"scopes": {
			Description: "List of scopes assigned to the robot user.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ips": {
			Description: "List of allowed IP addresses.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"regenerate_trigger": {
			Description: "Change this value to trigger credential regeneration.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"access_key": {
			Description: "Robot user's access key.",
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
		"secret_key": {
			Description: "Robot user's secret key. Only available at creation or after regeneration.",
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
		},
		"user_id": {
			Description: "Robot user's UUID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_active": {
			Description: "Whether the robot user is active.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_at": {
			Description: "Time and date of creation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func ResourceTaikunRobot() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Robot User",
		CreateContext: resourceTaikunRobotCreate,
		ReadContext:   generateResourceTaikunRobotReadWithoutRetries(),
		UpdateContext: resourceTaikunRobotUpdate,
		DeleteContext: resourceTaikunRobotDelete,
		Schema:        resourceTaikunRobotSchema(),
	}
}

func resourceTaikunRobotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateRobotUserCommand{}
	body.SetAccountId(int32(d.Get("account_id").(int)))
	body.SetName(d.Get("name").(string))
	body.SetDescription(d.Get("description").(string))

	expiresAt, err := time.Parse(time.RFC3339, d.Get("expires_at").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	body.SetExpiresAt(expiresAt)

	if v, ok := d.GetOk("organization_id"); ok {
		body.SetOrganizationId(int32(v.(int)))
	}

	if v, ok := d.GetOk("scopes"); ok {
		scopes := make([]string, len(v.([]interface{})))
		for i, s := range v.([]interface{}) {
			scopes[i] = s.(string)
		}
		body.SetScopes(scopes)
	}

	if v, ok := d.GetOk("ips"); ok {
		ips := make([]string, len(v.([]interface{})))
		for i, s := range v.([]interface{}) {
			ips[i] = s.(string)
		}
		body.SetIps(ips)
	}

	response, res, err := apiClient.Client.RobotAPI.RobotCreate(ctx).CreateRobotUserCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	_ = d.Set("access_key", response.GetAccessKey())
	_ = d.Set("secret_key", response.GetSecretKey())

	robot, diags := findRobotByAccessKey(ctx, apiClient, int32(d.Get("account_id").(int)), response.GetAccessKey())
	if diags != nil {
		return diags
	}

	d.SetId(robot.GetUserId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunRobotReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunRobotReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunRobotRead(true)
}

func generateResourceTaikunRobotReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunRobotRead(false)
}

func generateResourceTaikunRobotRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		robot, diags := findRobotByUserId(ctx, apiClient, d.Id())
		if diags != nil {
			return diags
		}

		if robot == nil {
			if withRetries {
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			d.SetId("")
			return nil
		}

		err := utils.SetResourceDataFromMap(d, flattenTaikunRobot(robot))
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceTaikunRobotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	if d.HasChanges("name", "description", "expires_at", "ips") {
		body := tkcore.EditRobotUserCommand{}
		body.SetId(d.Id())
		body.SetName(d.Get("name").(string))
		body.SetDescription(d.Get("description").(string))

		expiresAt, err := time.Parse(time.RFC3339, d.Get("expires_at").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetExpiresAt(expiresAt)

		if v, ok := d.GetOk("ips"); ok {
			ips := make([]string, len(v.([]interface{})))
			for i, s := range v.([]interface{}) {
				ips[i] = s.(string)
			}
			body.SetIps(ips)
		} else {
			body.SetIps([]string{})
		}

		res, err := apiClient.Client.RobotAPI.RobotUpdate(ctx).EditRobotUserCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChange("scopes") {
		body := tkcore.UpdateRobotScopeCommand{}
		body.SetId(d.Id())

		if v, ok := d.GetOk("scopes"); ok {
			scopes := make([]string, len(v.([]interface{})))
			for i, s := range v.([]interface{}) {
				scopes[i] = s.(string)
			}
			body.SetScopes(scopes)
		} else {
			body.SetScopes([]string{})
		}

		res, err := apiClient.Client.RobotAPI.RobotUpdateScope(ctx).UpdateRobotScopeCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
	}

	if d.HasChange("regenerate_trigger") {
		body := tkcore.RegenerateRobotTokenCommand{}
		body.SetId(d.Id())

		expiresAt, err := time.Parse(time.RFC3339, d.Get("expires_at").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetExpiresAt(expiresAt)

		response, res, err := apiClient.Client.RobotAPI.RobotRegenerate(ctx).RegenerateRobotTokenCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		_ = d.Set("access_key", response.GetAccessKey())
		_ = d.Set("secret_key", response.GetSecretKey())
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunRobotReadWithRetries(), ctx, d, meta)
}

func resourceTaikunRobotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	res, err := apiClient.Client.RobotAPI.RobotDelete(ctx, d.Id()).Execute()
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func findRobotByUserId(ctx context.Context, apiClient *tk.Client, userId string) (*tkcore.RobotUsersListDto, diag.Diagnostics) {
	response, res, err := apiClient.Client.RobotAPI.RobotList(ctx).SearchId(userId).Execute()
	if err != nil {
		return nil, diag.FromErr(tk.CreateError(res, err))
	}

	for i := range response.GetData() {
		item := response.GetData()[i]
		if item.GetUserId() == userId {
			return &item, nil
		}
	}

	return nil, nil
}

func findRobotByAccessKey(ctx context.Context, apiClient *tk.Client, accountId int32, accessKey string) (*tkcore.RobotUsersListDto, diag.Diagnostics) {
	response, res, err := apiClient.Client.RobotAPI.RobotList(ctx).AccountId(accountId).Search(accessKey).Execute()
	if err != nil {
		return nil, diag.FromErr(tk.CreateError(res, err))
	}

	for i := range response.GetData() {
		item := response.GetData()[i]
		if item.GetAccessKey() == accessKey {
			return &item, nil
		}
	}

	return nil, diag.Errorf("robot user with access key not found after creation")
}

func flattenTaikunRobot(robot *tkcore.RobotUsersListDto) map[string]interface{} {
	result := map[string]interface{}{
		"name":        robot.GetName(),
		"account_id":  robot.GetAccountId(),
		"description": robot.GetDescription(),
		"user_id":     robot.GetUserId(),
		"is_active":   robot.GetIsActive(),
		"created_at":  robot.GetCreatedAt(),
		"scopes":      robot.GetScopes(),
	}

	if robot.HasOrganizationId() {
		result["organization_id"] = robot.GetOrganizationId()
	}

	if robot.HasExpiresAt() {
		result["expires_at"] = robot.GetExpiresAt()
	}

	if ips := robot.GetIps(); len(ips) > 0 {
		result["ips"] = ips
	}

	return result
}
