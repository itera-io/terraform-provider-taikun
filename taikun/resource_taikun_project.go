package taikun

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/itera-io/taikungoclient/client/project_quotas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/itera-io/taikungoclient/client/backup"
	"github.com/itera-io/taikungoclient/client/flavors"

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
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to store the project.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
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
		"expiration_date": {
			Description:      "Project's expiration date in the format: 'dd/mm/yyyy'.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsDate,
		},
		"flavors": {
			Description: "List of flavors bound to the project.",
			Type:        schema.TypeSet,
			Optional:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{}, nil
			},
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
		"quota_cpu_units": {
			Description: "Maximum CPU units.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"quota_disk_size": {
			Description: "Maximum disk size in GBs. Unlimited if unspecified.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"quota_id": {
			Description: "ID of the project quota.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"quota_ram_size": {
			Description: "Maximum RAM size in GBs. Unlimited if unspecified.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"router_id_end_range": {
			Description:  "Router ID end range (only used if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 255),
			RequiredWith: []string{"router_id_start_range", "taikun_lb_flavor"},
		},
		"router_id_start_range": {
			Description:  "Router ID start range (only used if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 255),
			RequiredWith: []string{"router_id_end_range", "taikun_lb_flavor"},
		},
		"taikun_lb_flavor": {
			Description:  "OpenStack flavor for the Taikun load balancer (only used if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"router_id_end_range", "router_id_start_range"},
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

func resourceTaikunProjectCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := models.CreateProjectCommand{
		Name:         data.Get("name").(string),
		IsKubernetes: true,
	}
	body.CloudCredentialID, _ = atoi32(data.Get("cloud_credential_id").(string))
	flavorsData := data.Get("flavors").(*schema.Set).List()
	flavors := make([]string, len(flavorsData))
	for i, flavorData := range flavorsData {
		flavors[i] = flavorData.(string)
	}
	body.Flavors = flavors

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

	if taikunLBFlavor, taikunLBFlavorIsSet := data.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
		body.TaikunLBFlavor = taikunLBFlavor.(string)
		body.RouterIDStartRange = int32(data.Get("router_id_start_range").(int))
		body.RouterIDEndRange = int32(data.Get("router_id_end_range").(int))
	}

	params := projects.NewProjectsCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Projects.ProjectsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.Payload.ID)

	quotaCPU, quotaCPUIsSet := data.GetOk("quota_cpu_units")
	quotaDisk, quotaDiskIsSet := data.GetOk("quota_disk_size")
	quotaRAM, quotaRAMIsSet := data.GetOk("quota_ram_size")
	if quotaCPUIsSet || quotaDiskIsSet || quotaRAMIsSet {

		quotaEditBody := &models.ProjectQuotaUpdateDto{
			IsCPUUnlimited:      true,
			IsRAMUnlimited:      true,
			IsDiskSizeUnlimited: true,
		}

		if quotaCPUIsSet {
			quotaEditBody.CPU = int64(quotaCPU.(int))
			quotaEditBody.IsCPUUnlimited = false
		}

		if quotaDiskIsSet {
			quotaEditBody.DiskSize = int64(quotaDisk.(int))
			quotaEditBody.IsDiskSizeUnlimited = false
		}

		if quotaRAMIsSet {
			quotaEditBody.RAM = int64(quotaRAM.(int))
			quotaEditBody.IsDiskSizeUnlimited = false
		}

		projectId, _ := atoi32(response.Payload.ID)

		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectId) // TODO use /api/v1/projects endpoint?
		response, err := apiClient.client.Servers.ServersDetails(params, apiClient)

		if err == nil {
			quotaEditParams := project_quotas.NewProjectQuotasEditParams().WithV(ApiVersion).WithQuotaID(response.Payload.Project.QuotaID).WithBody(quotaEditBody)
			_, err := apiClient.client.ProjectQuotas.ProjectQuotasEdit(quotaEditParams, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return readAfterCreateWithRetries(resourceTaikunProjectRead, ctx, data, meta)
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

	boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(projectDetailsDTO.ProjectID, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	quotaParams := project_quotas.NewProjectQuotasListParams().WithV(ApiVersion).WithID(&projectDetailsDTO.QuotaID)
	quotaResponse, err := apiClient.client.ProjectQuotas.ProjectQuotasList(quotaParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(quotaResponse.Payload.Data) != 1 {
		return nil
	}

	err = setResourceDataFromMap(data, flattenTaikunProject(projectDetailsDTO, boundFlavorDTOs, quotaResponse.Payload.Data[0]))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(id)

	return nil
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
	if data.HasChange("enable_monitoring") {
		body := models.MonitoringOperationsCommand{ProjectID: id}
		params := projects.NewProjectsMonitoringOperationsParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.Projects.ProjectsMonitoringOperations(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
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
	if data.HasChange("flavors") {
		oldFlavorData, newFlavorData := data.GetChange("flavors")
		oldFlavors := oldFlavorData.(*schema.Set)
		newFlavors := newFlavorData.(*schema.Set)
		flavorsToUnbind := oldFlavors.Difference(newFlavors)
		flavorsToBind := newFlavors.Difference(oldFlavors).List()
		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(id, apiClient)
		if err != nil {
			return diag.FromErr(err)
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
			if _, err := apiClient.client.Flavors.FlavorsUnbindFromProject(unbindParams, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
		if len(flavorsToBind) != 0 {
			flavorsToBindNames := make([]string, len(flavorsToBind))
			for i, flavorToBind := range flavorsToBind {
				flavorsToBindNames[i] = flavorToBind.(string)
			}
			bindBody := models.BindFlavorToProjectCommand{ProjectID: id, Flavors: flavorsToBindNames}
			bindParams := flavors.NewFlavorsBindToProjectParams().WithV(ApiVersion).WithBody(&bindBody)
			if _, err := apiClient.client.Flavors.FlavorsBindToProject(bindParams, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if data.HasChanges("quota_cpu_units", "quota_disk_size", "quota_ram_size") {
		quotaId, _ := atoi32(data.Get("quota_id").(string))

		quotaEditBody := &models.ProjectQuotaUpdateDto{
			IsCPUUnlimited:      true,
			IsRAMUnlimited:      true,
			IsDiskSizeUnlimited: true,
		}

		if quotaCPU, quotaCPUIsSet := data.GetOk("quota_cpu_units"); quotaCPUIsSet {
			quotaEditBody.CPU = int64(quotaCPU.(int))
			quotaEditBody.IsCPUUnlimited = false
		}

		if quotaDisk, quotaDiskIsSet := data.GetOk("quota_disk_size"); quotaDiskIsSet {
			quotaEditBody.DiskSize = int64(quotaDisk.(int))
			quotaEditBody.IsDiskSizeUnlimited = false
		}

		if quotaRAM, quotaRAMIsSet := data.GetOk("quota_ram_size"); quotaRAMIsSet {
			quotaEditBody.RAM = int64(quotaRAM.(int))
			quotaEditBody.IsRAMUnlimited = false
		}

		quotaEditParams := project_quotas.NewProjectQuotasEditParams().WithV(ApiVersion).WithQuotaID(quotaId).WithBody(quotaEditBody)
		_, err := apiClient.client.ProjectQuotas.ProjectQuotasEdit(quotaEditParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(resourceTaikunProjectRead, ctx, data, meta)
}

func resourceTaikunProjectDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
func flattenTaikunProject(projectDetailsDTO *models.ProjectDetailsForServersDto, boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto, projectQuotaDTO *models.ProjectQuotaListDto) map[string]interface{} {
	flavors := make([]string, len(boundFlavorDTOs))
	for i, boundFlavorDTO := range boundFlavorDTOs {
		flavors[i] = boundFlavorDTO.Name
	}

	projectMap := map[string]interface{}{
		"access_profile_id":     i32toa(projectDetailsDTO.AccessProfileID),
		"alerting_profile_name": projectDetailsDTO.AlertingProfileName,
		"cloud_credential_id":   i32toa(projectDetailsDTO.CloudID),
		"enable_auto_upgrade":   projectDetailsDTO.IsAutoUpgrade,
		"enable_monitoring":     projectDetailsDTO.IsMonitoringEnabled,
		"expiration_date":       rfc3339DateTimeToDate(projectDetailsDTO.ExpiredAt),
		"flavors":               flavors,
		"id":                    i32toa(projectDetailsDTO.ProjectID),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.KubernetesProfileID),
		"name":                  projectDetailsDTO.ProjectName,
		"organization_id":       i32toa(projectDetailsDTO.OrganizationID),
		"quota_id":              i32toa(projectDetailsDTO.QuotaID),
	}

	var nullID int32
	if projectDetailsDTO.AlertingProfileID != nullID {
		projectMap["alerting_profile_id"] = i32toa(projectDetailsDTO.AlertingProfileID)
	}

	if projectDetailsDTO.IsBackupEnabled {
		projectMap["backup_credential_id"] = i32toa(projectDetailsDTO.S3CredentialID)
	}

	if !projectQuotaDTO.IsCPUUnlimited {
		projectMap["quota_cpu_units"] = projectQuotaDTO.CPU
	}

	if !projectQuotaDTO.IsDiskSizeUnlimited {
		projectMap["quota_disk_size"] = projectQuotaDTO.DiskSize
	}

	if !projectQuotaDTO.IsRAMUnlimited {
		projectMap["quota_ram_size"] = projectQuotaDTO.RAM
	}

	return projectMap
}

func resourceTaikunProjectGetBoundFlavorDTOs(projectID int32, apiClient *apiClient) ([]*models.BoundFlavorsForProjectsListDto, error) {
	var boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto
	for {
		boundFlavorsParams := flavors.NewFlavorsGetSelectedFlavorsForProjectParams().WithV(ApiVersion).WithProjectID(&projectID)
		response, err := apiClient.client.Flavors.FlavorsGetSelectedFlavorsForProject(boundFlavorsParams, apiClient)
		if err != nil {
			return nil, err
		}
		boundFlavorDTOs = append(boundFlavorDTOs, response.Payload.Data...)
		boundFlavorDTOsCount := int32(len(boundFlavorDTOs))
		if boundFlavorDTOsCount == response.Payload.TotalCount {
			break
		}
		boundFlavorsParams = boundFlavorsParams.WithOffset(&boundFlavorDTOsCount)
	}
	return boundFlavorDTOs, nil
}
