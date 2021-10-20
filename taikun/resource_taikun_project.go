package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_profile_id": {
			Description:      "ID of the project's access profile",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"alerting_profile_id": {
			Description:      "ID of the project's alerting profile",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to store the project",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"kubernetes_profile_id": {
			Description:      "ID of the project's kubernetes profile",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"name": {
			Description:  "Project name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "ID of the organization which owns the project.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
	}
}
