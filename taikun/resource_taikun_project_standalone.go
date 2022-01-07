package taikun

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/models"
)

func taikunVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_ip": {
			Description: "The public IP.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"cloud_init": {
			Description: "Cloud init",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			ForceNew:    true,
		},
		"created_by": {
			Description: "The creator of the VM.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"disk": {
			Description: "Disks associated with the VM.",
			Type:        schema.TypeSet,
			Optional:    true,
			Set:         hashAttributes("name", "device_name", "lun_id", "volume_type"),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Description: "Name of the device (Only valid with AWS).",
						Type:        schema.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^/dev/sd[a-z]$"),
							"Must be a valid device name",
						),
					},
					"disk_size": {
						Description:  "The disk size in GBs.",
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"lun_id": {
						Description:  "LUN ID (only valid with Azure).",
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 999),
					},
					"name": {
						Description: "Name of the Disk.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(3, 30),
							validation.StringMatch(
								regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
								"expected only alpha numeric characters or non alpha numeric (_-.)",
							),
						),
					},
					"volume_type": {
						Description: "Type of the volume (only valid with OpenStack).",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
		},
		"disk_size": {
			Description:  "The VM's disk size in GBs.",
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"flavor": {
			Description:  "The VM's flavor.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"id": {
			Description: "ID of the VM.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"image": {
			Description:  "The VM's image.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"ip": {
			Description: "IP of the VM.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "The time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the VM.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "Name of the VM.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 52),
		},
		"public_ip": {
			Description: "Whether a public IP will be available",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"standalone_profile_id": {
			Description:      "Standalone profile ID.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"status": {
			Description: "VM status.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"tag": {
			Description: "Tags linked to the VM.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Description: "Key of the tag.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(0, 63),
							validation.StringMatch(
								regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
								"expected only alpha numeric characters or non alpha numeric (_-.)",
							),
						),
					},
					"value": {
						Description: "Value of the tag.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(0, 63),
							validation.StringMatch(
								regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
								"expected only alpha numeric characters or non alpha numeric (_-.)",
							),
						),
					},
				},
			},
		},
		"volume_type": {
			Description: "Volume type (only valid with OpenStack).",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
		},
	}
}

func resourceTaikunProjectEditImages(d *schema.ResourceData, apiClient *apiClient, id int32) error {
	oldImageData, newImageData := d.GetChange("images")
	oldImages := oldImageData.(*schema.Set)
	newImages := newImageData.(*schema.Set)
	imagesToUnbind := oldImages.Difference(newImages)
	imagesToBind := newImages.Difference(oldImages).List()
	boundImageDTOs, err := resourceTaikunProjectGetBoundImageDTOs(id, apiClient)
	if err != nil {
		return err
	}
	if imagesToUnbind.Len() != 0 {
		var imageBindingsToUndo []int32
		for _, boundImageDTO := range boundImageDTOs {
			if imagesToUnbind.Contains(boundImageDTO.ImageID) {
				imageBindingsToUndo = append(imageBindingsToUndo, boundImageDTO.ID)
			}
		}
		unbindBody := models.DeleteImageFromProjectCommand{Ids: imageBindingsToUndo}
		unbindParams := images.NewImagesUnbindImagesFromProjectParams().WithV(ApiVersion).WithBody(&unbindBody)
		if _, err := apiClient.client.Images.ImagesUnbindImagesFromProject(unbindParams, apiClient); err != nil {
			return err
		}
	}
	if len(imagesToBind) != 0 {
		imagesToBindNames := make([]string, len(imagesToBind))
		for i, imageToBind := range imagesToBind {
			imagesToBindNames[i] = imageToBind.(string)
		}
		bindBody := models.BindImageToProjectCommand{ProjectID: id, Images: imagesToBindNames}
		bindParams := images.NewImagesBindImagesToProjectParams().WithV(ApiVersion).WithBody(&bindBody)
		if _, err := apiClient.client.Images.ImagesBindImagesToProject(bindParams, apiClient); err != nil {
			return err
		}
	}
	return nil
}
