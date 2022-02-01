package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunImages() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given cloud credential.",
		ReadContext: dataSourceTaikunImagesRead,
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
		},
	}
}

func dataSourceTaikunImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*apiClient)
	params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&cloudCredentialID)
	list, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(list.GetPayload().Azure) == 0 && len(list.GetPayload().Amazon) == 0 && len(list.GetPayload().Openstack) == 0 {
		return diag.Errorf("Cloud Credential not found")
	}

	var imageList []map[string]interface{}

	switch {
	case len(list.GetPayload().Azure) != 0:
		offer, offerIsSet := d.GetOk("azure_offer")
		publisher, publisherIsSet := d.GetOk("azure_publisher")
		SKU, SKUIsSet := d.GetOk("azure_sku")
		if !SKUIsSet || !publisherIsSet || !offerIsSet {
			return diag.Errorf("One of the following attributes is missing: azure_offer, azure_publisher, azure_sku")
		}
		params := images.NewImagesAzureImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		params.WithPublisherName(publisher.(string)).WithOffer(offer.(string)).WithSku(SKU.(string))

		for {
			response, err := apiClient.client.Images.ImagesAzureImages(params, apiClient)
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
	case len(list.GetPayload().Amazon) != 0:
		params := images.NewImagesAwsImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		var limit int32 = 0
		if limitData, limitIsSet := d.GetOk("aws_limit"); limitIsSet {
			limit = int32(limitData.(int))
		}
		for {
			response, err := apiClient.client.Images.ImagesAwsImages(params, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
			imageList = append(imageList, flattenTaikunImages(response.Payload.Data...)...)
			count := int32(len(imageList))
			if limit != 0 && count >= limit {
				break
			}
			if count == response.Payload.TotalCount {
				break
			}
			params = params.WithOffset(&count)
		}
		if limit != 0 && int32(len(imageList)) > limit {
			imageList = imageList[:limit]
		}
	default: // OpenStack
		params := images.NewImagesOpenstackImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)

		for {
			response, err := apiClient.client.Images.ImagesOpenstackImages(params, apiClient)
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
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}

func flattenTaikunImages(rawImages ...*models.CommonStringBasedDropdownDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.ID,
			"name": rawImage.Name,
		}
	}
	return images
}
