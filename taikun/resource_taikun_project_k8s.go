package taikun

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/backup"
	"github.com/itera-io/taikungoclient/client/flavors"
	"github.com/itera-io/taikungoclient/client/opa_profiles"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/itera-io/taikungoclient/models"
)

func taikunServerKubeworkerSchema() map[string]*schema.Schema {
	kubeworkerSchema := taikunServerSchemaWithKubernetesNodeLabels()
	removeForceNewsFromSchema(kubeworkerSchema)
	return kubeworkerSchema
}

func taikunServerSchemaWithKubernetesNodeLabels() map[string]*schema.Schema {
	serverSchema := taikunServerBasicSchema()
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
	}
}

func resourceTaikunProjectSetServers(d *schema.ResourceData, apiClient *taikungoclient.Client, projectID int32) error {

	bastions := d.Get("server_bastion")
	kubeMasters := d.Get("server_kubemaster")
	kubeWorkers := d.Get("server_kubeworker")

	// Bastion
	bastion := bastions.(*schema.Set).List()[0].(map[string]interface{})
	serverCreateBody := &models.ServerForCreateDto{
		Count:                1,
		DiskSize:             gibiByteToByte(bastion["disk_size"].(int)),
		Flavor:               bastion["flavor"].(string),
		KubernetesNodeLabels: nil,
		Name:                 bastion["name"].(string),
		ProjectID:            projectID,
		Role:                 100,
	}

	serverCreateParams := servers.NewServersCreateParams().WithV(ApiVersion).WithBody(serverCreateBody)
	serverCreateResponse, err := apiClient.Client.Servers.ServersCreate(serverCreateParams, apiClient)
	if err != nil {
		return err
	}
	bastion["id"] = serverCreateResponse.Payload.ID
	err = d.Set("server_bastion", []map[string]interface{}{bastion})
	if err != nil {
		return err
	}

	kubeMastersList := kubeMasters.(*schema.Set).List()
	for _, kubeMaster := range kubeMastersList {
		kubeMasterMap := kubeMaster.(map[string]interface{})

		serverCreateBody := &models.ServerForCreateDto{
			Count:                1,
			DiskSize:             gibiByteToByte(kubeMasterMap["disk_size"].(int)),
			Flavor:               kubeMasterMap["flavor"].(string),
			KubernetesNodeLabels: resourceTaikunProjectServerKubernetesLabels(kubeMasterMap),
			Name:                 kubeMasterMap["name"].(string),
			ProjectID:            projectID,
			Role:                 200,
		}

		serverCreateParams := servers.NewServersCreateParams().WithV(ApiVersion).WithBody(serverCreateBody)
		serverCreateResponse, err := apiClient.Client.Servers.ServersCreate(serverCreateParams, apiClient)
		if err != nil {
			return err
		}
		kubeMasterMap["id"] = serverCreateResponse.Payload.ID
	}
	err = d.Set("server_kubemaster", kubeMastersList)
	if err != nil {
		return err
	}

	kubeWorkersList := kubeWorkers.(*schema.Set).List()
	for _, kubeWorker := range kubeWorkersList {
		kubeWorkerMap := kubeWorker.(map[string]interface{})

		serverCreateBody := &models.ServerForCreateDto{
			Count:                1,
			DiskSize:             gibiByteToByte(kubeWorkerMap["disk_size"].(int)),
			Flavor:               kubeWorkerMap["flavor"].(string),
			KubernetesNodeLabels: resourceTaikunProjectServerKubernetesLabels(kubeWorkerMap),
			Name:                 kubeWorkerMap["name"].(string),
			ProjectID:            projectID,
			Role:                 300,
		}
		serverCreateParams := servers.NewServersCreateParams().WithV(ApiVersion).WithBody(serverCreateBody)
		serverCreateResponse, err := apiClient.Client.Servers.ServersCreate(serverCreateParams, apiClient)
		if err != nil {
			return err
		}
		kubeWorkerMap["id"] = serverCreateResponse.Payload.ID
	}
	err = d.Set("server_kubeworker", kubeWorkersList)
	if err != nil {
		return err
	}

	return nil
}

func resourceTaikunProjectCommit(apiClient *taikungoclient.Client, projectID int32) error {
	params := projects.NewProjectsCommitParams().WithV(ApiVersion).WithProjectID(projectID)
	_, err := apiClient.Client.Projects.ProjectsCommit(params, apiClient)
	if err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectPurgeServers(serversToPurge []interface{}, apiClient *taikungoclient.Client, projectID int32) error {
	serverIds := make([]int32, 0)

	for _, server := range serversToPurge {
		serverMap := server.(map[string]interface{})
		if serverIdStr, serverIdSet := serverMap["id"]; serverIdSet {
			serverId, _ := atoi32(serverIdStr.(string))
			if serverId != 0 {
				serverIds = append(serverIds, serverId)
			}
		}
	}

	if len(serverIds) != 0 {
		deleteServerBody := &models.DeleteServerCommand{
			ProjectID: projectID,
			ServerIds: serverIds,
		}
		deleteServerParams := servers.NewServersDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
		_, _, err := apiClient.Client.Servers.ServersDelete(deleteServerParams, apiClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectServerKubernetesLabels(data map[string]interface{}) []*models.KubernetesNodeLabelsDto {
	labels, labelsAreSet := data["kubernetes_node_label"]
	if !labelsAreSet {
		return []*models.KubernetesNodeLabelsDto{}
	}
	labelsList := labels.(*schema.Set).List()
	labelsToAdd := make([]*models.KubernetesNodeLabelsDto, len(labelsList))
	for i, labelData := range labelsList {
		label := labelData.(map[string]interface{})
		labelsToAdd[i] = &models.KubernetesNodeLabelsDto{
			Key:   label["key"].(string),
			Value: label["value"].(string),
		}
	}
	return labelsToAdd
}

func resourceTaikunProjectUpdateToggleServices(ctx context.Context, d *schema.ResourceData, apiClient *taikungoclient.Client) error {
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

func resourceTaikunProjectUpdateToggleMonitoring(ctx context.Context, d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	if d.HasChange("monitoring") {
		projectID, _ := atoi32(d.Id())
		body := models.MonitoringOperationsCommand{ProjectID: projectID}
		params := projects.NewProjectsMonitoringOperationsParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.Client.Projects.ProjectsMonitoringOperations(params, apiClient)
		if err != nil {
			return err
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableMonitoring", "DisableMonitoring"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectUpdateToggleBackup(ctx context.Context, d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	if d.HasChange("backup_credential_id") {
		projectID, _ := atoi32(d.Id())

		disableBody := &models.DisableBackupCommand{
			ProjectID: projectID,
		}
		disableParams := backup.NewBackupDisableBackupParams().WithV(ApiVersion).WithBody(disableBody)
		_, err := apiClient.Client.Backup.BackupDisableBackup(disableParams, apiClient)
		if err != nil {
			return err
		}

		newCredential, newCredentialIsSet := d.GetOk("backup_credential_id")

		if newCredentialIsSet {

			newCredentialID, _ := atoi32(newCredential.(string))

			// Wait for the backup to be disabled
			disableStateConf := &resource.StateChangeConf{
				Pending: []string{
					strconv.FormatBool(true),
				},
				Target: []string{
					strconv.FormatBool(false),
				},
				Refresh: func() (interface{}, string, error) {
					params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
					response, err := apiClient.Client.Servers.ServersDetails(params, apiClient)
					if err != nil {
						return 0, "", err
					}

					return response, strconv.FormatBool(response.Payload.Project.IsBackupEnabled), nil
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

			enableBody := &models.EnableBackupCommand{
				ProjectID:      projectID,
				S3CredentialID: newCredentialID,
			}
			enableParams := backup.NewBackupEnableBackupParams().WithV(ApiVersion).WithBody(enableBody)
			_, err = apiClient.Client.Backup.BackupEnableBackup(enableParams, apiClient)
			if err != nil {
				return err
			}
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableBackup", "DisableBackup"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectUpdateToggleOPA(ctx context.Context, d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	if d.HasChange("policy_profile_id") {
		projectID, _ := atoi32(d.Id())
		oldOPAProfile, _ := d.GetChange("policy_profile_id")

		if oldOPAProfile != "" {

			disableBody := &models.DisableGatekeeperCommand{
				ProjectID: projectID,
			}
			disableParams := opa_profiles.NewOpaProfilesDisableGatekeeperParams().WithV(ApiVersion).WithBody(disableBody)
			_, err := apiClient.Client.OpaProfiles.OpaProfilesDisableGatekeeper(disableParams, apiClient)
			if err != nil {
				return err
			}

		}

		newOPAProfile, newOPAProfileIsSet := d.GetOk("policy_profile_id")

		if newOPAProfileIsSet {

			newOPAProfilelID, _ := atoi32(newOPAProfile.(string))

			// Wait for the OPA to be disabled
			disableStateConf := &resource.StateChangeConf{
				Pending: []string{
					strconv.FormatBool(true),
				},
				Target: []string{
					strconv.FormatBool(false),
				},
				Refresh: func() (interface{}, string, error) {
					params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
					response, err := apiClient.Client.Servers.ServersDetails(params, apiClient)
					if err != nil {
						return 0, "", err
					}

					return response, strconv.FormatBool(response.Payload.Project.IsOpaEnabled), nil
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

			enableBody := &models.EnableGatekeeperCommand{
				ProjectID:    projectID,
				OpaProfileID: newOPAProfilelID,
			}
			enableParams := opa_profiles.NewOpaProfilesEnableGatekeeperParams().WithV(ApiVersion).WithBody(enableBody)
			_, err = apiClient.Client.OpaProfiles.OpaProfilesEnableGatekeeper(enableParams, apiClient)
			if err != nil {
				return err
			}
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableGatekeeper", "DisableGatekeeper"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectEditFlavors(d *schema.ResourceData, apiClient *taikungoclient.Client, id int32) error {
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
			if flavorsToUnbind.Contains(boundFlavorDTO.Name) {
				flavorBindingsToUndo = append(flavorBindingsToUndo, boundFlavorDTO.ID)
			}
		}
		unbindBody := models.UnbindFlavorFromProjectCommand{Ids: flavorBindingsToUndo}
		unbindParams := flavors.NewFlavorsUnbindFromProjectParams().WithV(ApiVersion).WithBody(&unbindBody)
		if _, err := apiClient.Client.Flavors.FlavorsUnbindFromProject(unbindParams, apiClient); err != nil {
			return err
		}
	}
	if len(flavorsToBind) != 0 {
		flavorsToBindNames := make([]string, len(flavorsToBind))
		for i, flavorToBind := range flavorsToBind {
			flavorsToBindNames[i] = flavorToBind.(string)
		}
		bindBody := models.BindFlavorToProjectCommand{ProjectID: id, Flavors: flavorsToBindNames}
		bindParams := flavors.NewFlavorsBindToProjectParams().WithV(ApiVersion).WithBody(&bindBody)
		if _, err := apiClient.Client.Flavors.FlavorsBindToProject(bindParams, apiClient); err != nil {
			return err
		}
	}
	return nil
}
