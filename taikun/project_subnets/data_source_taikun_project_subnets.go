package project_subnets

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunProjectSubnets() *schema.Resource {
	return &schema.Resource{
		Description: "Taikun Project Subnets Data Source",
		ReadContext: dataSourceTaikunProjectSubnetsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Project ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"subnets": {
				Description: "List of cloud subnets for the project.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Subnet ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Subnet type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunProjectSubnetsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	response, res, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	project := response.GetProject()
	subnets := project.GetCloudSubnets()
	result := make([]interface{}, len(subnets))
	for i, subnet := range subnets {
		result[i] = map[string]interface{}{
			"id":   subnet.GetSubnetId(),
			"type": string(subnet.GetSubnetType()),
		}
	}

	err = d.Set("subnets", result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("project_id").(string))

	return nil
}
