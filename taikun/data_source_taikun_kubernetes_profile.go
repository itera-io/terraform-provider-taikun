package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunKubernetesProfile() *schema.Resource {
	return &schema.Resource{
		Description: "Get a Kubernetes profile by its id.",
		ReadContext: dataSourceTaikunKubernetesProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "The id of the Kubernetes profile.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
			},
			"bastion_proxy_enabled": {
				Description: "Exposes the Service on each Node's IP at a static port, the NodePort. You'll be able to contact the NodePort Service, from outside the cluster, by requesting `<NodeIP>:<NodePort>`.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_by": {
				Description: "The creator of the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cni": {
				Description: "Container Network Interface(CNI) of the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_locked": {
				Description: "Indicates whether the Kubernetes profile is locked or not.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"last_modified": {
				Description: "Time of last modification.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The last user who modified the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"load_balancing_solution": {
				Description: "Load-balancing solution: `None`, `Octavia` or `Taikun`. `Octavia` and `Taikun` are only available for OpenStack cloud.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_id": {
				Description: "The id of the organization which owns the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_name": {
				Description: "The name of the organization which owns the Kubernetes profile.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceTaikunKubernetesProfileRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunKubernetesProfileRead(ctx, data, meta)
}
