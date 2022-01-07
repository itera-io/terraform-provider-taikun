package taikun

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/backup"
	"github.com/itera-io/taikungoclient/client/opa_profiles"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunProjectSetServers(d *schema.ResourceData, apiClient *apiClient, projectID int32) error {

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
	serverCreateResponse, err := apiClient.client.Servers.ServersCreate(serverCreateParams, apiClient)
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
		serverCreateResponse, err := apiClient.client.Servers.ServersCreate(serverCreateParams, apiClient)
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
		serverCreateResponse, err := apiClient.client.Servers.ServersCreate(serverCreateParams, apiClient)
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

func resourceTaikunProjectCommit(apiClient *apiClient, projectID int32) error {
	params := projects.NewProjectsCommitParams().WithV(ApiVersion).WithProjectID(projectID)
	_, err := apiClient.client.Projects.ProjectsCommit(params, apiClient)
	if err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectPurgeServers(serversToPurge []interface{}, apiClient *apiClient, projectID int32) error {
	serverIds := make([]int32, 0)

	for _, server := range serversToPurge {
		serverMap := server.(map[string]interface{})
		serverId, _ := atoi32(serverMap["id"].(string))
		serverIds = append(serverIds, serverId)
	}

	if len(serverIds) != 0 {
		deleteServerBody := &models.DeleteServerCommand{
			ProjectID: projectID,
			ServerIds: serverIds,
		}
		deleteServerParams := servers.NewServersDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
		_, _, err := apiClient.client.Servers.ServersDelete(deleteServerParams, apiClient)
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

func resourceTaikunProjectUpdateToggleServices(ctx context.Context, d *schema.ResourceData, apiClient *apiClient) error {
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

func resourceTaikunProjectUpdateToggleMonitoring(ctx context.Context, d *schema.ResourceData, apiClient *apiClient) error {
	if d.HasChange("monitoring") {
		projectID, _ := atoi32(d.Id())
		body := models.MonitoringOperationsCommand{ProjectID: projectID}
		params := projects.NewProjectsMonitoringOperationsParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.Projects.ProjectsMonitoringOperations(params, apiClient)
		if err != nil {
			return err
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"EnableMonitoring", "DisableMonitoring"}, apiClient, projectID); err != nil {
			return err
		}
	}
	return nil
}

func resourceTaikunProjectUpdateToggleBackup(ctx context.Context, d *schema.ResourceData, apiClient *apiClient) error {
	if d.HasChange("backup_credential_id") {
		projectID, _ := atoi32(d.Id())
		oldCredential, _ := d.GetChange("backup_credential_id")

		if oldCredential != "" {

			oldCredentialID, _ := atoi32(oldCredential.(string))

			disableBody := &models.DisableBackupCommand{
				ProjectID:      projectID,
				S3CredentialID: oldCredentialID,
			}
			disableParams := backup.NewBackupDisableBackupParams().WithV(ApiVersion).WithBody(disableBody)
			_, err := apiClient.client.Backup.BackupDisableBackup(disableParams, apiClient)
			if err != nil {
				return err
			}

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
					response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
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
			_, err = apiClient.client.Backup.BackupEnableBackup(enableParams, apiClient)
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

func resourceTaikunProjectUpdateToggleOPA(ctx context.Context, d *schema.ResourceData, apiClient *apiClient) error {
	if d.HasChange("policy_profile_id") {
		projectID, _ := atoi32(d.Id())
		oldOPAProfile, _ := d.GetChange("policy_profile_id")

		if oldOPAProfile != "" {

			disableBody := &models.DisableGatekeeperCommand{
				ProjectID: projectID,
			}
			disableParams := opa_profiles.NewOpaProfilesDisableGatekeeperParams().WithV(ApiVersion).WithBody(disableBody)
			_, err := apiClient.client.OpaProfiles.OpaProfilesDisableGatekeeper(disableParams, apiClient)
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
					params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID) // TODO use /api/v1/projects endpoint?
					response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
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
			_, err = apiClient.client.OpaProfiles.OpaProfilesEnableGatekeeper(enableParams, apiClient)
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
