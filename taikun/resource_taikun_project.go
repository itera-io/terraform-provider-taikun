package taikun

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

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
		"access_ip": {
			Description: "Public IP address of the bastion.",
			Type:        schema.TypeString,
			Computed:    true,
		},
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
		"auto_upgrade": {
			Description: "Kubespray version will be automatically upgraded if new version is available.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
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
		"lock": {
			Description: "Indicates whether to lock the project.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"monitoring": {
			Description: "Kubernetes cluster monitoring.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
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
			Description: "Maximum CPU units. Unlimited if unspecified.",
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
		"server_bastion": {
			Description:  "Bastion server.",
			Type:         schema.TypeSet,
			MaxItems:     1,
			Optional:     true,
			RequiredWith: []string{"server_kubemaster", "server_kubeworker"},
			Set:          hashAttributes("name", "disk_size", "flavor"),
			Elem: &schema.Resource{
				Schema: taikunServerBasicSchema(),
			},
		},
		"server_kubemaster": {
			Description:  "Kubemaster server.",
			Type:         schema.TypeSet,
			Optional:     true,
			RequiredWith: []string{"server_bastion", "server_kubeworker"},
			Set:          hashAttributes("name", "disk_size", "flavor", "kubernetes_node_label"),
			Elem: &schema.Resource{
				Schema: taikunServerSchemaWithKubernetesNodeLabels(),
			},
		},
		"server_kubeworker": {
			Description:  "Kubeworker server.",
			Type:         schema.TypeSet,
			Optional:     true,
			RequiredWith: []string{"server_bastion", "server_kubemaster"},
			Set:          hashAttributes("name", "disk_size", "flavor", "kubernetes_node_label"),
			Elem: &schema.Resource{
				Schema: taikunServerSchemaWithKubernetesNodeLabels(),
			},
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

func taikunServerSchemaWithKubernetesNodeLabels() map[string]*schema.Schema {
	serverSchema := taikunServerBasicSchema()
	serverSchema["kubernetes_node_label"] = &schema.Schema{
		Description: "Attach Kubernetes node labels.",
		Type:        schema.TypeList,
		Optional:    true,
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
			ValidateFunc: validation.IntAtLeast(30),
			Default:      30,
		},
		"flavor": {
			Description:  "The server's flavor.",
			Type:         schema.TypeString,
			Required:     true,
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

func resourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project",
		CreateContext: resourceTaikunProjectCreate,
		ReadContext:   generateResourceTaikunProjectRead(false),
		UpdateContext: resourceTaikunProjectUpdate,
		DeleteContext: resourceTaikunProjectDelete,
		Schema:        resourceTaikunProjectSchema(),
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue(
				"server_kubemaster",
				func(ctx context.Context, value, meta interface{}) error {
					set := value.(*schema.Set)

					if set.Len() != 0 && set.Len()%2 != 1 {
						return fmt.Errorf("there must be an odd number of server_kubemaster (currently %d)", set.Len())
					}
					return nil
				},
			),
			customdiff.ForceNewIfChange(
				"server_kubemaster",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return old.(*schema.Set).Len() != 0 && new.(*schema.Set).Len() != 0
				},
			),
			customdiff.ForceNewIfChange(
				"server_bastion",
				func(ctx context.Context, old, new, meta interface{}) bool {
					return old.(*schema.Set).Len() != 0 && new.(*schema.Set).Len() != 0
				},
			),
		),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
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
	if enableAutoUpgrade, enableAutoUpgradeIsSet := data.GetOk("auto_upgrade"); enableAutoUpgradeIsSet {
		body.IsAutoUpgrade = enableAutoUpgrade.(bool)
	}
	if enableMonitoring, enableMonitoringIsSet := data.GetOk("monitoring"); enableMonitoringIsSet {
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
	projectID, _ := atoi32(response.Payload.ID)

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

		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID) // TODO use /api/v1/projects endpoint?
		response, err := apiClient.client.Servers.ServersDetails(params, apiClient)

		if err == nil {
			quotaEditParams := project_quotas.NewProjectQuotasEditParams().WithV(ApiVersion).WithQuotaID(response.Payload.Project.QuotaID).WithBody(quotaEditBody)
			_, err := apiClient.client.ProjectQuotas.ProjectQuotasEdit(quotaEditParams, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	_, bastionsIsSet := data.GetOk("server_bastion")

	// Check if the project is not empty
	if bastionsIsSet {
		err = resourceTaikunProjectSetServers(data, apiClient, projectID)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectCommit(apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}
	}

	lock := data.Get("lock").(bool)
	if lock {
		lockMode := getLockMode(lock)
		params := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&projectID).WithMode(&lockMode)
		_, err := apiClient.client.Projects.ProjectsLockManager(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
		return diag.FromErr(err)
	}

	return readAfterCreateWithRetries(generateResourceTaikunProjectRead(true), ctx, data, meta)
}

func generateResourceTaikunProjectRead(isAfterUpdateOrCreate bool) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			if isAfterUpdateOrCreate {
				data.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
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
			if isAfterUpdateOrCreate {
				data.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		err = setResourceDataFromMap(data, flattenTaikunProject(projectDetailsDTO, response.Payload.Data, boundFlavorDTOs, quotaResponse.Payload.Data[0]))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(id)

		return nil
	}
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
	if data.HasChange("lock") {
		lock := data.Get("lock").(bool)
		lockMode := getLockMode(lock)
		params := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&id).WithMode(&lockMode)
		_, err := apiClient.client.Projects.ProjectsLockManager(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
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

	if data.HasChange("server_bastion") {
		o, n := data.GetChange("server_bastion")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		if oldSet.Len() == 0 {
			// The project was empty before
			if err := resourceTaikunProjectUpdateToggleBackupAndMonitoring(ctx, data, apiClient); err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectSetServers(data, apiClient, id); err != nil {
				return diag.FromErr(err)
			}

			if err := resourceTaikunProjectCommit(apiClient, id); err != nil {
				return diag.FromErr(err)
			}

		} else if newSet.Len() == 0 {
			// Purge
			oldKubeMasters, _ := data.GetChange("server_kubemaster")
			oldKubeWorkers, _ := data.GetChange("server_kubeworker")
			err = resourceTaikunProjectPurgeServers(
				o,
				oldKubeMasters,
				oldKubeWorkers,
				apiClient,
				id,
			)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectUpdateToggleBackupAndMonitoring(ctx, data, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if err := resourceTaikunProjectUpdateToggleBackupAndMonitoring(ctx, data, apiClient); err != nil {
			return diag.FromErr(err)
		}
		if data.HasChange("server_kubeworker") {
			o, n := data.GetChange("server_kubeworker")
			oldSet := o.(*schema.Set)
			newSet := n.(*schema.Set)
			toAdd := newSet.Difference(oldSet)
			toDel := oldSet.Difference(newSet)

			// Delete
			if toDel.Len() != 0 {
				serverIds := make([]int32, 0)

				for _, kubeWorker := range toDel.List() {
					kubeWorkerMap := kubeWorker.(map[string]interface{})
					kubeWorkerId, _ := atoi32(kubeWorkerMap["id"].(string))
					serverIds = append(serverIds, kubeWorkerId)
				}

				deleteServerBody := &models.DeleteServerCommand{
					ProjectID: id,
					ServerIds: serverIds,
				}
				deleteServerParams := servers.NewServersDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
				_, _, err := apiClient.client.Servers.ServersDelete(deleteServerParams, apiClient)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			// Create
			if toAdd.Len() != 0 {

				kubeWorkersList := oldSet.Intersection(newSet)

				for _, kubeWorker := range toAdd.List() {
					kubeWorkerMap := kubeWorker.(map[string]interface{})

					serverCreateBody := &models.ServerForCreateDto{
						Count:                1,
						DiskSize:             gibiByteToByte(kubeWorkerMap["disk_size"].(int)),
						Flavor:               kubeWorkerMap["flavor"].(string),
						KubernetesNodeLabels: resourceTaikunProjectServerKubernetesLabels(kubeWorkerMap),
						Name:                 kubeWorkerMap["name"].(string),
						ProjectID:            id,
						Role:                 300,
					}
					serverCreateParams := servers.NewServersCreateParams().WithV(ApiVersion).WithBody(serverCreateBody)
					serverCreateResponse, err := apiClient.client.Servers.ServersCreate(serverCreateParams, apiClient)
					if err != nil {
						return diag.FromErr(err)
					}
					kubeWorkerMap["id"] = serverCreateResponse.Payload.ID

					kubeWorkersList.Add(kubeWorkerMap)
				}

				err = data.Set("server_kubeworker", kubeWorkersList)
				if err != nil {
					return diag.FromErr(err)
				}
			}

			if err := resourceTaikunProjectCommit(apiClient, id); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunProjectRead(true), ctx, data, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readParams := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id) // TODO use /api/v1/projects endpoint?
	response, err := apiClient.client.Servers.ServersDetails(readParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.Project.IsLocked {
		unlockedMode := getLockMode(false)
		unlockParams := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&id).WithMode(&unlockedMode)
		if _, err := apiClient.client.Projects.ProjectsLockManager(unlockParams, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	// Purge all the servers
	err = resourceTaikunProjectPurgeServers(
		data.Get("server_bastion"),
		data.Get("server_kubemaster"),
		data.Get("server_kubeworker"),
		apiClient,
		id,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, id); err != nil {
		return diag.FromErr(err)
	}

	// Delete the project
	body := models.DeleteProjectCommand{ProjectID: id, IsForceDelete: false}
	params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.client.Projects.ProjectsDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceTaikunProjectUpdateToggleBackupAndMonitoring(ctx context.Context, data *schema.ResourceData, apiClient *apiClient) error {
	if err := resourceTaikunProjectUpdateToggleMonitoring(ctx, data, apiClient); err != nil {
		return err
	}
	if err := resourceTaikunProjectUpdateToggleBackup(ctx, data, apiClient); err != nil {
		return err
	}
	return nil
}

func resourceTaikunProjectUpdateToggleMonitoring(ctx context.Context, data *schema.ResourceData, apiClient *apiClient) error {
	if data.HasChange("monitoring") {
		projectID, _ := atoi32(data.Id())
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

func resourceTaikunProjectUpdateToggleBackup(ctx context.Context, data *schema.ResourceData, apiClient *apiClient) error {
	if data.HasChange("backup_credential_id") {
		projectID, _ := atoi32(data.Id())
		oldCredential, _ := data.GetChange("backup_credential_id")

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
					params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID) // TODO use /api/v1/projects endpoint?
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
				return errors.New(fmt.Sprintf("Error waiting for project (%s) to disable backup: %s", data.Id(), err))
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

// TODO change type of DTO if read endpoint is modified
func flattenTaikunProject(projectDetailsDTO *models.ProjectDetailsForServersDto, serverListDTO []*models.ServerListDto, boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto, projectQuotaDTO *models.ProjectQuotaListDto) map[string]interface{} {
	flavors := make([]string, len(boundFlavorDTOs))
	for i, boundFlavorDTO := range boundFlavorDTOs {
		flavors[i] = boundFlavorDTO.Name
	}

	projectMap := map[string]interface{}{
		"access_ip":             projectDetailsDTO.AccessIP,
		"access_profile_id":     i32toa(projectDetailsDTO.AccessProfileID),
		"alerting_profile_name": projectDetailsDTO.AlertingProfileName,
		"cloud_credential_id":   i32toa(projectDetailsDTO.CloudID),
		"auto_upgrade":          projectDetailsDTO.IsAutoUpgrade,
		"monitoring":            projectDetailsDTO.IsMonitoringEnabled,
		"expiration_date":       rfc3339DateTimeToDate(projectDetailsDTO.ExpiredAt),
		"flavors":               flavors,
		"id":                    i32toa(projectDetailsDTO.ProjectID),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.KubernetesProfileID),
		"lock":                  projectDetailsDTO.IsLocked,
		"name":                  projectDetailsDTO.ProjectName,
		"organization_id":       i32toa(projectDetailsDTO.OrganizationID),
		"quota_id":              i32toa(projectDetailsDTO.QuotaID),
	}

	bastions := make([]map[string]interface{}, 0)
	kubeMasters := make([]map[string]interface{}, 0)
	kubeWorkers := make([]map[string]interface{}, 0)
	for _, server := range serverListDTO {
		serverMap := map[string]interface{}{
			"created_by":       server.CreatedBy,
			"disk_size":        byteToGibiByte(server.DiskSize),
			"flavor":           server.OpenstackFlavor,
			"id":               i32toa(server.ID),
			"ip":               server.IPAddress,
			"last_modified":    server.LastModified,
			"last_modified_by": server.LastModifiedBy,
			"name":             server.Name,
			"status":           server.Status,
		}
		// Bastion
		if server.Role == "Bastion" {
			bastions = append(bastions, serverMap)
		} else {
			labels := make([]map[string]interface{}, len(server.KubernetesNodeLabels))
			for i, rawLabel := range server.KubernetesNodeLabels {
				labels[i] = map[string]interface{}{
					"key":   rawLabel.Key,
					"value": rawLabel.Value,
				}
			}
			serverMap["kubernetes_node_label"] = labels

			if server.Role == "Kubemaster" {
				kubeMasters = append(kubeMasters, serverMap)
			} else {
				kubeWorkers = append(kubeWorkers, serverMap)
			}
		}
	}
	projectMap["server_bastion"] = bastions
	projectMap["server_kubemaster"] = kubeMasters
	projectMap["server_kubeworker"] = kubeWorkers

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

func resourceTaikunProjectPurgeServers(bastions interface{}, kubeMasters interface{}, kubeWorkers interface{}, apiClient *apiClient, projectID int32) error {
	serverIds := make([]int32, 0)
	for _, bastion := range bastions.(*schema.Set).List() {
		bastionMap := bastion.(map[string]interface{})
		bastionId, _ := atoi32(bastionMap["id"].(string))
		serverIds = append(serverIds, bastionId)
	}

	for _, kubeMaster := range kubeMasters.(*schema.Set).List() {
		kubeMasterMap := kubeMaster.(map[string]interface{})
		kubeMasterId, _ := atoi32(kubeMasterMap["id"].(string))
		serverIds = append(serverIds, kubeMasterId)
	}

	for _, kubeWorker := range kubeWorkers.(*schema.Set).List() {
		kubeWorkerMap := kubeWorker.(map[string]interface{})
		kubeWorkerId, _ := atoi32(kubeWorkerMap["id"].(string))
		serverIds = append(serverIds, kubeWorkerId)
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

func resourceTaikunProjectSetServers(data *schema.ResourceData, apiClient *apiClient, projectID int32) error {

	bastions := data.Get("server_bastion")
	kubeMasters := data.Get("server_kubemaster")
	kubeWorkers := data.Get("server_kubeworker")

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
	err = data.Set("server_bastion", []map[string]interface{}{bastion})
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
	err = data.Set("server_kubemaster", kubeMastersList)
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
	err = data.Set("server_kubeworker", kubeWorkersList)
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

func resourceTaikunProjectWaitForStatus(ctx context.Context, targetList []string, pendingList []string, apiClient *apiClient, projectID int32) error {

	createStateConf := &resource.StateChangeConf{
		Pending: pendingList,
		Target:  targetList,
		Refresh: func() (interface{}, string, error) {
			params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
			resp, err := apiClient.client.Servers.ServersDetails(params, apiClient)
			if err != nil {
				return nil, "", err
			}

			return resp, resp.Payload.Project.ProjectStatus, nil
		},
		Timeout:                   40 * time.Minute,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for project (%d) to be in status %s: %s", projectID, targetList, err)
	}
	return nil
}

func resourceTaikunProjectServerKubernetesLabels(data map[string]interface{}) []*models.KubernetesNodeLabelsDto {
	labels, labelsAreSet := data["kubernetes_node_label"]
	if !labelsAreSet {
		return []*models.KubernetesNodeLabelsDto{}
	}
	labelsList := labels.([]interface{})
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
