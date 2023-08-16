package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DEPRECATED: this data source is deprecated in favour of `taikun_images_aws`, `taikun_images_azure`, `taikun_images_gcp` and `taikun_images_openstack`.

func dataSourceTaikunImages() *schema.Resource {
	return &schema.Resource{
		Description:        "Retrieve images for a given cloud credential.",
		DeprecationMessage: "This data source is deprecated in favour of taikun_images_aws, taikun_images_azure, taikun_images_gcp and taikun_images_openstack.",
		ReadContext:        dataSourceTaikunImagesRead,
		Schema: map[string]*schema.Schema{
			"aws_limit": {
				Description:  "Limit the number of listed AWS images (highly recommended as fetching the entire list of images can take a long time) (only valid with AWS cloud credential ID).",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"azure_offer": {
				Description: "Azure offer (only valid with Azure Cloud Credential ID).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"azure_publisher": {
				Description: "Azure publisher (only valid with Azure Cloud Credential ID).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"azure_sku": {
				Description: "Azure sku (only valid with Azure Cloud Credential ID).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cloud_credential_id": {
				Description:      "Cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"images": {
				Description: "List of retrieved images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"personal": {
				Description: "If the image is personal",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func dataSourceTaikunImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*tk.Client)
	list, res, err := apiClient.Client.CloudCredentialApi.CloudcredentialsDashboardList(ctx).Id(cloudCredentialID).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(list.GetAzure()) == 0 && len(list.GetAmazon()) == 0 && len(list.GetOpenstack()) == 0 && len(list.GetGoogle()) == 0 {
		return diag.Errorf("Cloud Credential not found")
	}

	var imageList []map[string]interface{}

	switch {
	case len(list.GetAzure()) != 0:
		offer, offerIsSet := d.GetOk("azure_offer")
		publisher, publisherIsSet := d.GetOk("azure_publisher")
		SKU, SKUIsSet := d.GetOk("azure_sku")
		if !SKUIsSet || !publisherIsSet || !offerIsSet {
			return diag.Errorf("All of the following attributes must be set: azure_offer, azure_publisher, azure_sku")
		}
		params := apiClient.Client.ImagesApi.ImagesAzureImages(ctx, cloudCredentialID, publisher.(string), offer.(string), SKU.(string)).Latest(false)
		var offset int32 = 0

		for {
			response, res, err := params.Offset(offset).Execute()
			if err != nil {
				return diag.FromErr(tk.CreateError(res, err))
			}
			imageList = append(imageList, flattenTaikunImages(response.GetData()...)...)
			if len(imageList) == int(response.GetTotalCount()) {
				break
			}
			offset = int32(len(imageList))
		}
	case len(list.GetAmazon()) != 0:
		var limit int32 = 0
		if limitData, limitIsSet := d.GetOk("aws_limit"); limitIsSet {
			limit = int32(limitData.(int))
		}

		params := apiClient.Client.ImagesApi.ImagesAwsCommonImages(ctx, cloudCredentialID)

		response, res, err := params.Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		imageList = flattenTaikunImagesAwsOwnerDetails(response)
		if limit != 0 && int32(len(imageList)) > limit {
			imageList = imageList[:limit]
		}
	default: // OpenStack
		var offset int32 = 0
		params := apiClient.Client.ImagesApi.ImagesOpenstackImages(ctx, cloudCredentialID).Personal(d.Get("personal").(bool)).Personal(false)

		for {
			response, res, err := params.Offset(offset).Execute()
			if err != nil {
				return diag.FromErr(tk.CreateError(res, err))
			}
			imageList = append(imageList, flattenTaikunImages(response.GetData()...)...)
			if len(imageList) == int(response.GetTotalCount()) {
				break
			}
			offset = int32(len(imageList))
		}
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}

func flattenTaikunImages(rawImages ...tkcore.CommonStringBasedDropdownDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.GetId(),
			"name": rawImage.GetName(),
		}
	}
	return images
}

func flattenTaikunImagesAwsOwnerDetails(rawImages []tkcore.AwsOwnerDetails) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.Image.GetId(),
			"name": rawImage.Image.GetName(),
		}
	}
	return images
}
