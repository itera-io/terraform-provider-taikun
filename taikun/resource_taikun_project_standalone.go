package taikun

import (
	"context"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/client/stand_alone"
	"github.com/itera-io/taikungoclient/client/stand_alone_vm_disks"
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
			Description: "Cloud init (updating this field will recreate the VM).",
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
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_name": {
						Description: "Name of the device (required with AWS).",
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^/dev/sd[a-z]$"),
							"Must be a valid device name",
						),
					},
					"id": {
						Description: "ID of the disk.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "Name of the disk.",
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
						Computed:    true,
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
			Description:  "The VM's image ID (updating this field will recreate the VM).",
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
			Description:  "Name of the VM (updating this field will recreate the VM).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 52),
		},
		"public_ip": {
			Description: "Whether a public IP will be available (updating this field will recreate the VM if the project isn't hosted on OpenStack).",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"standalone_profile_id": {
			Description:      "Standalone profile ID bound to the VM (updating this field will recreate the VM).",
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
			Description: "Tags linked to the VM (updating this field will recreate the VM).",
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
		"username": {
			Description:  "The VM's username (required for Azure).",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 20),
		},
		"volume_size": {
			Description:  "The VM's volume size in GBs (updating this field will recreate the VM).",
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"volume_type": {
			Description: "Volume type (updating this field will recreate the VM).",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
	}
}

func resourceTaikunProjectSetVMs(d *schema.ResourceData, apiClient *taikungoclient.Client, projectID int32) error {

	vms := d.Get("vm")

	vmsList := vms.([]interface{})
	for _, vm := range vmsList {
		vmMap := vm.(map[string]interface{})

		vmId, unreadableProperties, err := resourceTaikunProjectAddVM(vmMap, apiClient, projectID)
		if err != nil {
			return err
		}
		vmMap["id"] = vmId

		for key, value := range unreadableProperties {
			vmMap[key] = value
		}
	}

	err := d.Set("vm", vmsList)
	if err != nil {
		return err
	}

	return nil
}

func findWithId(searchMap []map[string]interface{}, id string) map[string]interface{} {
	for _, f := range searchMap {
		if f["id"] == id && id != "" {
			return f
		}
	}
	return nil
}

func hasChanges(old map[string]interface{}, new map[string]interface{}, labels ...string) bool {
	for _, label := range labels {
		// Special compare function for sets
		if set, isSet := old[label].(*schema.Set); isSet {
			if !set.Equal(new[label]) {
				return true
			}
		} else if !reflect.DeepEqual(old[label], new[label]) {
			return true
		}
	}
	return false
}

func genVmRecreateFunc(cloudType string) func(old, new map[string]interface{}) bool {
	return func(old, new map[string]interface{}) bool {

		if cloudType != cloudTypeOpenStack && hasChanges(old, new, "public_ip") {
			return true
		}

		// ForceNew fields within the VM subresource
		return hasChanges(old, new,
			"cloud_init",
			"image_id",
			"name",
			"standalone_profile_id",
			"tag",
			"username",
			"volume_size",
			"volume_type",
		)
	}
}

func shouldRecreateDisk(old map[string]interface{}, new map[string]interface{}) bool {
	return hasChanges(old, new, "device_name", "name", "volume_type")
}

func computeDiff(oldMap []map[string]interface{}, newMap []map[string]interface{}, recreateFunc func(old map[string]interface{}, new map[string]interface{}) bool) ([]map[string]interface{}, []map[string]interface{}, []map[string]interface{}) {
	toDelete, toAdd, intersection := make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)

	// Vms which don't have id will be added
	for _, e := range newMap {
		if e["id"] == nil || e["id"].(string) == "" {
			toAdd = append(toAdd, e)
			continue
		}
		intersection = append(intersection, e)
	}

	// Vms which are no longer in the list will be deleted
	for _, e := range oldMap {
		id := e["id"].(string)
		if findWithId(newMap, id) != nil {
			continue
		}

		toDelete = append(toDelete, e)
	}

	// Vms which have ForceNew changes will be deleted and added
	for _, new := range intersection {
		id := new["id"].(string)
		if old := findWithId(oldMap, id); old != nil {
			if recreateFunc(old, new) {
				toDelete = append(toDelete, old)
				toAdd = append(toAdd, new)
			}
		}
		// Shouldn't happen
	}

	return toDelete, toAdd, intersection
}

func resourceTaikunProjectUpdateVMs(ctx context.Context, d *schema.ResourceData, apiClient *taikungoclient.Client, projectID int32) error {

	oldVms, newVms := d.GetChange("vm")
	oldVmsList := oldVms.([]interface{})
	newVmsList := newVms.([]interface{})
	oldMap, newMap := make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)
	for _, e := range oldVmsList {
		oldMap = append(oldMap, e.(map[string]interface{}))
	}
	for _, e := range newVmsList {
		newMap = append(newMap, e.(map[string]interface{}))
	}

	cloudCredentialID, _ := atoi32(d.Get("cloud_credential_id").(string))
	cloudType, err := resourceTaikunProjectGetCloudType(cloudCredentialID, apiClient)
	if err != nil {
		return err
	}

	toDelete, toAdd, intersection := computeDiff(oldMap, newMap, genVmRecreateFunc(cloudType))

	vmIds := make([]int32, 0)

	for _, vmMap := range toDelete {
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
		_, err := apiClient.Client.StandAlone.StandAloneDelete(deleteVMParams, apiClient)
		if err != nil {
			return err
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging", "Deleting", "PendingDelete"}, apiClient, projectID); err != nil {
			return err
		}
	}

	if len(toAdd) != 0 {
		vmsList := intersection

		for _, vmMap := range toAdd {

			vmId, unreadableProperties, err := resourceTaikunProjectAddVM(vmMap, apiClient, projectID)
			if err != nil {
				return err
			}
			vmMap["id"] = vmId

			for key, value := range unreadableProperties {
				vmMap[key] = value
			}

			vmsList = append(vmsList, vmMap)
		}
		err := d.Set("vm", vmsList)
		if err != nil {
			return err
		}

		if err := resourceTaikunProjectStandaloneCommit(apiClient, projectID); err != nil {
			return err
		}
		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return err
		}
	}

	repairNeeded := false
	for _, new := range intersection {
		id := new["id"].(string)
		vmId, _ := atoi32(id)
		if old := findWithId(oldMap, id); old != nil {

			if hasChanges(old, new, "public_ip") {
				repairNeeded = true
				mode := "enable"
				if !new["public_ip"].(bool) {
					mode = "disable"
				}
				body := &models.StandAloneVMIPManagementCommand{
					ID:   int32Address(vmId),
					Mode: stringAddress(mode),
				}
				params := stand_alone.NewStandAloneIPManagementParams().WithV(ApiVersion).WithBody(body)
				_, err := apiClient.Client.StandAlone.StandAloneIPManagement(params, apiClient)
				if err != nil {
					return err
				}
			}
			if hasChanges(old, new, "flavor") {
				repairNeeded = true
				body := &models.UpdateStandAloneVMFlavorCommand{
					ID:     int32Address(vmId),
					Flavor: stringAddress(new["flavor"]),
				}
				params := stand_alone.NewStandAloneUpdateFlavorParams().WithV(ApiVersion).WithBody(body)
				_, err := apiClient.Client.StandAlone.StandAloneUpdateFlavor(params, apiClient)
				if err != nil {
					return err
				}
			}

			if hasChanges(old, new, "disk") {
				repairNeeded = true
				err := resourceTaikunProjectUpdateVMDisks(ctx, old["disk"], new["disk"], apiClient, vmId, projectID)
				if err != nil {
					return err
				}
			}

		}
		// Shouldn't happen
	}
	if repairNeeded {
		body := &models.RepairStandAloneVMCommand{ProjectID: projectID}
		params := stand_alone.NewStandAloneRepairParams().WithV(ApiVersion).WithBody(body)
		_, err := apiClient.Client.StandAlone.StandAloneRepair(params, apiClient)
		if err != nil {
			return err
		}
		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return err
		}
	}

	return nil
}

func resourceTaikunProjectUpdateVMDisks(ctx context.Context, oldDisks interface{}, newDisks interface{}, apiClient *taikungoclient.Client, vmID int32, projectID int32) error {
	oldDisksList := oldDisks.([]interface{})
	newDisksList := newDisks.([]interface{})
	oldMap, newMap := make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)
	for _, e := range oldDisksList {
		oldMap = append(oldMap, e.(map[string]interface{}))
	}
	for _, e := range newDisksList {
		newMap = append(newMap, e.(map[string]interface{}))
	}

	toDelete, toAdd, intersection := computeDiff(oldMap, newMap, shouldRecreateDisk)

	diskIds := make([]int32, 0)

	for _, diskMap := range toDelete {
		if diskIdStr, diskIdSet := diskMap["id"]; diskIdSet {
			diskId, _ := atoi32(diskIdStr.(string))
			diskIds = append(diskIds, diskId)
		}
	}

	if len(diskIds) != 0 {
		deleteDiskBody := &models.DeleteStandAloneVMDiskCommand{
			StandaloneVMID: vmID,
			VMDiskIds:      diskIds,
		}
		deleteDiskParams := stand_alone_vm_disks.NewStandAloneVMDisksDeleteParams().WithV(ApiVersion).WithBody(deleteDiskBody)
		_, err := apiClient.Client.StandAloneVMDisks.StandAloneVMDisksDelete(deleteDiskParams, apiClient)
		if err != nil {
			return err
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return err
		}
	}

	for _, diskMap := range toAdd {
		err := resourceTaikunProjectAddDisk(diskMap, apiClient, vmID)
		if err != nil {
			return err
		}
	}

	for _, new := range intersection {
		id := new["id"].(string)
		diskId, _ := atoi32(id)
		if old := findWithId(oldMap, id); old != nil {
			if hasChanges(old, new, "size") {
				body := &models.UpdateStandaloneVMDiskSizeCommand{
					ID:   diskId,
					Size: int64(new["size"].(int)),
				}
				params := stand_alone_vm_disks.NewStandAloneVMDisksUpdateDiskSizeParams().WithV(ApiVersion).WithBody(body)
				_, err := apiClient.Client.StandAloneVMDisks.StandAloneVMDisksUpdateDiskSize(params, apiClient)
				if err != nil {
					return err
				}
			}
		}
		// Shouldn't happen
	}

	return nil
}

func resourceTaikunProjectAddVM(vmMap map[string]interface{}, apiClient *taikungoclient.Client, projectID int32) (string, map[string]interface{}, error) {

	standaloneProfileId, _ := atoi32(vmMap["standalone_profile_id"].(string))
	unreadableProperties := map[string]interface{}{}

	vmCreateBody := &models.CreateStandAloneVMCommand{
		CloudInit:           vmMap["cloud_init"].(string),
		Count:               1,
		FlavorName:          stringAddress(vmMap["flavor"]),
		Image:               stringAddress(vmMap["image_id"]),
		Name:                stringAddress(vmMap["name"]),
		ProjectID:           int32Address(projectID),
		PublicIPEnabled:     vmMap["public_ip"].(bool),
		StandAloneMetaDatas: make([]*models.StandAloneMetaDataDto, 0),
		StandAloneProfileID: int32Address(standaloneProfileId),
		StandAloneVMDisks:   make([]*models.StandAloneVMDiskDto, 0),
		VolumeSize:          int64(vmMap["volume_size"].(int)),
	}

	if vmMap["username"] != nil {
		vmCreateBody.Username = vmMap["username"].(string)
		unreadableProperties["username"] = vmCreateBody.Username
	}

	if vmMap["volume_type"] != nil {
		vmCreateBody.VolumeType = vmMap["volume_type"].(string)
	}

	if vmMap["tag"] != nil {
		rawTags := vmMap["tag"].(*schema.Set).List()
		tagsList := make([]*models.StandAloneMetaDataDto, len(rawTags))
		for i, e := range rawTags {
			rawTag := e.(map[string]interface{})
			tagsList[i] = &models.StandAloneMetaDataDto{
				Key:   stringAddress(rawTag["key"]),
				Value: stringAddress(rawTag["value"]),
			}
		}
		vmCreateBody.StandAloneMetaDatas = tagsList
	}

	if vmMap["disk"] != nil {
		rawDisks := vmMap["disk"].([]interface{})
		disksList := make([]*models.StandAloneVMDiskDto, len(rawDisks))
		for i, e := range rawDisks {
			rawDisk := e.(map[string]interface{})
			disksList[i] = &models.StandAloneVMDiskDto{
				DeviceName: rawDisk["device_name"].(string),
				Name:       stringAddress(rawDisk["name"]),
				Size:       int64(rawDisk["size"].(int)),
				VolumeType: rawDisk["volume_type"].(string),
			}
		}
		vmCreateBody.StandAloneVMDisks = disksList
	}

	vmCreateParams := stand_alone.NewStandAloneCreateParams().WithV(ApiVersion).WithBody(vmCreateBody)
	vmCreateResponse, err := apiClient.Client.StandAlone.StandAloneCreate(vmCreateParams, apiClient)
	if err != nil {
		return "", nil, err
	}

	return vmCreateResponse.Payload.ID, unreadableProperties, nil
}

func resourceTaikunProjectAddDisk(diskMap map[string]interface{}, apiClient *taikungoclient.Client, vmId int32) error {

	diskCreateBody := &models.CreateStandAloneDiskCommand{
		DeviceName:     diskMap["device_name"].(string),
		Name:           stringAddress(diskMap["name"]),
		Size:           int64(diskMap["size"].(int)),
		VolumeType:     diskMap["volume_type"].(string),
		StandaloneVMID: int32Address(vmId),
	}

	diskCreateParams := stand_alone_vm_disks.NewStandAloneVMDisksCreateParams().WithV(ApiVersion).WithBody(diskCreateBody)
	_, err := apiClient.Client.StandAloneVMDisks.StandAloneVMDisksCreate(diskCreateParams, apiClient)
	if err != nil {
		return err
	}

	return nil
}

func resourceTaikunProjectStandaloneCommit(apiClient *taikungoclient.Client, projectID int32) error {
	body := &models.CommitStandAloneVMCommand{ProjectID: projectID}
	params := stand_alone.NewStandAloneCommitParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.Client.StandAlone.StandAloneCommit(params, apiClient)
	if err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectEditImages(d *schema.ResourceData, apiClient *taikungoclient.Client, id int32) error {
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
		if _, err := apiClient.Client.Images.ImagesUnbindImagesFromProject(unbindParams, apiClient); err != nil {
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
		if _, err := apiClient.Client.Images.ImagesBindImagesToProject(bindParams, apiClient); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectPurgeVMs(vmsToPurge []interface{}, apiClient *taikungoclient.Client, projectID int32) error {
	vmIds := make([]int32, 0)

	for _, vm := range vmsToPurge {
		vmMap := vm.(map[string]interface{})
		if vmIdStr, vmIdSet := vmMap["id"]; vmIdSet {
			vmId, _ := atoi32(vmIdStr.(string))
			if vmId != 0 {
				vmIds = append(vmIds, vmId)
			}
		}
	}

	if len(vmIds) != 0 {
		deleteServerBody := &models.DeleteStandAloneVMCommand{
			ProjectID: projectID,
			VMIds:     vmIds,
		}
		deleteVMParams := stand_alone.NewStandAloneDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
		_, err := apiClient.Client.StandAlone.StandAloneDelete(deleteVMParams, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}
