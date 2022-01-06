package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunImages() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given cloud credential.",
		ReadContext: dataSourceTaikunImagesRead,
		Schema: map[string]*schema.Schema{
			"aws_owner": {
				Description: "AWS owner (only valid with AWS Cloud Credential ID).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"aws_platform": {
				Description: "AWS Platform (only valid with AWS Cloud Credential ID).",
				Type:        schema.TypeString,
				Optional:    true,
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

func dataSourceTaikunImagesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	cloudCredentialID, err := atoi32(data.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	//params := cloud_credentials.NewCloudCredentialsAllFlavorsParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)

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
		offer, offerIsSet := data.GetOk("azure_offer")
		publisher, publisherIsSet := data.GetOk("azure_publisher")
		SKU, SKUIsSet := data.GetOk("azure_sku")
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
		owner, ownerIsSet := data.GetOk("aws_owner")
		platform, platformIsSet := data.GetOk("aws_platform")
		if !ownerIsSet || !platformIsSet {
			return diag.Errorf("One of the following attributes is missing: aws_owner, aws_platform")
		}
		params := images.NewImagesAwsImagesParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		params.WithOwner(owner.(string)).WithPlatform(platform.(string))

		for {
			response, err := apiClient.client.Images.ImagesAwsImages(params, apiClient)
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

	if err := data.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(cloudCredentialID))
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
