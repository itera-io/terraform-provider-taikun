package taikun

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/aws"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunImagesAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given AWS cloud credential.",
		ReadContext: dataSourceTaikunImagesAWSRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "AWS cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"images": {
				Description: "List of retrieved AWS images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "AWS image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "AWS image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"latest": {
				Description: "Retrieve latest AWS images.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"owners": {
				Description: "List of AWS image owners",
				Type:        schema.TypeSet,
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []string{
						"Canonical",
						"Debian",
						"SUSE",
						"Red Hat",
						"Amazon",
						"Microsoft",
					}, nil
				},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func dataSourceTaikunImagesAWSRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*taikungoclient.Client)
	owners, err := dataSourceTaikunImagesAWSGetOwnerID(apiClient, d.Get("owners").(*schema.Set).List())

	body := models.AwsImagesPostListCommand{
		Latest:  d.Get("latest").(bool),
		CloudID: cloudCredentialID,
		Owners:  owners,
	}

	params := images.NewImagesAwsImagesAsPostParams().WithV(ApiVersion).WithBody(&body)

	var imageList []map[string]interface{}
	for {
		response, err := apiClient.Client.Images.ImagesAwsImagesAsPost(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		imageList = append(imageList, flattenTaikunImagesAWS(response.Payload.Data...)...)
		if len(imageList) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(imageList))
		body.Offset = offset
		params = params.WithBody(&body)
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}

// Converts a list of AWS owner names to a list of AWS owner IDs
func dataSourceTaikunImagesAWSGetOwnerID(apiClient *taikungoclient.Client, ownerNames []interface{}) (ownerIds []string, err error) {

	// Get list of Owners with ID and Name from API
	params := aws.NewAwsAwsOwnersParams().WithV(ApiVersion)
	response, err := apiClient.Client.Aws.AwsAwsOwners(params, apiClient)
	if err != nil {
		return
	}

	// Create owner name to owner ID map
	awsOwnerIdNameMap := make(map[string]string)
	for _, owner := range response.Payload {
		awsOwnerIdNameMap[owner.Name] = owner.ID
	}

	// Make list of owner IDs from the list of owner names
	ownerIds = []string{}
	for _, ownerName := range ownerNames {
		ownerId, ok := awsOwnerIdNameMap[ownerName.(string)]
		if !ok {
			return nil, fmt.Errorf("%s AWS owner not found.", ownerName)
		}
		ownerIds = append(ownerIds, ownerId)
	}

	return
}

func flattenTaikunImagesAWS(rawImages ...*models.AwsExtendedImagesListDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.ID,
			"name": rawImage.Name,
		}
	}
	return images
}
