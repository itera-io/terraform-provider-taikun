package cc_zadara

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunImagesZadara() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve images for a given Zadara cloud credential.",
		ReadContext: dataSourceTaikunImagesZadaraRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "Zadara cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
			"images": {
				Description: "List of retrieved Zadara images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Zadara image ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Zadara image name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"latest": {
				Description: "Retrieve latest Zadara images.",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func dataSourceTaikunImagesZadaraRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudCredentialID, err := utils.Atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*tk.Client)

	params := apiClient.Client.ImagesAPI.ImagesZadaraImagesList(context.TODO(), cloudCredentialID).Latest(d.Get("latest").(bool))

	var offset int32 = 0
	var imageList []map[string]interface{}
	for {
		params.Offset(offset)
		response, res, err := params.Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		imageList = append(imageList, flattenTaikunImagesZadara(response.GetData()...)...)
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

func flattenTaikunImagesZadara(rawImages ...tkcore.CommonStringBasedDropdownDto) []map[string]interface{} {

	images := make([]map[string]interface{}, len(rawImages))
	for i, rawImage := range rawImages {
		images[i] = map[string]interface{}{
			"id":   rawImage.GetId(),
			"name": rawImage.GetName(),
		}
	}
	return images
}
