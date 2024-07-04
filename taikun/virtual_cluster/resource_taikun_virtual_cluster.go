package virtual_cluster

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
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
			Description:  "The hostname that will be used for the virtual cluster. If left empty, you are assigned a hostname based on your IP an virtual cluster name.",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func ResourceTaikunVirtualCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Virtual Cluster project in Taikun.",
		CreateContext: resourceTaikunVirtualClusterCreate,
		ReadContext:   generateResourceTaikunVirtualClusterRead(),
		UpdateContext: resourceTaikunVirtualClusterUpdate,
		DeleteContext: resourceTaikunVirtualClusterDelete,
		Schema:        resourceTaikunVirtualClusterSchema(),
	}
}

func resourceTaikunVirtualClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunVirtualClusterRead(), ctx, d, meta)
}

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
	parentId, err := utils.Atoi32(d.Get("parentId").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	bodyCreate.SetName(name)
	bodyCreate.SetProjectId(parentId)
	bodyCreate.SetAlertingProfileId(parentId)
	bodyCreate.SetDeleteOnExpiration(false)
	bodyCreate.SetExposeHostname(d.Get("hostname").(string))

	response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterCreate(ctx).CreateVirtualClusterCommand(*bodyCreate).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunVirtualClusterRead(), ctx, d, meta)
}

func generateResourceTaikunVirtualClusterRead() schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		id, err := utils.Atoi32(d.Get("id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		parentId, err := utils.Atoi32(d.Get("parentId").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		data, response, err := apiClient.Client.VirtualClusterAPI.VirtualClusterList(ctx, parentId).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}

		// Load all the found data to the local object
		err = utils.SetResourceDataFromMap(d, flattenTaikunVirtualCluster(&data.GetData()[0]))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("id").(string)) // We need to tell provider that object was created

		return nil
	}
}

func flattenTaikunVirtualCluster(rawVirtualProject *tkcore.VClusterListDto) map[string]interface{} {
	return map[string]interface{}{
		"id":       rawVirtualProject.GetId(),
		"name":     rawVirtualProject.GetName(),
		"hostname": rawVirtualProject.GetAccessIp(),
	}
}
