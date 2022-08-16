package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/images"
)

func dataSourceTaikunImagesOpenStack() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given OpenStack cloud credential.",
		ReadContext: dataSourceTaikunImagesOpenStackRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "OpenStack cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"images": {
				Description: "List of retrieved OpenStack images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "OpenStack image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "OpenStack image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunImagesOpenStackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*taikungoclient.Client)
	params := images.NewImagesOpenstackImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)

	var imageList []map[string]interface{}
	for {
		response, err := apiClient.Client.Images.ImagesOpenstackImages(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		imageList = append(imageList, flattenTaikunImages(response.Payload.Data...)...)
		if len(imageList) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(imageList))
		params = params.WithOffset(&offset)
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}
