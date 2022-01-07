package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/models"
)

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
