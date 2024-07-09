package virtual_cluster

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"time"
)

func resourceTaikunVirtualClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the Virtual cluster project.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the virtual cluster.",
			Type:         schema.TypeString,
			ValidateFunc: validation.StringLenBetween(3, 30),
			Required:     true,
			ForceNew:     true,
		},
		"parent_id": {
			Description:  "The ID of the parent of the virtual cluster.",
			Type:         schema.TypeString,
			ValidateFunc: validation.StringLenBetween(3, 30),
			Required:     true,
			ForceNew:     true,
		},
		"hostname": {
			Description: "The hostname that will be used for the virtual cluster. If left empty, you are assigned a hostname based on your IP an virtual cluster name.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
		},
		"hostname_generated": {
			Description: "IP-based resolvable hostname generated by Taikun.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func ResourceTaikunVirtualCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Virtual Cluster project in Taikun.",
		CreateContext: resourceTaikunVirtualClusterCreate,
		ReadContext:   generateResourceTaikunVirtualClusterRead(),
		//UpdateContext: resourceTaikunVirtualClusterUpdate,
		DeleteContext: resourceTaikunVirtualClusterDelete,
		Schema:        resourceTaikunVirtualClusterSchema(),
	}
}

//func resourceTaikunVirtualClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunVirtualClusterRead(), ctx, d, meta)
//}

func resourceTaikunVirtualClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	deleteCommand := tkcore.DeleteVirtualClusterCommand{}
	virtualClusterId, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	deleteCommand.SetProjectId(virtualClusterId)
	response, err2 := apiClient.Client.VirtualClusterAPI.VirtualClusterDelete(ctx).DeleteVirtualClusterCommand(deleteCommand).Execute()
	if err2 != nil {
		return diag.FromErr(tk.CreateError(response, err2))
	}
	d.SetId("")
	return nil
}

func resourceTaikunVirtualClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	// Send create query
	bodyCreate := &tkcore.CreateVirtualClusterCommand{}
	name := d.Get("name").(string)
	parentId, err := utils.Atoi32(d.Get("parent_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	bodyCreate.SetName(name)
	bodyCreate.SetProjectId(parentId)
	bodyCreate.SetDeleteOnExpiration(false)
	bodyCreate.SetExposeHostname(d.Get("hostname").(string))

	response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterCreate(ctx).CreateVirtualClusterCommand(*bodyCreate).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	err = resourceTaikunVirtualClusterWaitForReady(name, parentId, ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunVirtualClusterRead(), ctx, d, meta)
}

func generateResourceTaikunVirtualClusterRead() schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		virtualClusterName := d.Get("name").(string)
		parentId, err := utils.Atoi32(d.Get("parent_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		data, response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterList(ctx, parentId).Search(virtualClusterName).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}

		foundMatch := false
		var rawVirtualProject tkcore.VClusterListDto
		for _, virtualProject := range data.GetData() {
			if virtualProject.GetName() == virtualClusterName {
				foundMatch = true
				rawVirtualProject = virtualProject
				break
			}
		}
		if !foundMatch {
			return diag.FromErr(fmt.Errorf("Created Virtual project not found in Taikun response."))
		}

		// Load all the found data to the local object
		err = utils.SetResourceDataFromMap(d, flattenTaikunVirtualCluster(&rawVirtualProject))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("id").(string)) // We need to tell provider that object was created

		return nil
	}
}

func flattenTaikunVirtualCluster(rawVirtualProject *tkcore.VClusterListDto) map[string]interface{} {
	return map[string]interface{}{
		"id":                 utils.I32toa(rawVirtualProject.GetId()),
		"name":               rawVirtualProject.GetName(),
		"hostname_generated": rawVirtualProject.GetAccessIp(),
	}
}

// After we create a virtual cluster, we wait. It takes some time (~60s).
func resourceTaikunVirtualClusterWaitForReady(virtualClusterName string, parentId int32, ctx context.Context, meta interface{}) error {
	apiClient := meta.(*tk.Client)

	pendingStates := []string{"pending"}
	targetStates := []string{"finished"}

	// Try to get the virtual project
	createStateConf := &retry.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			data, response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterList(ctx, parentId).Search(virtualClusterName).Execute()
			if err != nil {
				return nil, "", tk.CreateError(response, err)
			}

			foundMatch := "pending"
			for _, virtualProject := range data.GetData() {
				if virtualProject.GetName() == virtualClusterName {
					if virtualProject.GetStatus() == tkcore.PROJECTSTATUS_READY {
						foundMatch = "finished"
						break
					}
				}
			}

			return data, foundMatch, nil
		},
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for virtual cluster (%s) to be read: %s", virtualClusterName, err)
	}

	return nil
}
