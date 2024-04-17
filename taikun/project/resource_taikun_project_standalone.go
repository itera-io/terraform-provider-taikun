package project

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		"hypervisor": {
			Description:      "Hypervisor used for this VM (required for Proxmox, required for vSphere when DRS is disabled).",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.IgnoreChangeFromEmpty,
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
					//"device_name": {
					//	Description: "Name of the device (required with AWS).",
					//	Type:        schema.TypeString,
					//	Optional:    true,
					//	Computed:    true,
					//	ValidateFunc: validation.StringMatch(
					//		regexp.MustCompile("^/dev/sd[a-z]$"),
					//		"Must be a valid device name",
					//	),
					//},
					"id": {
						Description: "ID of the disk.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"lun_id": {
						Description:  "LUN ID (required with Azure).",
						Type:         schema.TypeInt,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IntBetween(0, 999),
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
		"spot_vm": {
			Description: "Enable if this to create standalone VM on spot instances",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"spot_vm_max_price": {
			Description: "The maximum price you are willing to pay for the spot instance (USD) - Any changes made to this attribute after project creation are ignored by terraform provider. If not specified, the current on-demand price is used.",
			Type:        schema.TypeFloat,
			Optional:    true,
			// Ignore all changes to max price (API returns/sets on-demand spotPrice if we send null spotPrice)
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return old != "" // Set the value only first time
			},
		},
		"standalone_profile_id": {
			Description:      "Standalone profile ID bound to the VM (updating this field will recreate the VM).",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: utils.StringIsInt,
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
		"zone": {
			Description:      "Availability zone for this VM (only for AWS, Azure and GCP). If not specified, the first valid zone is used.",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.IgnoreChangeFromEmpty,
		},
	}
}

func resourceTaikunProjectSetVMs(d *schema.ResourceData, apiClient *tk.Client, projectID int32) error {

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

		if cloudType != utils.CloudTypeOpenStack && hasChanges(old, new, "public_ip") {
			return true
		}

		// ForceNew fields within the VM subresource
		return hasChanges(old, new,
			"cloud_init",
			"flavor",
			"image_id",
			"name",
			"standalone_profile_id",
			"tag",
			"username",
			"volume_size",
			"volume_type",
			"spot_vm",
			"hypervisor",
			"zone",
		)
	}
}

func shouldRecreateDisk(old map[string]interface{}, new map[string]interface{}) bool {
	//return hasChanges(old, new, "device_name", "name", "volume_type")
	return hasChanges(old, new, "name", "volume_type")
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

func resourceTaikunProjectUpdateVMs(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client, projectID int32) error {

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

	cloudCredentialID, _ := utils.Atoi32(d.Get("cloud_credential_id").(string))
	cloudType, err := ResourceTaikunProjectGetCloudType(cloudCredentialID, apiClient)
	if err != nil {
		return err
	}

	toDelete, toAdd, intersection := computeDiff(oldMap, newMap, genVmRecreateFunc(cloudType))

	vmIds := make([]int32, 0)

	for _, vmMap := range toDelete {
		if vmIdStr, vmIdSet := vmMap["id"]; vmIdSet {
			vmId, _ := utils.Atoi32(vmIdStr.(string))
			vmIds = append(vmIds, vmId)
		}
	}

	if len(vmIds) != 0 {
		deleteServerBody := tkcore.DeleteStandAloneVmCommand{}
		deleteServerBody.SetProjectId(projectID)
		deleteServerBody.SetVmIds(vmIds)

		res, err := apiClient.Client.StandaloneAPI.StandaloneDelete(ctx).DeleteStandAloneVmCommand(deleteServerBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
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
		vmId, _ := utils.Atoi32(id)
		if old := findWithId(oldMap, id); old != nil {

			if hasChanges(old, new, "public_ip") {
				repairNeeded = true
				mode := "enable"
				if !new["public_ip"].(bool) {
					mode = "disable"
				}
				body := tkcore.StandAloneVmIpManagementCommand{}
				body.SetId(vmId)
				body.SetMode(mode)

				res, err := apiClient.Client.StandaloneAPI.StandaloneIpManagement(ctx).StandAloneVmIpManagementCommand(body).Execute()
				if err != nil {
					return tk.CreateError(res, err)
				}
			}
			if hasChanges(old, new, "flavor") {
				repairNeeded = true
				body := tkcore.UpdateStandAloneVmFlavorCommand{}
				body.SetId(vmId)
				body.SetFlavor(new["flavor"].(string))

				res, err := apiClient.Client.StandaloneAPI.StandaloneUpdateFlavor(ctx).UpdateStandAloneVmFlavorCommand(body).Execute()
				if err != nil {
					return tk.CreateError(res, err)
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
		body := tkcore.RepairStandAloneVmCommand{}
		body.SetProjectId(projectID)
		res, err := apiClient.Client.StandaloneAPI.StandaloneRepair(ctx).RepairStandAloneVmCommand(body).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return err
		}
	}

	return nil
}

func resourceTaikunProjectUpdateVMDisks(ctx context.Context, oldDisks interface{}, newDisks interface{}, apiClient *tk.Client, vmID int32, projectID int32) error {
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
			diskId, _ := utils.Atoi32(diskIdStr.(string))
			diskIds = append(diskIds, diskId)
		}
	}

	if len(diskIds) != 0 {
		deleteDiskBody := tkcore.DeleteStandAloneVmDiskCommand{}
		deleteDiskBody.SetVmDiskIds(diskIds)
		deleteDiskBody.SetStandaloneVmId(vmID)

		res, err := apiClient.Client.StandaloneVMDisksAPI.StandalonevmdisksDelete(ctx).DeleteStandAloneVmDiskCommand(deleteDiskBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
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
		diskId, _ := utils.Atoi32(id)
		if old := findWithId(oldMap, id); old != nil {
			if hasChanges(old, new, "size") {
				body := tkcore.UpdateStandaloneVmDiskSizeCommand{}
				body.SetId(diskId)
				body.SetSize(int64(new["size"].(int)))

				res, err := apiClient.Client.StandaloneVMDisksAPI.StandalonevmdisksUpdateSize(ctx).UpdateStandaloneVmDiskSizeCommand(body).Execute()
				if err != nil {
					return tk.CreateError(res, err)
				}
			}
		}
		// Shouldn't happen
	}

	return nil
}

func resourceTaikunProjectAddVM(vmMap map[string]interface{}, apiClient *tk.Client, projectID int32) (string, map[string]interface{}, error) {

	standaloneProfileId, _ := utils.Atoi32(vmMap["standalone_profile_id"].(string))
	unreadableProperties := map[string]interface{}{}

	vmCreateBody := tkcore.CreateStandAloneVmCommand{}
	vmCreateBody.SetCount(1)
	vmCreateBody.SetFlavorName(vmMap["flavor"].(string))
	vmCreateBody.SetImage(vmMap["image_id"].(string))
	vmCreateBody.SetName(vmMap["name"].(string))
	vmCreateBody.SetProjectId(projectID)
	vmCreateBody.SetPublicIpEnabled(vmMap["public_ip"].(bool))
	vmCreateBody.SetStandAloneMetaDatas(make([]tkcore.StandAloneMetaDataDto, 0))
	vmCreateBody.SetStandAloneProfileId(standaloneProfileId)
	vmCreateBody.SetStandAloneVmDisks(make([]tkcore.StandAloneVmDiskDto, 0))
	vmCreateBody.SetVolumeSize(int64(vmMap["volume_size"].(int)))

	if vmMap["cloud_init"] != nil {
		vmCreateBody.SetCloudInit(vmMap["cloud_init"].(string))
	}

	if vmMap["hypervisor"] != nil {
		vmCreateBody.SetHypervisor(vmMap["hypervisor"].(string))
	}

	if vmMap["username"] != nil {
		vmCreateBody.SetUsername(vmMap["username"].(string))
	}

	if vmMap["volume_type"] != nil {
		vmCreateBody.SetVolumeType(vmMap["volume_type"].(string))
	}

	if vmMap["zone"] != nil {
		vmCreateBody.SetAvailabilityZone(vmMap["zone"].(string))
	}

	if vmMap["tag"] != nil {
		rawTags := vmMap["tag"].(*schema.Set).List()
		tagsList := make([]tkcore.StandAloneMetaDataDto, len(rawTags))
		for i, e := range rawTags {
			rawTag := e.(map[string]interface{})
			tagsList[i] = tkcore.StandAloneMetaDataDto{}
			tagsList[i].SetKey(rawTag["key"].(string))
			tagsList[i].SetValue(rawTag["value"].(string))
		}
		vmCreateBody.SetStandAloneMetaDatas(tagsList)
	}

	if vmMap["disk"] != nil {
		rawDisks := vmMap["disk"].([]interface{})
		disksList := make([]tkcore.StandAloneVmDiskDto, len(rawDisks))
		for i, e := range rawDisks {
			rawDisk := e.(map[string]interface{})
			disksList[i] = tkcore.StandAloneVmDiskDto{}
			//deviceName := rawDisk["device_name"].(string)
			//if deviceName != "" {
			//	disksList[i].SetDeviceName(deviceName)
			//} else {
			//	disksList[i].SetDeviceNameNil()
			//}
			disksList[i].SetName(rawDisk["name"].(string))
			disksList[i].SetSize(int64(rawDisk["size"].(int)))
			disksList[i].SetVolumeType(rawDisk["volume_type"].(string))
		}
		vmCreateBody.SetStandAloneVmDisks(disksList)
	}

	// Standalone VM spots
	if (vmMap["spot_vm_max_price"].(float64) != 0) && (!vmMap["spot_vm"].(bool)) {
		return "", nil, fmt.Errorf("Spot VM max price is set, but the VM does not have spot enabled.")
	}
	if vmMap["spot_vm"] != nil {
		spotForThisVm := vmMap["spot_vm"].(bool)
		vmCreateBody.SetSpotInstance(spotForThisVm)
		vmCreateBody.SetSpotPrice(vmMap["spot_vm_max_price"].(float64))
		if vmMap["spot_vm_max_price"].(float64) == 0 {
			vmCreateBody.UnsetSpotPrice() // Send null if the user did not specify anything
		}
	}

	vmCreateResponse, res, err := apiClient.Client.StandaloneAPI.StandaloneCreate(context.TODO()).CreateStandAloneVmCommand(vmCreateBody).Execute()
	if err != nil {
		return "", nil, tk.CreateError(res, err)
	}

	return vmCreateResponse.GetId(), unreadableProperties, nil
}

func resourceTaikunProjectAddDisk(diskMap map[string]interface{}, apiClient *tk.Client, vmId int32) error {

	diskCreateBody := tkcore.CreateStandAloneDiskCommand{}
	// diskCreateBody.SetDeviceName(diskMap["device_name"].(string)) // Removed from API at 27.11.2023 by Arzu. Method forever in our hearts.
	diskCreateBody.SetName(diskMap["name"].(string))
	diskCreateBody.SetSize(int64(diskMap["size"].(int)))
	diskCreateBody.SetVolumeType(diskMap["volume_type"].(string))
	diskCreateBody.SetStandaloneVmId(vmId)

	_, res, err := apiClient.Client.StandaloneVMDisksAPI.StandalonevmdisksCreate(context.TODO()).CreateStandAloneDiskCommand(diskCreateBody).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}

	return nil
}

func resourceTaikunProjectStandaloneCommit(apiClient *tk.Client, projectID int32) error {
	body := tkcore.CommitStandAloneVmCommand{}
	body.SetProjectId(projectID)
	res, err := apiClient.Client.StandaloneAPI.StandaloneCommit(context.TODO()).CommitStandAloneVmCommand(body).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	return nil
}

func resourceTaikunProjectEditImages(d *schema.ResourceData, apiClient *tk.Client, id int32) error {
	// Get cloud type (because not all cloud types use image.id, some use image.name instead)
	ccID, err := utils.Atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return err
	}
	cloudType, err := ResourceTaikunProjectGetCloudType(ccID, apiClient)
	if err != nil {
		return err
	}

	// Bind / unbind images
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
			if cloudType == string(tkcore.CLOUDTYPE_GOOGLE) {
				// GCP uses names to identify images
				if imagesToUnbind.Contains(boundImageDTO.GetName()) {
					imageBindingsToUndo = append(imageBindingsToUndo, boundImageDTO.GetId())
				}
			} else {
				// All other providers use ids to identify images
				if imagesToUnbind.Contains(boundImageDTO.GetImageId()) {
					imageBindingsToUndo = append(imageBindingsToUndo, boundImageDTO.GetId())
				}
			}
		}
		unbindBody := tkcore.DeleteImageFromProjectCommand{}
		unbindBody.SetIds(imageBindingsToUndo)
		res, err := apiClient.Client.ImagesAPI.ImagesUnbindImagesFromProject(context.TODO()).DeleteImageFromProjectCommand(unbindBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	if len(imagesToBind) != 0 {
		imagesToBindNames := make([]string, len(imagesToBind))
		for i, imageToBind := range imagesToBind {
			imagesToBindNames[i] = imageToBind.(string)
		}
		bindBody := tkcore.BindImageToProjectCommand{}
		bindBody.SetProjectId(id)
		bindBody.SetImages(imagesToBindNames)
		res, err := apiClient.Client.ImagesAPI.ImagesBindImagesToProject(context.TODO()).BindImageToProjectCommand(bindBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	return nil
}

func resourceTaikunProjectPurgeVMs(vmsToPurge []interface{}, apiClient *tk.Client, projectID int32) error {
	vmIds := make([]int32, 0)

	for _, vm := range vmsToPurge {
		vmMap := vm.(map[string]interface{})
		if vmIdStr, vmIdSet := vmMap["id"]; vmIdSet {
			vmId, _ := utils.Atoi32(vmIdStr.(string))
			if vmId != 0 {
				vmIds = append(vmIds, vmId)
			}
		}
	}

	if len(vmIds) != 0 {
		deleteServerBody := tkcore.DeleteStandAloneVmCommand{}
		deleteServerBody.SetProjectId(projectID)
		deleteServerBody.SetVmIds(vmIds)

		res, err := apiClient.Client.StandaloneAPI.StandaloneDelete(context.TODO()).DeleteStandAloneVmCommand(deleteServerBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	return nil
}

func resourceTaikunProjectToggleVmsSpot(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	bodyToggle := tkcore.SpotVmOperationCommand{}
	bodyToggle.SetId(projectID)

	if d.Get("spot_vms").(bool) {
		bodyToggle.SetMode("enable")
	} else if !d.Get("spot_full").(bool) {
		bodyToggle.SetMode("disable")
	}

	res, err := apiClient.Client.ProjectsAPI.ProjectsToggleSpotVms(ctx).SpotVmOperationCommand(bodyToggle).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	return nil
}
