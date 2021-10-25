package taikun

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/itera-io/taikungoclient/client/backup"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_profile_id": {
			Description:      "ID of the project's access profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"alerting_profile_id": {
			Description:      "ID of the project's alerting profile.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"alerting_profile_name": {
			Description: "Name of the project's alerting profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"backup_credential_id": {
			Description:      "ID of the backup credential. If unspecified, backups are disabled.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"enable_auto_upgrade": {
			Description: "Kubespray version will be automatically upgraded if new version is available.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"enable_monitoring": {
			Description: "Kubernetes cluster monitoring.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to store the project.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"expiration_date": {
			Description:      "Project's expiration date in the format: 'dd/mm/yyyy'.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsDate,
		},
		"id": {
			Description: "Project ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"kubernetes_profile_id": {
			Description:      "ID of the project's kubernetes profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"name": {
			Description: "Project name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or non alpha numeric (-)",
				),
			),
			ForceNew: true,
		},
		"organization_id": {
			Description:      "ID of the organization which owns the project.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
	}
}

func resourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project",
		CreateContext: resourceTaikunProjectCreate,
		ReadContext:   resourceTaikunProjectRead,
		UpdateContext: resourceTaikunProjectUpdate,
		DeleteContext: resourceTaikunProjectDelete,
		Schema:        resourceTaikunProjectSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunProjectRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id := data.Id()
	id32, err := atoi32(id)
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id32) // TODO use /api/v1/projects endpoint?
	response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
	if err != nil {
		return nil
	}

	projectDetailsDTO := response.Payload.Project
	err = setResourceDataFromMap(data, flattenTaikunProject(projectDetailsDTO))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(id)

	return nil
}

func resourceTaikunProjectCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateProjectCommand{
		Name:         data.Get("name").(string),
		IsKubernetes: true,
	}
	body.CloudCredentialID, _ = atoi32(data.Get("cloud_credential_id").(string))
	if accessProfileID, accessProfileIDIsSet := data.GetOk("access_profile_id"); accessProfileIDIsSet {
		body.AccessProfileID, _ = atoi32(accessProfileID.(string))
	}
	if alertingProfileID, alertingProfileIDIsSet := data.GetOk("alerting_profile_id"); alertingProfileIDIsSet {
		body.AlertingProfileID, _ = atoi32(alertingProfileID.(string))
	}
	if backupCredentialID, backupCredentialIDIsSet := data.GetOk("backup_credential_id"); backupCredentialIDIsSet {
		body.IsBackupEnabled = true
		body.S3CredentialID, _ = atoi32(backupCredentialID.(string))
	}
	if enableAutoUpgrade, enableAutoUpgradeIsSet := data.GetOk("enable_auto_upgrade"); enableAutoUpgradeIsSet {
		body.IsAutoUpgrade = enableAutoUpgrade.(bool)
	}
	if enableMonitoring, enableMonitoringIsSet := data.GetOk("enable_monitoring"); enableMonitoringIsSet {
		body.IsMonitoringEnabled = enableMonitoring.(bool)
	}
	if expirationDate, expirationDateIsSet := data.GetOk("expiration_date"); expirationDateIsSet {
		dateTime := dateToDateTime(expirationDate.(string))
		body.ExpiredAt = &dateTime
	} else {
		body.ExpiredAt = nil
	}
	if kubernetesProfileID, kubernetesProfileIDIsSet := data.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		body.KubernetesProfileID, _ = atoi32(kubernetesProfileID.(string))
	}
	if organizationID, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		body.OrganizationID, _ = atoi32(organizationID.(string))
	}

	params := projects.NewProjectsCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Projects.ProjectsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.Payload.ID)

	return resourceTaikunProjectRead(ctx, data, meta)
}

func resourceTaikunProjectUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("alerting_profile_id") {
		body := models.AttachDetachAlertingProfileCommand{
			ProjectID: id,
		}
		detachParams := alerting_profiles.NewAlertingProfilesDetachParams().WithV(ApiVersion).WithBody(&body)
		if _, err := apiClient.client.AlertingProfiles.AlertingProfilesDetach(detachParams, apiClient); err != nil {
			return diag.FromErr(err)
		}
		if newAlertingProfileIDData, newAlertingProfileIDProvided := data.GetOk("alerting_profile_id"); newAlertingProfileIDProvided {
			newAlertingProfileID, _ := atoi32(newAlertingProfileIDData.(string))
			body.AlertingProfileID = newAlertingProfileID
			attachParams := alerting_profiles.NewAlertingProfilesAttachParams().WithV(ApiVersion).WithBody(&body)
			if _, err := apiClient.client.AlertingProfiles.AlertingProfilesAttach(attachParams, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if data.HasChange("enable_monitoring") {
		body := models.MonitoringOperationsCommand{ProjectID: id}
		params := projects.NewProjectsMonitoringOperationsParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.Projects.ProjectsMonitoringOperations(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if data.HasChange("backup_credential_id") {
		oldCredential, _ := data.GetChange("backup_credential_id")

		if oldCredential != "" {

			oldCredentialID, _ := atoi32(oldCredential.(string))

			disableBody := &models.DisableBackupCommand{
				ProjectID:      id,
				S3CredentialID: oldCredentialID,
			}
			disableParams := backup.NewBackupDisableBackupParams().WithV(ApiVersion).WithBody(disableBody)
			_, err = apiClient.client.Backup.BackupDisableBackup(disableParams, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		newCredential, newCredentialIsSet := data.GetOk("backup_credential_id")

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
					params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id) // TODO use /api/v1/projects endpoint?
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
			_, err = disableStateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf("Error waiting for project (%s) to disable backup: %s", data.Id(), err)
			}

			enableBody := &models.EnableBackupCommand{
				ProjectID:      id,
				S3CredentialID: newCredentialID,
			}
			enableParams := backup.NewBackupEnableBackupParams().WithV(ApiVersion).WithBody(enableBody)
			_, err = apiClient.client.Backup.BackupEnableBackup(enableParams, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if data.HasChange("expiration_date") {
		body := models.ProjectExtendLifeTimeCommand{
			ProjectID: id,
		}
		if expirationDate, expirationDateIsSet := data.GetOk("expiration_date"); expirationDateIsSet {
			dateTime := dateToDateTime(expirationDate.(string))
			body.ExpireAt = &dateTime
		} else {
			body.ExpireAt = nil
		}
		params := projects.NewProjectsExtendLifeTimeParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.Projects.ProjectsExtendLifeTime(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunProjectRead(ctx, data, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := models.DeleteProjectCommand{ProjectID: id, IsForceDelete: false}
	params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.client.Projects.ProjectsDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

// TODO change type of DTO if read endpoint is modified
func flattenTaikunProject(projectDetailsDTO *models.ProjectDetailsForServersDto) map[string]interface{} {
	projectMap := map[string]interface{}{
		"access_profile_id":     i32toa(projectDetailsDTO.AccessProfileID),
		"alerting_profile_name": projectDetailsDTO.AlertingProfileName,
		"cloud_credential_id":   i32toa(projectDetailsDTO.CloudID),
		"enable_auto_upgrade":   projectDetailsDTO.IsAutoUpgrade,
		"enable_monitoring":     projectDetailsDTO.IsMonitoringEnabled,
		"expiration_date":       rfc3339DateTimeToDate(projectDetailsDTO.ExpiredAt),
		"id":                    i32toa(projectDetailsDTO.ProjectID),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.KubernetesProfileID),
		"name":                  projectDetailsDTO.ProjectName,
		"organization_id":       i32toa(projectDetailsDTO.OrganizationID),
	}

	var nullID int32
	if projectDetailsDTO.AlertingProfileID != nullID {
		projectMap["alerting_profile_id"] = i32toa(projectDetailsDTO.AlertingProfileID)
	}

	if projectDetailsDTO.IsBackupEnabled {
		projectMap["backup_credential_id"] = i32toa(projectDetailsDTO.S3CredentialID)
	}

	return projectMap
}
