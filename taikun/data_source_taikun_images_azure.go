package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceTaikunImagesAzureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var offset int32 = 0

	latest := d.Get("latest").(bool)

	params := apiClient.Client.ImagesApi.ImagesAzureImages(ctx, cloudCredentialID, d.Get("publisher").(string),
		d.Get("offer").(string), d.Get("sku").(string)).Latest(latest)

	var imageList []map[string]interface{}
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

	if err := d.Set("images", imageList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}
