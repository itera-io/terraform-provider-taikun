package robot

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunRobots() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Robot Users Data Source",
		ReadContext: dataSourceTaikunRobotsRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account ID to list robot users for.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"organization_id": {
				Description: "Filter by organization ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"robots": {
				Description: "List of robot users.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Description: "Robot user's UUID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Robot user's name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Robot user's description.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"account_id": {
							Description: "Account ID.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"organization_id": {
							Description: "Organization ID.",
							Type:        schema.TypeInt,
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
						"expires_at": {
							Description: "Expiration date.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"scopes": {
							Description: "List of scopes.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunRobotsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	accountId := int32(d.Get("account_id").(int))
	robots := make([]tkcore.RobotUsersListDto, 0)
	var offset int32 = 0

	for {
		request := apiClient.Client.RobotAPI.RobotList(ctx).AccountId(accountId).Offset(offset)

		if v, ok := d.GetOk("organization_id"); ok {
			request = request.OrganizationId(int32(v.(int)))
		}

		response, res, err := request.Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		robots = append(robots, response.GetData()...)

		if !response.GetHasMore() {
			break
		}
		offset = response.GetNextOffset()
	}

	err := d.Set("robots", flattenRobots(robots))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("robots")

	return nil
}

func flattenRobots(robots []tkcore.RobotUsersListDto) []interface{} {
	result := make([]interface{}, len(robots))
	for i, robot := range robots {
		result[i] = map[string]interface{}{
			"user_id":         robot.GetUserId(),
			"name":            robot.GetName(),
			"description":     robot.GetDescription(),
			"account_id":      robot.GetAccountId(),
			"organization_id": robot.GetOrganizationId(),
			"is_active":       robot.GetIsActive(),
			"created_at":      robot.GetCreatedAt(),
			"expires_at":      robot.GetExpiresAt(),
			"scopes":          robot.GetScopes(),
		}
	}
	return result
}
