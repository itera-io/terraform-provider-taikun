package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunKubeconfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_scope": {
			Description: "Who can use the kubeconfig: `personal` (only you), `managers` (managers only) or `all` (all users with access to this project).",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"all",
				"managers",
				"personal",
			}, false),
		},
		"name": {
			Description:  "Kubeconfig's name.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Description:      "ID of the kubeconfig's project.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"project_name": {
			Description: "Name of the kubeconfig's project.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"role": {
			Description: "Kubeconfig's role.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"cluster-admin",
				"admin",
				"edit",
				"view",
			}, false),
		},
		"user_id": {
			Description: "ID of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"user_name": {
			Description: "Name of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"user_role": {
			Description: "Role of the kubeconfig's user, if the kubeconfig is personal.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}
