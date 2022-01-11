package taikun

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/client/stand_alone"
	"github.com/itera-io/taikungoclient/models"
)

func taikunVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_ip": {
			Description: "Access IP of the VM.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"cloud_init": {
			Description: "Cloud init.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
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
						Description: "Name of the device (required with AWS).",
						Type:        schema.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^/dev/sd[a-z]$"),
							"Must be a valid device name",
						),
					},
					"lun_id": {
						Description:  "LUN ID (required with Azure).",
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
					"size": {
						Description:  "The disk size in GBs.",
						Type:         schema.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"volume_type": {
						Description: "Type of the volume (only valid with OpenStack).",
						Type:        schema.TypeString,
						Optional:    true,
					},
				},
			},
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
		"image_id": {
			Description:  "The VM's image id.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"image_name": {
			Description: "The VM's image name.",
			Type:        schema.TypeString,
			Computed:    true,
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
			ValidateFunc: validation.StringLenBetween(1, 52),
		},
		"public_ip": {
			Description: "Whether a public IP will be available.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"standalone_profile_id": {
			Description:      "Standalone profile ID bound to the VM.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"status": {
			Description: "VM status.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"tag": {
			Description: "Tags linked to the VM.",
			Type:        schema.TypeSet,
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
		"volume_size": {
			Description:  "The VM's volume size in GBs.",
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"volume_type": {
			Description: "Volume type (only valid with OpenStack).",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
		},
	}
}

func resourceTaikunProjectSetVMs(d *schema.ResourceData, apiClient *apiClient, projectID int32) error {

	vms := d.Get("vm")

	vmsList := vms.(*schema.Set).List()
	for _, vm := range vmsList {
		vmMap := vm.(map[string]interface{})

		vmId, err := resourceTaikunProjectAddVM(vmMap, apiClient, projectID)
		if err != nil {
			return err
		}
		vmMap["id"] = vmId
	}
	err := d.Set("vm", vmsList)
	if err != nil {
		return err
	}

	return nil
}

func resourceTaikunProjectUpdateVMs(d *schema.ResourceData, apiClient *apiClient, projectID int32) error {

	oldVms, newVms := d.GetChange("vm")
	oldVmsSet := oldVms.(*schema.Set)
	newVmsSet := newVms.(*schema.Set)

	toDelete := oldVmsSet.Difference(newVmsSet).List()
	toAdd := newVmsSet.Difference(oldVmsSet).List()

	vmIds := make([]int32, 0)

	for _, vm := range toDelete {
		vmMap := vm.(map[string]interface{})
		if vmIdStr, vmIdSet := vmMap["id"]; vmIdSet {
			vmId, _ := atoi32(vmIdStr.(string))
			vmIds = append(vmIds, vmId)
		}
	}

	if len(vmIds) != 0 {
		deleteServerBody := &models.DeleteStandAloneVMCommand{
			ProjectID: projectID,
			VMIds:     vmIds,
		}
		deleteVMParams := stand_alone.NewStandAloneDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
		_, err := apiClient.client.StandAlone.StandAloneDelete(deleteVMParams, apiClient)
		if err != nil {
			return err
		}
	}

	vmsList := oldVmsSet.Intersection(newVmsSet).List()

	for _, vm := range toAdd {
		vmMap := vm.(map[string]interface{})

		vmId, err := resourceTaikunProjectAddVM(vmMap, apiClient, projectID)
		if err != nil {
			return err
		}
		vmMap["id"] = vmId

		vmsList = append(vmsList, vmMap)
	}
	err := d.Set("vm", vmsList)
	if err != nil {
		return err
	}

	return nil
}

func resourceTaikunProjectAddVM(vmMap map[string]interface{}, apiClient *apiClient, projectID int32) (string, error) {

	standaloneProfileId, _ := atoi32(vmMap["standalone_profile_id"].(string))

	vmCreateBody := &models.CreateStandAloneVMCommand{
		CloudInit:           vmMap["cloud_init"].(string),
		Count:               1,
		FlavorName:          vmMap["flavor"].(string),
		Image:               vmMap["image_id"].(string),
		Name:                vmMap["name"].(string),
		ProjectID:           projectID,
		PublicIPEnabled:     vmMap["public_ip"].(bool),
		StandAloneProfileID: standaloneProfileId,
		VolumeSize:          int64(vmMap["volume_size"].(int)),
		VolumeType:          vmMap["volume_type"].(string),
	}

	if vmMap["tag"] != nil {
		rawTags := vmMap["tag"].(*schema.Set).List()
		tagsList := make([]*models.StandAloneMetaDataDto, len(rawTags))
		for i, e := range rawTags {
			rawTag := e.(map[string]interface{})
			tagsList[i] = &models.StandAloneMetaDataDto{
				Key:   rawTag["key"].(string),
				Value: rawTag["value"].(string),
			}
		}
		vmCreateBody.StandAloneMetaDatas = tagsList
	}

	if vmMap["disk"] != nil {
		rawDisks := vmMap["disk"].(*schema.Set).List()
		disksList := make([]*models.StandAloneVMDiskDto, len(rawDisks))
		for i, e := range rawDisks {
			rawDisk := e.(map[string]interface{})
			disksList[i] = &models.StandAloneVMDiskDto{
				DeviceName: rawDisk["device_name"].(string),
				LunID:      int32(rawDisk["lun_id"].(int)),
				Name:       rawDisk["name"].(string),
				Size:       int64(rawDisk["size"].(int)),
				VolumeType: rawDisk["volume_type"].(string),
			}
		}
		vmCreateBody.StandAloneVMDisks = disksList
	}

	vmCreateParams := stand_alone.NewStandAloneCreateParams().WithV(ApiVersion).WithBody(vmCreateBody)
	vmCreateResponse, err := apiClient.client.StandAlone.StandAloneCreate(vmCreateParams, apiClient)
	if err != nil {
		return "", err
	}

	return vmCreateResponse.Payload.ID, nil
}

func resourceTaikunProjectStandaloneCommit(apiClient *apiClient, projectID int32) error {
	body := &models.CommitStandAloneVMCommand{ProjectID: projectID}
	params := stand_alone.NewStandAloneCommitParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.client.StandAlone.StandAloneCommit(params, apiClient)
	if err != nil {
		return err
	}
	return nil
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

func resourceTaikunProjectPurgeVMs(vmsToPurge []interface{}, apiClient *apiClient, projectID int32) error {
	vmIds := make([]int32, 0)

	for _, vm := range vmsToPurge {
		vmMap := vm.(map[string]interface{})
		if vmIdStr, vmIdSet := vmMap["id"]; vmIdSet {
			vmId, _ := atoi32(vmIdStr.(string))
			vmIds = append(vmIds, vmId)
		}
	}

	if len(vmIds) != 0 {
		deleteServerBody := &models.DeleteStandAloneVMCommand{
			ProjectID: projectID,
			VMIds:     vmIds,
		}
		deleteVMParams := stand_alone.NewStandAloneDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
		_, err := apiClient.client.StandAlone.StandAloneDelete(deleteVMParams, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}
