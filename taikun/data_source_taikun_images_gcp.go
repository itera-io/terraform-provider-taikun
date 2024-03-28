package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceTaikunImagesGCP() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given GCP cloud credential.",
		ReadContext: dataSourceTaikunImagesGCPRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "GCP cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"images": {
				Description: "List of retrieved GCP images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "GCP image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "GCP image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"latest": {
				Description: "Retrieve latest GCP images.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"type": {
				Description: "GCP image type.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				ValidateFunc: validation.StringInSlice([]string{
					"all",
					"ubuntu",
					"debian",
					"windows",
				}, false),
			},
		},
	}
}

func dataSourceTaikunImagesGCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	latest := d.Get("latest").(bool)
	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var offset int32 = 0
	apiClient := meta.(*tk.Client)
	params := apiClient.Client.ImagesAPI.ImagesGoogleImages(ctx, cloudCredentialID, d.Get("type").(string))

	var imageList []map[string]interface{}
	for {
		response, res, err := params.Offset(offset).Latest(latest).Execute()
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
