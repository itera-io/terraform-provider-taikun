package project

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func taikunServerKubeworkerSchema() map[string]*schema.Schema {
	kubeworkerSchema := taikunServerSchemaWithKubernetesNodeLabels()
	kubeworkerSchema["proxmox_extra_disk_size"] = &schema.Schema{
		Description: "Specify the size of the Proxmox extra storage to enable proxmox storage. Proxmox storage type will be chosen automatically base on the Kubernetes profile used.",
		Type:        schema.TypeInt,
		Optional:    true,
	}
	utils.RemoveForceNewsFromSchema(kubeworkerSchema)
	return kubeworkerSchema
}

// Only for Controlplane and Workers
func taikunServerSchemaWithKubernetesNodeLabels() map[string]*schema.Schema {
	serverSchema := taikunServerBasicSchema()
	serverSchema["wasm"] = &schema.Schema{
		Description: "Enable if the server should support WASM.",
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
	}
	serverSchema["kubernetes_node_label"] = &schema.Schema{
		Description: "Attach Kubernetes node labels.",
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Description: "Kubernetes node label key.",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: validation.All(
						validation.StringLenBetween(1, 63),
						validation.StringMatch(
							regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
							"expected only alpha numeric characters or non alpha numeric (_-.)",
						),
					),
				},
				"value": {
					Description: "Kubernetes node label value.",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: validation.All(
						validation.StringLenBetween(1, 63),
						validation.StringMatch(
							regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
							"expected only alpha numeric characters or non alpha numeric (_-.)",
						),
					),
				},
			},
		},
	}
	return serverSchema
}

func taikunServerBasicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the server.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"disk_size": {
			Description:  "The server's disk size in GBs.",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(30),
			Default:      30,
		},
		"flavor": {
			Description:  "The server's flavor.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"id": {
			Description: "ID of the server.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"ip": {
			Description: "IP of the server.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "The time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the server.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name of the server.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or non alpha numeric (-)",
				),
			),
		},
		"status": {
			Description: "Server status.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"spot_server": {
			Description: "Enable if this to create kubernetes servers with spot instances",
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"spot_server_max_price": {
			Description: "The maximum price you are willing to pay for the spot instance (USD) - Any changes made to this attribute after project creation are ignored by terraform provider.  If not specified, the current on-demand price is used.",
			Type:        schema.TypeFloat,
			Optional:    true,
		},
		"zone": {
			Description:      "Availability zone for this server (only for AWS, Azure and GCP). If not specified, the first valid zone is used.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: utils.IgnoreChangeFromEmpty,
			Default:          "",
		},
		"hypervisor": {
			Description: "Hypervisor used for this server from Proxmox/vSphere Cloud credential (required for Proxmox, required for vSphere when DRS is disabled).",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
		},
	}
}

func resourceTaikunProjectSetServers(d *schema.ResourceData, apiClient *tk.Client, projectID int32) error {

	bastions := d.Get("server_bastion")
	kubeMasters := d.Get("server_kubemaster")
	kubeWorkers := d.Get("server_kubeworker")

	// Bastion
	bastion := bastions.(*schema.Set).List()[0].(map[string]interface{})
	serverCreateBody := tkcore.ServerForCreateDto{}
	serverCreateBody.SetCount(1)
	serverCreateBody.SetDiskSize(utils.GibiByteToByte(bastion["disk_size"].(int)))
	serverCreateBody.SetFlavor(bastion["flavor"].(string))
	serverCreateBody.SetName(bastion["name"].(string))
	serverCreateBody.SetAvailabilityZone(bastion["zone"].(string))
	serverCreateBody.SetHypervisor(bastion["hypervisor"].(string))
	serverCreateBody.SetProjectId(projectID)
	serverCreateBody.SetRole(tkcore.CLOUDROLE_BASTION)
	serverCreateBody, err := resourceTaikunProjectSetServerSpots(bastion, serverCreateBody) // Spots
	if err != nil {
		return err
	}

	serverCreateResponse, res, err := apiClient.Client.ServersAPI.ServersCreate(context.TODO()).ServerForCreateDto(serverCreateBody).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	bastion["id"] = serverCreateResponse.GetId()
	err = d.Set("server_bastion", []map[string]interface{}{bastion})
	if err != nil {
		return err
	}

	kubeMastersList := kubeMasters.(*schema.Set).List()
	for _, kubeMaster := range kubeMastersList {
		kubeMasterMap := kubeMaster.(map[string]interface{})
		serverCreateBody = tkcore.ServerForCreateDto{}
		serverCreateBody.SetCount(1)
		serverCreateBody.SetDiskSize(utils.GibiByteToByte(kubeMasterMap["disk_size"].(int)))
		serverCreateBody.SetFlavor(kubeMasterMap["flavor"].(string))
		serverCreateBody.SetKubernetesNodeLabels(resourceTaikunProjectServerKubernetesLabels(kubeMasterMap))
		serverCreateBody.SetName(kubeMasterMap["name"].(string))
		serverCreateBody.SetProjectId(projectID)
		serverCreateBody.SetWasmEnabled(kubeMasterMap["wasm"].(bool))
		serverCreateBody.SetAvailabilityZone(kubeMasterMap["zone"].(string))
		serverCreateBody.SetHypervisor(kubeMasterMap["hypervisor"].(string))
		serverCreateBody.SetRole(tkcore.CLOUDROLE_KUBEMASTER)
		serverCreateBody, err = resourceTaikunProjectSetServerSpots(kubeMasterMap, serverCreateBody) // Spots
		if err != nil {
			return err
		}

		serverCreateResponse, res, newErr := apiClient.Client.ServersAPI.ServersCreate(context.TODO()).ServerForCreateDto(serverCreateBody).Execute()
		if newErr != nil {
			return tk.CreateError(res, newErr)
		}
		kubeMasterMap["id"] = serverCreateResponse.GetId()
	}
	err = d.Set("server_kubemaster", kubeMastersList)
	if err != nil {
		return err
	}

	kubeWorkersList := kubeWorkers.(*schema.Set).List()
	for _, kubeWorker := range kubeWorkersList {
		kubeWorkerMap := kubeWorker.(map[string]interface{})
		serverCreateBody.SetCount(1)
		serverCreateBody.SetDiskSize(utils.GibiByteToByte(kubeWorkerMap["disk_size"].(int)))
		serverCreateBody.SetFlavor(kubeWorkerMap["flavor"].(string))
		//serverCreateBody.SetKubernetesNodeLabels(resourceTaikunProjectServerKubernetesLabels(kubeWorkerMap))
		serverCreateBody.SetName(kubeWorkerMap["name"].(string))
		serverCreateBody.SetProjectId(projectID)
		serverCreateBody.SetWasmEnabled(kubeWorkerMap["wasm"].(bool))
		serverCreateBody.SetAvailabilityZone(kubeWorkerMap["zone"].(string))
		serverCreateBody.SetHypervisor(kubeWorkerMap["hypervisor"].(string))

		if kubeWorkerMap["proxmox_extra_disk_size"].(int) != 0 {
			proxmoxStorageString, err1 := utils.GetProxmoxStorageStringForServer(projectID, apiClient)
			if err1 != nil {
				return err1
			}
			proxmoxRole, err2 := tkcore.NewProxmoxRoleFromValue(proxmoxStorageString)
			if err2 != nil {
				return err2
			}
			proxmoxExtraDiskSize := int32(kubeWorkerMap["proxmox_extra_disk_size"].(int))
			serverCreateBody.SetProxmoxRole(*proxmoxRole)
			serverCreateBody.SetProxmoxExtraDiskSize(proxmoxExtraDiskSize)
		}

		serverCreateBody.SetRole(tkcore.CLOUDROLE_KUBEWORKER)
		serverCreateBody, err = resourceTaikunProjectSetServerSpots(kubeWorkerMap, serverCreateBody) // Spots
		if err != nil {
			return err
		}

		serverCreateResponse, res, newErr := apiClient.Client.ServersAPI.ServersCreate(context.TODO()).ServerForCreateDto(serverCreateBody).Execute()
		if newErr != nil {
			return tk.CreateError(res, newErr)
		}
		kubeWorkerMap["id"] = serverCreateResponse.GetId()
	}
	err = d.Set("server_kubeworker", kubeWorkersList)
	if err != nil {
		return err
	}

	return nil
}

// Kubernetes server spots
func resourceTaikunProjectSetServerSpots(serverMap map[string]interface{}, serverCreateBody tkcore.ServerForCreateDto) (tkcore.ServerForCreateDto, error) {
	if (serverMap["spot_server_max_price"].(float64) != 0) && (!serverMap["spot_server"].(bool)) {
		return serverCreateBody, fmt.Errorf("Spot server max price is set, but the server does not have spot enabled.")
	}
	if serverMap["spot_server"] != nil {
		spotForThisVm := serverMap["spot_server"].(bool)
		serverCreateBody.SetSpotInstance(spotForThisVm)
		if serverMap["spot_server_max_price"].(float64) == 0 {
			serverCreateBody.UnsetSpotPrice() // Send null if the user did not specify anything
		} else {
			serverCreateBody.SetSpotPrice(serverMap["spot_server_max_price"].(float64))
		}
	}
	return serverCreateBody, nil
}

func resourceTaikunProjectCommit(apiClient *tk.Client, projectID int32) error {
	commitCommand := &tkcore.ProjectDeploymentCommitCommand{ProjectId: &projectID}
	res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentCommit(context.TODO()).ProjectDeploymentCommitCommand(*commitCommand).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	return nil
}

func resourceTaikunProjectPurgeServers(serversToPurge []interface{}, apiClient *tk.Client, projectID int32) error {
	serverIds := make([]int32, 0)

	for _, server := range serversToPurge {
		serverMap := server.(map[string]interface{})
		if serverIdStr, serverIdSet := serverMap["id"]; serverIdSet {
			serverId, _ := utils.Atoi32(serverIdStr.(string))
			if serverId != 0 {
				serverIds = append(serverIds, serverId)
			}
		}
	}

	if len(serverIds) != 0 {
		deleteServerBody := tkcore.ProjectDeploymentDeleteServersCommand{}
		deleteServerBody.SetProjectId(projectID)
		deleteServerBody.SetServerIds(serverIds)
		deleteServerBody.SetForceDeleteVClusters(true)
		deleteServerBody.SetDeleteAutoscalingServers(true)

		res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentDelete(context.TODO()).ProjectDeploymentDeleteServersCommand(deleteServerBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	return nil
}

func resourceTaikunProjectServerKubernetesLabels(data map[string]interface{}) []tkcore.KubernetesNodeLabelsDto {
	labels, labelsAreSet := data["kubernetes_node_label"]
	if !labelsAreSet {
		return []tkcore.KubernetesNodeLabelsDto{}
	}
	labelsList := labels.(*schema.Set).List()
	labelsToAdd := make([]tkcore.KubernetesNodeLabelsDto, len(labelsList))
	for i, labelData := range labelsList {
		label := labelData.(map[string]interface{})
		labelsToAdd[i] = tkcore.KubernetesNodeLabelsDto{}
		fmt.Println(label)
		labelsToAdd[i].SetKey(*label["key"].(tkcore.NullableString).Get())
		labelsToAdd[i].SetValue(*label["value"].(tkcore.NullableString).Get())
	}
	return labelsToAdd
}

func resourceTaikunProjectUpdateToggleServices(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	if err := resourceTaikunProjectUpdateToggleMonitoring(ctx, d, apiClient); err != nil {
		return err
	}
	if err := resourceTaikunProjectUpdateToggleBackup(ctx, d, apiClient); err != nil {
		return err
	}
	if err := resourceTaikunProjectUpdateToggleOPA(ctx, d, apiClient); err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectUpdateToggleMonitoring(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	//if d.HasChange("monitoring") {
	//	projectID, _ := utils.Atoi32(d.Id())
	//	//body := tkcore.MonitoringOperationsCommand{}
	//	//body.SetProjectId(projectID)
	//	//res, err := apiClient.Client.ProjectsAPI.ProjectsMonitoring(ctx).MonitoringOperationsCommand(body).Execute()
	//	body := tkcore.ApiProjectDeploymentEnableMonitoringRequest{}
	//	body.SetProjectId(projectID)
	//	res, err := apiClient.Client.ProjectsAPI.ProjectsMonitoring(ctx).MonitoringOperationsCommand(body).Execute()
	//	if err != nil {
	//		return tk.CreateError(res, err)
	//	}
	//
	//	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableMonitoring", "DisableMonitoring"}, apiClient, projectID); err != nil {
	//		return err
	//	}
	//}
	//return nil
	if d.HasChange("monitoring") {
		projectID, _ := utils.Atoi32(d.Id())

		// Get the current state of monitoring. If its already disabled, skip disabling query.
		data, response, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
		project := data.GetProject()
		monitoringCurrentyEnabled := project.GetIsMonitoringEnabled()
		if monitoringCurrentyEnabled {
			disableBody := tkcore.DeploymentDisableMonitoringCommand{}
			disableBody.SetProjectId(projectID)
			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentDisableMonitoring(context.TODO()).DeploymentDisableMonitoringCommand(disableBody).Execute()
			if err != nil {
				return tk.CreateError(res, err)
			}
		}

		_, newMonitoringIsSet := d.GetOk("monitoring")

		if newMonitoringIsSet {
			// Wait for the backup to be disabled
			disableStateConf := &retry.StateChangeConf{
				Pending: []string{
					strconv.FormatBool(true),
				},
				Target: []string{
					strconv.FormatBool(false),
				},
				Refresh: func() (interface{}, string, error) {
					response, _, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
					if err != nil {
						return 0, "", err
					}
					project := response.GetProject()

					return response, strconv.FormatBool(project.GetIsBackupEnabled()), nil
				},
				Timeout:                   5 * time.Minute,
				Delay:                     2 * time.Second,
				MinTimeout:                5 * time.Second,
				ContinuousTargetOccurence: 1,
			}
			_, err := disableStateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmt.Errorf("error waiting for project (%s) to enable monitoring: %s", d.Id(), err)
			}

			enableBody := tkcore.DeploymentEnableMonitoringCommand{}
			enableBody.SetProjectId(projectID)
			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentEnableMonitoring(context.TODO()).DeploymentEnableMonitoringCommand(enableBody).Execute()
			if err != nil {
				return tk.CreateError(res, err)
			}
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableBackup", "DisableBackup"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectUpdateToggleBackup(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	if d.HasChange("backup_credential_id") {
		projectID, _ := utils.Atoi32(d.Id())

		// Get the current state of backups. If they are already disabled, skip disabling query.
		data, response, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
		project := data.GetProject()
		backupCurrentyEnabled := project.GetIsBackupEnabled()
		if backupCurrentyEnabled {
			disableBody := tkcore.DeploymentDisableBackupCommand{}
			disableBody.SetProjectId(projectID)
			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentDisableBackup(context.TODO()).DeploymentDisableBackupCommand(disableBody).Execute()
			if err != nil {
				return tk.CreateError(res, err)
			}
		}

		newCredential, newCredentialIsSet := d.GetOk("backup_credential_id")

		if newCredentialIsSet {

			newCredentialID, _ := utils.Atoi32(newCredential.(string))

			// Wait for the backup to be disabled
			disableStateConf := &retry.StateChangeConf{
				Pending: []string{
					strconv.FormatBool(true),
				},
				Target: []string{
					strconv.FormatBool(false),
				},
				Refresh: func() (interface{}, string, error) {
					response, _, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
					if err != nil {
						return 0, "", err
					}
					project := response.GetProject()

					return response, strconv.FormatBool(project.GetIsBackupEnabled()), nil
				},
				Timeout:                   5 * time.Minute,
				Delay:                     2 * time.Second,
				MinTimeout:                5 * time.Second,
				ContinuousTargetOccurence: 1,
			}
			_, err := disableStateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmt.Errorf("error waiting for project (%s) to disable backup: %s", d.Id(), err)
			}

			enableBody := tkcore.DeploymentEnableBackupCommand{}
			enableBody.SetProjectId(projectID)
			enableBody.SetS3CredentialId(newCredentialID)
			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentEnableBackup(context.TODO()).DeploymentEnableBackupCommand(enableBody).Execute()
			if err != nil {
				return tk.CreateError(res, err)
			}
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableBackup", "DisableBackup"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectUpdateToggleOPA(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	if d.HasChange("policy_profile_id") {
		projectID, _ := utils.Atoi32(d.Id())
		oldOPAProfile, _ := d.GetChange("policy_profile_id")

		if oldOPAProfile != "" {

			disableBody := tkcore.DeploymentDisableOpaCommand{}
			disableBody.SetProjectId(projectID)

			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentDisableOpa(ctx).DeploymentDisableOpaCommand(disableBody).Execute()

			if err != nil {
				return tk.CreateError(res, err)
			}

		}

		newOPAProfile, newOPAProfileIsSet := d.GetOk("policy_profile_id")

		if newOPAProfileIsSet {

			newOPAProfilelID, _ := utils.Atoi32(newOPAProfile.(string))

			// Wait for the OPA to be disabled
			disableStateConf := &retry.StateChangeConf{
				Pending: []string{
					strconv.FormatBool(true),
				},
				Target: []string{
					strconv.FormatBool(false),
				},
				Refresh: func() (interface{}, string, error) {
					response, res, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
					if err != nil {
						return 0, "", tk.CreateError(res, err)
					}

					project := response.GetProject()
					return response, strconv.FormatBool(project.GetIsOpaEnabled()), nil
				},
				Timeout:                   5 * time.Minute,
				Delay:                     2 * time.Second,
				MinTimeout:                5 * time.Second,
				ContinuousTargetOccurence: 1,
			}
			_, err := disableStateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmt.Errorf("error waiting for project (%s) to disable OPA: %s", d.Id(), err)
			}

			enableBody := tkcore.DeploymentOpaEnableCommand{}
			enableBody.SetProjectId(projectID)
			enableBody.SetOpaCredentialId(newOPAProfilelID)

			res, err := apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentEnableOpa(ctx).DeploymentOpaEnableCommand(enableBody).Execute()
			if err != nil {
				return tk.CreateError(res, err)
			}
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableGatekeeper", "DisableGatekeeper"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectEditFlavors(d *schema.ResourceData, apiClient *tk.Client, id int32) error {
	oldFlavorData, newFlavorData := d.GetChange("flavors")
	oldFlavors := oldFlavorData.(*schema.Set)
	newFlavors := newFlavorData.(*schema.Set)
	flavorsToUnbind := oldFlavors.Difference(newFlavors)
	flavorsToBind := newFlavors.Difference(oldFlavors).List()
	boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(id, apiClient)
	if err != nil {
		return err
	}
	if flavorsToUnbind.Len() != 0 {
		var flavorBindingsToUndo []int32
		for _, boundFlavorDTO := range boundFlavorDTOs {
			if flavorsToUnbind.Contains(boundFlavorDTO.GetName()) {
				flavorBindingsToUndo = append(flavorBindingsToUndo, boundFlavorDTO.GetId())
			}
		}
		unbindBody := tkcore.UnbindFlavorFromProjectCommand{}
		unbindBody.SetIds(flavorBindingsToUndo)
		res, err := apiClient.Client.FlavorsAPI.FlavorsUnbindFromProject(context.TODO()).UnbindFlavorFromProjectCommand(unbindBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	if len(flavorsToBind) != 0 {
		flavorsToBindNames := make([]string, len(flavorsToBind))
		for i, flavorToBind := range flavorsToBind {
			flavorsToBindNames[i] = flavorToBind.(string)
		}
		bindBody := tkcore.BindFlavorToProjectCommand{}
		bindBody.SetProjectId(id)
		bindBody.SetFlavors(flavorsToBindNames)
		res, err := apiClient.Client.FlavorsAPI.FlavorsBindToProject(context.TODO()).BindFlavorToProjectCommand(bindBody).Execute()
		if err != nil {
			return tk.CreateError(res, err)
		}
	}
	return nil
}

func resourceTaikunProjectUpdateAutoscaler(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	body := tkcore.EditAutoscalingCommand{}
	body.SetProjectId(projectID)
	body.SetMinSize(int32(d.Get("autoscaler_min_size").(int)))
	body.SetMaxSize(int32(d.Get("autoscaler_max_size").(int)))

	res, err := apiClient.Client.AutoscalingAPI.AutoscalingEdit(ctx).EditAutoscalingCommand(body).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableAutoscaler", "DisableAutoscaler"}, apiClient, projectID); err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectRecreateAutoscaler(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	// Is autoscaler enabled or disabled?
	projectID, _ := utils.Atoi32(d.Id())
	data, response, err := apiClient.Client.ServersAPI.ServersDetails(ctx, projectID).Execute()
	if err != nil {
		return tk.CreateError(response, err)
	}
	if *data.GetProject().IsAutoscalingEnabled {
		// Autoscaler was enabled -> Disable autoscaler
		err := resourceTaikunProjectDisableAutoscaler(ctx, d, apiClient)
		if err != nil {
			return err
		}
	}
	// else autoscaler was disabled -> keep calm and carry on

	// Enable autoscaler with new values
	err = resourceTaikunProjectEnableAutoscaler(ctx, d, apiClient)
	if err != nil {
		return err
	}

	return nil
}

func resourceTaikunProjectDisableAutoscaler(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	bodyDisable := tkcore.DisableAutoscalingCommand{}
	bodyDisable.SetProjectId(projectID)
	res, err := apiClient.Client.AutoscalingAPI.AutoscalingDisable(ctx).DisableAutoscalingCommand(bodyDisable).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableAutoscaler", "DisableAutoscaler"}, apiClient, projectID); err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectEnableAutoscaler(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	bodyEnable := tkcore.EnableAutoscalingCommand{}
	bodyEnable.SetId(projectID)
	bodyEnable.SetAutoscalingGroupName(d.Get("autoscaler_name").(string))
	bodyEnable.SetFlavor(d.Get("autoscaler_flavor").(string))
	bodyEnable.SetMaxSize(int32(d.Get("autoscaler_max_size").(int)))
	bodyEnable.SetMinSize(int32(d.Get("autoscaler_min_size").(int)))
	bodyEnable.SetDiskSize(float64(utils.GibiByteToByte(d.Get("autoscaler_disk_size").(int))))
	bodyEnable.SetSpotEnabled(d.Get("autoscaler_spot_enabled").(bool))

	res, err := apiClient.Client.AutoscalingAPI.AutoscalingEnable(ctx).EnableAutoscalingCommand(bodyEnable).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableAutoscaler", "DisableAutoscaler"}, apiClient, projectID); err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectToggleFullSpot(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	bodyToggle := tkcore.FullSpotOperationCommand{}
	bodyToggle.SetId(projectID)

	if d.Get("spot_full").(bool) {
		bodyToggle.SetMode("enable")
	} else if !d.Get("spot_full").(bool) {
		bodyToggle.SetMode("disable")
	}

	res, err := apiClient.Client.ProjectsAPI.ProjectsToggleFullSpot(ctx).FullSpotOperationCommand(bodyToggle).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	return nil
}

func resourceTaikunProjectToggleWorkerSpot(ctx context.Context, d *schema.ResourceData, apiClient *tk.Client) error {
	projectID, _ := utils.Atoi32(d.Id())
	bodyToggle := tkcore.SpotWorkerOperationCommand{}
	bodyToggle.SetId(projectID)

	if d.Get("spot_worker").(bool) {
		bodyToggle.SetMode("enable")
	} else if !d.Get("spot_full").(bool) {
		bodyToggle.SetMode("disable")
	}

	res, err := apiClient.Client.ProjectsAPI.ProjectsToggleSpotWorkers(ctx).SpotWorkerOperationCommand(bodyToggle).Execute()
	if err != nil {
		return tk.CreateError(res, err)
	}
	return nil
}
