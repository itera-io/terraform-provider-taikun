package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTaikunProjectUserAttachment() *schema.Resource {
	return &schema.Resource{
		Description:        "Taikun Project User Attachment (Deprecated)",
		DeprecationMessage: "This resource has been removed from the Taikun API and is no longer functional.",
		CreateContext:      resourceTaikunProjectUserAttachmentCreate,
		ReadContext:        resourceTaikunProjectUserAttachmentRead,
		DeleteContext:      resourceTaikunProjectUserAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The ID of the project.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"user_id": {
				Description: "The ID of the user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceTaikunProjectUserAttachmentCreate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Errorf("The taikun_project_user_attachment resource is deprecated and no longer supported by the Taikun API.")
}

func resourceTaikunProjectUserAttachmentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Return nil to allow existing state to be read without error,
	// though it can't be refreshed from the API.
	return nil
}

func resourceTaikunProjectUserAttachmentDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Simply remove from state
	d.SetId("")
	return nil
}
