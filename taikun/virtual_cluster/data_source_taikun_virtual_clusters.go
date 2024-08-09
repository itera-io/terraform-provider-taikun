package virtual_cluster

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func DataSourceTaikunVirtualClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Virtual Clusters.",
		ReadContext: dataSourceTaikunVirtualClustersRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
			"virtual_clusters": {
				Description: "List of retrieved Virtual Clusters.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunVirtualClusterSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunVirtualClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0
	var virtualClustersList []tkcore.VClusterListDto
	params := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := utils.Atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	// Iterate through all projects
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		// Is it a virtual cluster?
		for _, project := range response.GetData() {
			offset += 1
			if project.GetIsVirtualCluster() {
				// Save ID and parent project ID
				virtualProjectId := project.GetId()
				virtualProjectParentId := project.GetParentProjectId()
				// Get Virtual cluster details
				data, response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterList(ctx, virtualProjectParentId).Id(virtualProjectId).Execute()
				if err != nil {
					return diag.FromErr(tk.CreateError(response, err))
				}
				// Append virtual cluster
				virtualClustersList = append(virtualClustersList, data.GetData()[0])
			}
		}

		if offset == response.GetTotalCount() {
			break
		}
	}

	virtualClusters := make([]map[string]interface{}, len(virtualClustersList))
	for i, rawVirtualCluster := range virtualClustersList {
		virtualClusters[i] = flattenTaikunVirtualCluster(&rawVirtualCluster)
	}
	if err := d.Set("virtual_clusters", virtualClusters); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
