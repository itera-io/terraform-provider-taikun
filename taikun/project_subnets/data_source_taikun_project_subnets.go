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
						"name": {
							Description: "Subnet name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cidr": {
							Description: "Subnet CIDR.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"zone": {
							Description: "Availability zone.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"available_ip_count": {
							Description: "Number of available IP addresses.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"node_count": {
							Description: "Number of nodes.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"is_default": {
							Description: "Whether the subnet is default.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"owner_id": {
							Description: "Owner ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"vpc_id": {
							Description: "VPC ID.",
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
			"id":                 subnet.GetSubnetId(),
			"type":               string(subnet.GetSubnetType()),
			"name":               subnet.GetName(),
			"cidr":               subnet.GetCidr(),
			"zone":               subnet.GetZone(),
			"available_ip_count": int(subnet.GetAvailableIpCount()),
			"node_count":         int(subnet.GetNodeCount()),
			"is_default":         subnet.GetIsDefault(),
			"owner_id":           subnet.GetOwnerId(),
			"vpc_id":             subnet.GetVpcId(),
		}
	}

	err = d.Set("subnets", result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("project_id").(string))

	return nil
}
