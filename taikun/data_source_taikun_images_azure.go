package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/images"
)

func dataSourceTaikunImagesAzure() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given Azure cloud credential.",
		ReadContext: dataSourceTaikunImagesAzureRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "Azure cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"images": {
				Description: "List of retrieved Azure images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Azure image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Azure image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"latest": {
				Description: "Retrieve latest Azure images.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"offer": {
				Description: "Azure offer.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"publisher": {
				Description: "Azure publisher.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"sku": {
				Description: "Azure SKU.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func dataSourceTaikunImagesAzureRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	latest := d.Get("latest").(bool)

	params := images.NewImagesAzureImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
	params = params.WithPublisherName(d.Get("publisher").(string))
	params = params.WithOffer(d.Get("offer").(string))
	params = params.WithSku(d.Get("sku").(string))
	params = params.WithLatest(&latest)

	var imageList []map[string]interface{}
	for {
		response, err := apiClient.Client.Images.ImagesAzureImages(params, apiClient)
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
