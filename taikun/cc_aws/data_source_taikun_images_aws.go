package cc_aws

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceTaikunImagesAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given AWS cloud credential.",
		ReadContext: dataSourceTaikunImagesAWSRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "AWS cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.StringIsInt,
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
	cloudCredentialID, err := utils.Atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*tk.Client)
	owners, err := dataSourceTaikunImagesAWSGetOwnerID(apiClient, d.Get("owners").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.AwsImagesPostListCommand{}
	body.SetLatest(d.Get("latest").(bool))
	body.SetCloudId(cloudCredentialID)
	body.SetOwners(owners)

	var offset int32 = 0
	var imageList []map[string]interface{}
	for {
		body.SetOffset(offset)
		response, res, err := apiClient.Client.ImagesAPI.ImagesAwsImagesList(context.TODO()).AwsImagesPostListCommand(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		imageList = append(imageList, flattenTaikunImagesAWS(response.GetData()...)...)
		if len(imageList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(imageList))
	}

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.I32toa(cloudCredentialID))
	return nil
}

// Converts a list of AWS owner names to a list of AWS owner IDs
func dataSourceTaikunImagesAWSGetOwnerID(apiClient *tk.Client, ownerNames []interface{}) (ownerIds []string, err error) {

	// Get list of Owners with ID and Name from API
	response, res, err := apiClient.Client.AWSCloudCredentialAPI.AwsOwners(context.TODO()).Execute()
	if err != nil {
		err = tk.CreateError(res, err)
		return
	}

	// Create owner name to owner ID map
	awsOwnerIdNameMap := make(map[string]string)
	for _, owner := range response {
		awsOwnerIdNameMap[owner.GetName()] = owner.GetId()
	}

	// Make list of owner IDs from the list of owner names
	ownerIds = []string{}
	for _, ownerName := range ownerNames {
		ownerId, ok := awsOwnerIdNameMap[ownerName.(string)]
		if !ok {
			return nil, fmt.Errorf("%s AWS owner not found", ownerName)
		}
		ownerIds = append(ownerIds, ownerId)
	}

	return
}

func flattenTaikunImagesAWS(rawImages ...tkcore.CommonStringBasedDropdownDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.GetId(),
			"name": rawImage.GetName(),
		}
	}
	return images
}
