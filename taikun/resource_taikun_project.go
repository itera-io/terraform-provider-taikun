package taikun

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/images"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/client/project_quotas"
	"github.com/itera-io/taikungoclient/client/stand_alone"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			Description:      "ID of the project's access profile. Defaults to the default access profile of the project's organization.",
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
			Description: "If enabled, the Kubespray version will be automatically upgraded when a new version is available.",
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
			Description:      "ID of the cloud credential used to create the project's servers.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"delete_on_expiration": {
			Description:  "If enabled, the project will be deleted on the expiration date and it will not be possible to recover it.",
			Type:         schema.TypeBool,
			Optional:     true,
			Default:      false,
			ForceNew:     true,
			RequiredWith: []string{"expiration_date"},
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
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"id": {
			Description: "Project ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"images": {
			Description: "List of images bound to the project.",
			Type:        schema.TypeSet,
			Optional:    true,
			DefaultFunc: func() (interface{}, error) {
				return []interface{}{}, nil
			},
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"kubernetes_profile_id": {
			Description:      "ID of the project's Kubernetes profile. Defaults to the default Kubernetes profile of the project's organization.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"kubernetes_version": {
			Description: "Kubernetes Version at project creation. Use the meta-argument `ignore_changes` to ignore future upgrades.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`),
				"Kubernets version must be in the format vMAJOR.MINOR.PATCH",
			),
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
		"policy_profile_id": {
			Description:      "ID of the Policy profile. If unspecified, Gatekeeper is disabled.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"quota_cpu_units": {
			Description:  "Maximum CPU units.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1000000,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_disk_size": {
			Description:  "Maximum disk size in GBs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      102400, // 100 TB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_ram_size": {
			Description:  "Maximum RAM size in GBs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      102400, // 100 TB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_cpu_units": {
			Description:  "Maximum CPU units for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1000000,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_volume_size": {
			Description:  "Maximum volume size in GBs for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      102400, // 100 TB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_ram_size": {
			Description:  "Maximum RAM size in GBs for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      102400, // 100 TB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"router_id_end_range": {
			Description:  "Router ID end range (specify only if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 255),
			RequiredWith: []string{"router_id_start_range", "taikun_lb_flavor"},
		},
		"router_id_start_range": {
			Description:  "Router ID start range (specify only if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
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
				Schema: taikunServerKubeworkerSchema(),
			},
		},
		"taikun_lb_flavor": {
			Description:  "OpenStack flavor for the Taikun load balancer (specify only if using OpenStack cloud credentials with Taikun Load Balancer enabled).",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"router_id_end_range", "router_id_start_range"},
		},
		"vm": {
			Description: "Virtual machines.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: taikunVMSchema(),
			},
		},
	}
}

func resourceTaikunProject() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Project",
		CreateContext: resourceTaikunProjectCreate,
		ReadContext:   generateResourceTaikunProjectReadWithoutRetries(),
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
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {

				names := make([]string, 0)

				for _, attribute := range []string{"server_bastion", "server_kubeworker", "server_kubemaster"} {
					if servers, serversIsSet := d.GetOk(attribute); serversIsSet {
						serversSet := servers.(*schema.Set)
						for _, server := range serversSet.List() {
							serverMap := server.(map[string]interface{})
							names = append(names, serverMap["name"].(string))
						}
					}
				}

				visitedMap := make(map[string]bool)
				for _, name := range names {
					if visitedMap[name] {
						return fmt.Errorf("server names must be unique: %s", name)
					}
					visitedMap[name] = true
				}

				return nil
			},
		),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(80 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	ctx, cancel := context.WithTimeout(ctx, 80*time.Minute)
	defer cancel()

	body := models.CreateProjectCommand{
		Name:         d.Get("name").(string),
		IsKubernetes: true,
	}
	body.CloudCredentialID, _ = atoi32(d.Get("cloud_credential_id").(string))
	flavorsData := d.Get("flavors").(*schema.Set).List()
	flavors := make([]string, len(flavorsData))
	for i, flavorData := range flavorsData {
		flavors[i] = flavorData.(string)
	}
	body.Flavors = flavors

	var projectOrganizationID int32 = -1

	if alertingProfileID, alertingProfileIDIsSet := d.GetOk("alerting_profile_id"); alertingProfileIDIsSet {
		body.AlertingProfileID, _ = atoi32(alertingProfileID.(string))
	}
	if backupCredentialID, backupCredentialIDIsSet := d.GetOk("backup_credential_id"); backupCredentialIDIsSet {
		body.IsBackupEnabled = true
		body.S3CredentialID, _ = atoi32(backupCredentialID.(string))
	}
	if enableAutoUpgrade, enableAutoUpgradeIsSet := d.GetOk("auto_upgrade"); enableAutoUpgradeIsSet {
		body.IsAutoUpgrade = enableAutoUpgrade.(bool)
	}
	if enableMonitoring, enableMonitoringIsSet := d.GetOk("monitoring"); enableMonitoringIsSet {
		body.IsMonitoringEnabled = enableMonitoring.(bool)
	}
	if deleteOnExpiration, deleteOnExpirationIsSet := d.GetOk("delete_on_expiration"); deleteOnExpirationIsSet {
		body.DeleteOnExpiration = deleteOnExpiration.(bool)
	}
	if expirationDate, expirationDateIsSet := d.GetOk("expiration_date"); expirationDateIsSet {
		dateTime := dateToDateTime(expirationDate.(string))
		body.ExpiredAt = &dateTime
	} else {
		body.ExpiredAt = nil
	}
	if PolicyProfileID, PolicyProfileIDIsSet := d.GetOk("policy_profile_id"); PolicyProfileIDIsSet {
		body.OpaProfileID, _ = atoi32(PolicyProfileID.(string))
	}

	if organizationID, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		projectOrganizationID, _ = atoi32(organizationID.(string))
		body.OrganizationID = projectOrganizationID
	}

	if taikunLBFlavor, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
		body.TaikunLBFlavor = taikunLBFlavor.(string)
		body.RouterIDStartRange = int32(d.Get("router_id_start_range").(int))
		body.RouterIDEndRange = int32(d.Get("router_id_end_range").(int))
	}
	if err := resourceTaikunProjectValidateKubernetesProfileLB(d, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if accessProfileID, accessProfileIDIsSet := d.GetOk("access_profile_id"); accessProfileIDIsSet {
		body.AccessProfileID, _ = atoi32(accessProfileID.(string))
	} else {
		if projectOrganizationID == -1 {
			if err := getDefaultOrganization(&projectOrganizationID, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
		defaultAccessProfileID, found, err := resourceTaikunProjectGetDefaultAccessProfile(projectOrganizationID, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if found {
			body.AccessProfileID = defaultAccessProfileID
		}
	}
	if kubernetesProfileID, kubernetesProfileIDIsSet := d.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		body.KubernetesProfileID, _ = atoi32(kubernetesProfileID.(string))
	} else {
		if projectOrganizationID == -1 {
			if err := getDefaultOrganization(&projectOrganizationID, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
		defaultKubernetesProfileID, found, err := resourceTaikunProjectGetDefaultKubernetesProfile(projectOrganizationID, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if found {
			body.KubernetesProfileID = defaultKubernetesProfileID
		}
	}

	if kubernetesVersion, kubernetesVersionIsSet := d.GetOk("kubernetes_version"); kubernetesVersionIsSet {
		body.KubernetesVersion = kubernetesVersion.(string)
	}

	// Send project creation request
	params := projects.NewProjectsCreateParams().WithV(ApiVersion).WithBody(&body).WithContext(ctx)
	response, err := apiClient.Client.Projects.ProjectsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Payload.ID)
	projectID, _ := atoi32(response.Payload.ID)

	if resourceTaikunProjectQuotaIsSet(d) {
		if err = resourceTaikunProjectEditQuotas(d, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, imagesIsSet := d.GetOk("images"); imagesIsSet {
		err := resourceTaikunProjectEditImages(d, apiClient, projectID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Check if the project is not empty
	if _, bastionsIsSet := d.GetOk("server_bastion"); bastionsIsSet {

		if err := resourceTaikunProjectSetServers(d, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectCommit(apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, vmIsSet := d.GetOk("vm"); vmIsSet {

		if err := resourceTaikunProjectSetVMs(d, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectStandaloneCommit(apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunProjectLock(projectID, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunProjectReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectRead(true)
}
func generateResourceTaikunProjectReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectRead(false)
}
func generateResourceTaikunProjectRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id := d.Id()
		id32, err := atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id32)
		response, err := apiClient.Client.Servers.ServersDetails(params, apiClient)
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		paramsVM := stand_alone.NewStandAloneDetailsParams().WithV(ApiVersion).WithProjectID(id32)
		responseVM, err := apiClient.Client.StandAlone.StandAloneDetails(paramsVM, apiClient)
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		projectDetailsDTO := response.Payload.Project
		serverList := response.Payload.Data
		vmList := responseVM.Payload.Data

		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(projectDetailsDTO.ProjectID, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		boundImageDTOs, err := resourceTaikunProjectGetBoundImageDTOs(projectDetailsDTO.ProjectID, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		quotaParams := project_quotas.NewProjectQuotasListParams().WithV(ApiVersion).WithID(&id32)
		quotaResponse, err := apiClient.Client.ProjectQuotas.ProjectQuotasList(quotaParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(quotaResponse.Payload.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		projectMap := flattenTaikunProject(projectDetailsDTO, serverList, vmList, boundFlavorDTOs, boundImageDTOs, quotaResponse.Payload.Data[0])
		usernames := resourceTaikunProjectGetResourceDataVmUsernames(d)
		if err := setResourceDataFromMap(d, projectMap); err != nil {
			return diag.FromErr(err)
		}

		if err := resourceTaikunProjectRestoreResourceDataVmUsernames(d, usernames); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)

		return nil
	}
}

func resourceTaikunProjectGetResourceDataVmUsernames(d *schema.ResourceData) (usernames map[string]string) {
	usernames = map[string]string{}

	vmListData, ok := d.GetOk("vm")
	if !ok {
		return
	}

	vmList, ok := vmListData.([]interface{})
	if !ok {
		return
	}
	for _, vmData := range vmList {
		vm, ok := vmData.(map[string]interface{})
		if !ok {
			return
		}

		vmIdData, ok := vm["id"]
		if !ok {
			continue
		}

		vmId, ok := vmIdData.(string)
		if !ok {
			continue
		}

		usernameData, ok := vm["username"]
		if !ok {
			continue
		}

		if username, ok := usernameData.(string); ok {
			usernames[vmId] = username
		}
	}

	return usernames
}

func resourceTaikunProjectRestoreResourceDataVmUsernames(d *schema.ResourceData, usernames map[string]string) error {
	if len(usernames) == 0 {
		return nil
	}

	vmListData, ok := d.GetOk("vm")
	if !ok {
		return nil
	}

	vmList, ok := vmListData.([]interface{})
	if !ok {
		return nil
	}

	for usernameVmId, username := range usernames {
		for _, vmData := range vmList {
			vm, ok := vmData.(map[string]interface{})
			if !ok {
				return nil
			}

			vmIdData, ok := vm["id"]
			if !ok {
				continue
			}

			vmId, ok := vmIdData.(string)
			if !ok {
				continue
			}

			if vmId == usernameVmId {
				vm["username"] = username
			}
		}
	}

	return d.Set("vm", vmList)
}

func resourceTaikunProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("alerting_profile_id") {
		body := models.AttachDetachAlertingProfileCommand{
			ProjectID: id,
		}
		detachParams := alerting_profiles.NewAlertingProfilesDetachParams().WithV(ApiVersion).WithBody(&body)
		if _, err := apiClient.Client.AlertingProfiles.AlertingProfilesDetach(detachParams, apiClient); err != nil {
			return diag.FromErr(err)
		}
		if newAlertingProfileIDData, newAlertingProfileIDProvided := d.GetOk("alerting_profile_id"); newAlertingProfileIDProvided {
			newAlertingProfileID, _ := atoi32(newAlertingProfileIDData.(string))
			body.AlertingProfileID = newAlertingProfileID
			attachParams := alerting_profiles.NewAlertingProfilesAttachParams().WithV(ApiVersion).WithBody(&body)
			if _, err := apiClient.Client.AlertingProfiles.AlertingProfilesAttach(attachParams, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("expiration_date") {
		body := models.ProjectExtendLifeTimeCommand{
			ProjectID: id,
		}
		if expirationDate, expirationDateIsSet := d.GetOk("expiration_date"); expirationDateIsSet {
			dateTime := dateToDateTime(expirationDate.(string))
			body.ExpireAt = &dateTime
		} else {
			body.ExpireAt = nil
		}
		params := projects.NewProjectsExtendLifeTimeParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.Client.Projects.ProjectsExtendLifeTime(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("flavors") {
		if err := resourceTaikunProjectEditFlavors(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("images") {
		if err := resourceTaikunProjectEditImages(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("quota_cpu_units", "quota_disk_size", "quota_ram_size", "quota_vm_cpu_units", "quota_vm_ram_size", "quota_vm_volume_size") {
		if err := resourceTaikunProjectEditQuotas(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("server_bastion") {
		oldBastions, newBastions := d.GetChange("server_bastion")
		oldSet := oldBastions.(*schema.Set)
		newSet := newBastions.(*schema.Set)

		if oldSet.Len() == 0 {
			// The project was empty before
			if err := resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectSetServers(d, apiClient, id); err != nil {
				return diag.FromErr(err)
			}

			if err := resourceTaikunProjectCommit(apiClient, id); err != nil {
				return diag.FromErr(err)
			}

			if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
				return diag.FromErr(err)
			}

		} else if newSet.Len() == 0 {
			// Purge
			oldKubeMasters, _ := d.GetChange("server_kubemaster")
			oldKubeWorkers, _ := d.GetChange("server_kubeworker")
			serversToPurge := resourceTaikunProjectFlattenServersData(oldBastions, oldKubeMasters, oldKubeWorkers)
			err = resourceTaikunProjectPurgeServers(serversToPurge, apiClient, id)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if err := resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("server_kubeworker") {
			o, n := d.GetChange("server_kubeworker")
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
				_, _, err := apiClient.Client.Servers.ServersDelete(deleteServerParams, apiClient)
				if err != nil {
					return diag.FromErr(err)
				}

				if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Deleting", "PendingDelete"}, apiClient, id); err != nil {
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
					serverCreateResponse, err := apiClient.Client.Servers.ServersCreate(serverCreateParams, apiClient)
					if err != nil {
						return diag.FromErr(err)
					}
					kubeWorkerMap["id"] = serverCreateResponse.Payload.ID

					kubeWorkersList.Add(kubeWorkerMap)
				}

				err = d.Set("server_kubeworker", kubeWorkersList)
				if err != nil {
					return diag.FromErr(err)
				}

				if err := resourceTaikunProjectCommit(apiClient, id); err != nil {
					return diag.FromErr(err)
				}

				if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("vm") {
		err = resourceTaikunProjectUpdateVMs(ctx, d, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunProjectLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, d, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	serversToPurge := resourceTaikunProjectFlattenServersData(
		d.Get("server_bastion"),
		d.Get("server_kubemaster"),
		d.Get("server_kubeworker"),
	)
	if len(serversToPurge) != 0 {
		err = resourceTaikunProjectPurgeServers(serversToPurge, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if vms := d.Get("vm").([]interface{}); len(vms) != 0 {
		err = resourceTaikunProjectPurgeVMs(vms, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}

	// Delete the project
	body := models.DeleteProjectCommand{ProjectID: id, IsForceDelete: false}
	params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.Client.Projects.ProjectsDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceTaikunProjectUnlockIfLocked(projectID int32, apiClient *taikungoclient.Client) error {
	readParams := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
	response, err := apiClient.Client.Servers.ServersDetails(readParams, apiClient)
	if err != nil {
		return err
	}

	if response.Payload.Project.IsLocked {
		if err := resourceTaikunProjectLock(projectID, false, apiClient); err != nil {
			return err
		}
	}

	return nil
}

func resourceTaikunProjectEditQuotas(d *schema.ResourceData, apiClient *taikungoclient.Client, projectID int32) (err error) {

	body := &models.UpdateQuotaCommand{
		QuotaID: projectID,
	}

	if cpu, ok := d.GetOk("quota_cpu_units"); ok {
		body.ServerCPU = int64(cpu.(int))
	}

	if ram, ok := d.GetOk("quota_ram_size"); ok {
		body.ServerRAM = gibiByteToByte(ram.(int))
	}

	if disk, ok := d.GetOk("quota_disk_size"); ok {
		body.ServerDiskSize = gibiByteToByte(disk.(int))
	}

	if vmCpu, ok := d.GetOk("quota_vm_cpu_units"); ok {
		body.VMCPU = int64(vmCpu.(int))
	}

	if vmRam, ok := d.GetOk("quota_vm_ram_size"); ok {
		body.VMRAM = gibiByteToByte(vmRam.(int))
	}

	if vmVolume, ok := d.GetOk("quota_vm_volume_size"); ok {
		body.VMVolumeSize = int64(vmVolume.(int)) // No conversion needed, API takes GBs
	}

	params := project_quotas.NewProjectQuotasEditParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.Client.ProjectQuotas.ProjectQuotasEdit(params, apiClient)
	return
}

func flattenTaikunProject(
	projectDetailsDTO *models.ProjectDetailsForServersDto,
	serverListDTO []*models.ServerListDto,
	vmListDTO []*models.StandaloneVmsListForDetailsDto,
	boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto,
	boundImageDTOs []*models.BoundImagesForProjectsListDto,
	projectQuotaDTO *models.ProjectQuotaListDto,
) map[string]interface{} {

	flavors := make([]string, len(boundFlavorDTOs))
	for i, boundFlavorDTO := range boundFlavorDTOs {
		flavors[i] = boundFlavorDTO.Name
	}

	images := make([]string, len(boundImageDTOs))
	for i, boundImageDTO := range boundImageDTOs {
		images[i] = boundImageDTO.ImageID
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
		"images":                images,
		"id":                    i32toa(projectDetailsDTO.ProjectID),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.KubernetesProfileID),
		"kubernetes_version":    projectDetailsDTO.KubernetesCurrentVersion,
		"lock":                  projectDetailsDTO.IsLocked,
		"name":                  projectDetailsDTO.ProjectName,
		"organization_id":       i32toa(projectDetailsDTO.OrganizationID),
		"quota_cpu_units":       projectQuotaDTO.ServerCPU,
		"quota_ram_size":        byteToGibiByte(projectQuotaDTO.ServerRAM),
		"quota_disk_size":       byteToGibiByte(projectQuotaDTO.ServerDiskSize),
		"quota_vm_cpu_units":    projectQuotaDTO.VMCPU,
		"quota_vm_ram_size":     byteToGibiByte(projectQuotaDTO.VMRAM),
		"quota_vm_volume_size":  projectQuotaDTO.VMVolumeSize,
	}

	bastions := make([]map[string]interface{}, 0)
	kubeMasters := make([]map[string]interface{}, 0)
	kubeWorkers := make([]map[string]interface{}, 0)
	for _, server := range serverListDTO {
		serverMap := map[string]interface{}{
			"created_by":       server.CreatedBy,
			"disk_size":        byteToGibiByte(server.DiskSize),
			"id":               i32toa(server.ID),
			"ip":               server.IPAddress,
			"last_modified":    server.LastModified,
			"last_modified_by": server.LastModifiedBy,
			"name":             server.Name,
			"status":           server.Status,
		}

		switch strings.ToLower(server.CloudType) {
		case "aws":
			serverMap["flavor"] = server.AwsInstanceType
		case "azure":
			serverMap["flavor"] = server.AzureVMSize
		case "openstack":
			serverMap["flavor"] = server.OpenstackFlavor
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

	vms := make([]map[string]interface{}, 0)
	for _, vm := range vmListDTO {
		vmMap := map[string]interface{}{
			"access_ip":             vm.PublicIP,
			"cloud_init":            vm.CloudInit,
			"created_by":            vm.CreatedBy,
			"flavor":                vm.TargetFlavor,
			"id":                    i32toa(vm.ID),
			"image_id":              vm.ImageID,
			"image_name":            vm.ImageName,
			"ip":                    vm.IPAddress,
			"last_modified":         vm.LastModified,
			"last_modified_by":      vm.LastModifiedBy,
			"name":                  vm.Name,
			"public_ip":             vm.PublicIPEnabled,
			"standalone_profile_id": i32toa(vm.Profile.ID),
			"status":                vm.Status,
			"volume_size":           vm.VolumeSize,
			"volume_type":           vm.VolumeType,
		}

		tags := make([]map[string]interface{}, len(vm.StandAloneMetaDatas))
		for i, rawTag := range vm.StandAloneMetaDatas {
			tags[i] = map[string]interface{}{
				"key":   rawTag.Key,
				"value": rawTag.Value,
			}
		}
		vmMap["tag"] = tags

		disks := make([]map[string]interface{}, len(vm.Disks))
		for i, rawDisk := range vm.Disks {
			lunId, _ := atoi32(rawDisk.LunID)
			disks[i] = map[string]interface{}{
				"device_name": rawDisk.DeviceName,
				"lun_id":      lunId,
				"id":          i32toa(rawDisk.ID),
				"name":        rawDisk.Name,
				"size":        rawDisk.CurrentSize,
				"volume_type": rawDisk.VolumeType,
			}
		}
		vmMap["disk"] = disks

		vms = append(vms, vmMap)
	}
	projectMap["vm"] = vms

	var nullID int32
	if projectDetailsDTO.AlertingProfileID != nullID {
		projectMap["alerting_profile_id"] = i32toa(projectDetailsDTO.AlertingProfileID)
	}

	if projectDetailsDTO.IsBackupEnabled {
		projectMap["backup_credential_id"] = i32toa(projectDetailsDTO.S3CredentialID)
	}

	if projectDetailsDTO.IsOpaEnabled {
		projectMap["policy_profile_id"] = i32toa(projectDetailsDTO.OpaProfileID)
	}

	return projectMap
}

func resourceTaikunProjectGetBoundFlavorDTOs(projectID int32, apiClient *taikungoclient.Client) ([]*models.BoundFlavorsForProjectsListDto, error) {
	var boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto
	boundFlavorsParams := flavors.NewFlavorsGetSelectedFlavorsForProjectParams().WithV(ApiVersion).WithProjectID(&projectID)
	for {
		response, err := apiClient.Client.Flavors.FlavorsGetSelectedFlavorsForProject(boundFlavorsParams, apiClient)
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

func resourceTaikunProjectGetBoundImageDTOs(projectID int32, apiClient *taikungoclient.Client) ([]*models.BoundImagesForProjectsListDto, error) {
	var boundImageDTOs []*models.BoundImagesForProjectsListDto
	boundImageParams := images.NewImagesGetSelectedImagesForProjectParams().WithV(ApiVersion).WithProjectID(&projectID)
	for {
		response, err := apiClient.Client.Images.ImagesGetSelectedImagesForProject(boundImageParams, apiClient)
		if err != nil {
			return nil, err
		}
		boundImageDTOs = append(boundImageDTOs, response.Payload.Data...)
		boundFlavorDTOsCount := int32(len(boundImageDTOs))
		if boundFlavorDTOsCount == response.Payload.TotalCount {
			break
		}
		boundImageParams = boundImageParams.WithOffset(&boundFlavorDTOsCount)
	}
	return boundImageDTOs, nil
}

func resourceTaikunProjectFlattenServersData(bastionsData interface{}, kubeMastersData interface{}, kubeWorkersData interface{}) []interface{} {
	var servers []interface{}
	servers = append(servers, bastionsData.(*schema.Set).List()...)
	servers = append(servers, kubeMastersData.(*schema.Set).List()...)
	servers = append(servers, kubeWorkersData.(*schema.Set).List()...)
	return servers
}

func resourceTaikunProjectWaitForStatus(ctx context.Context, targetList []string, pendingList []string, apiClient *taikungoclient.Client, projectID int32) error {
	createStateConf := &resource.StateChangeConf{
		Pending: pendingList,
		Target:  targetList,
		Refresh: func() (interface{}, string, error) {
			params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
			resp, err := apiClient.Client.Servers.ServersDetails(params, apiClient)
			if err != nil {
				return nil, "", err
			}

			return resp, resp.Payload.Project.ProjectStatus, nil
		},
		Timeout:                   80 * time.Minute,
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

func resourceTaikunProjectValidateKubernetesProfileLB(d *schema.ResourceData, apiClient *taikungoclient.Client) error {
	if kubernetesProfileIDData, kubernetesProfileIsSet := d.GetOk("kubernetes_profile_id"); kubernetesProfileIsSet {
		kubernetesProfileID, _ := atoi32(kubernetesProfileIDData.(string))
		lbSolution, err := resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID, apiClient)
		if err != nil {
			return err
		}
		if lbSolution == loadBalancerTaikun {
			cloudCredentialID, _ := atoi32(d.Get("cloud_credential_id").(string))
			cloudType, err := resourceTaikunProjectGetCloudType(cloudCredentialID, apiClient)
			if err != nil {
				return err
			}
			if _, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); !taikunLBFlavorIsSet {
				return fmt.Errorf("If Taikun load balancer is enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor must be set")
			}
			if cloudType != cloudTypeOpenStack {
				return fmt.Errorf("If Taikun load balancer is enabled, cloud type should be OpenStack; is %s", cloudType)
			}
		} else if _, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
			return fmt.Errorf("If Taikun load balancer is not enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor should not be set")
		}
	}
	return nil
}

func resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID int32, apiClient *taikungoclient.Client) (string, error) {
	params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&kubernetesProfileID)
	response, err := apiClient.Client.KubernetesProfiles.KubernetesProfilesList(params, apiClient)
	if err != nil {
		return "", err
	}
	if len(response.Payload.Data) == 0 {
		return "", fmt.Errorf("kubernetes profile with ID %d not found", kubernetesProfileID)
	}
	kubernetesProfile := response.Payload.Data[0]
	return getLoadBalancingSolution(kubernetesProfile.OctaviaEnabled, kubernetesProfile.TaikunLBEnabled), nil
}

func resourceTaikunProjectGetCloudType(cloudCredentialID int32, apiClient *taikungoclient.Client) (string, error) {
	params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&cloudCredentialID)
	response, err := apiClient.Client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
	if err != nil {
		return "", err
	}
	if len(response.Payload.Amazon) == 1 {
		return cloudTypeAWS, nil
	}
	if len(response.Payload.Azure) == 1 {
		return cloudTypeAzure, nil
	}
	if len(response.Payload.Openstack) == 1 {
		return cloudTypeOpenStack, nil
	}
	return "", fmt.Errorf("cloud credential with ID %d not found", cloudCredentialID)
}

const defaultAccessProfileName = "default"

func resourceTaikunProjectGetDefaultAccessProfile(organizationID int32, apiClient *taikungoclient.Client) (accessProfileID int32, found bool, err error) {
	params := access_profiles.NewAccessProfilesAccessProfilesForOrganizationListParams().WithV(ApiVersion).WithOrganizationID(&organizationID)
	response, err := apiClient.Client.AccessProfiles.AccessProfilesAccessProfilesForOrganizationList(params, apiClient)
	if err != nil {
		return 0, false, err
	}
	for _, profile := range response.Payload {
		if profile.Name == defaultAccessProfileName {
			return profile.ID, true, nil
		}
	}
	return 0, false, nil
}

const defaultKubernetesProfileName = "default"

func resourceTaikunProjectGetDefaultKubernetesProfile(organizationID int32, apiClient *taikungoclient.Client) (kubernetesProfileID int32, found bool, err error) {
	params := kubernetes_profiles.NewKubernetesProfilesKubernetesProfilesForOrganizationListParams().WithV(ApiVersion).WithOrganizationID(&organizationID)
	response, err := apiClient.Client.KubernetesProfiles.KubernetesProfilesKubernetesProfilesForOrganizationList(params, apiClient)
	if err != nil {
		return 0, false, err
	}
	for _, profile := range response.Payload {
		if profile.Name == defaultKubernetesProfileName {
			return profile.ID, true, nil
		}
	}
	return 0, false, nil
}

func resourceTaikunProjectLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	lockMode := getLockMode(lock)
	params := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&id).WithMode(&lockMode)
	_, err := apiClient.Client.Projects.ProjectsLockManager(params, apiClient)
	return err
}

func resourceTaikunProjectQuotaIsSet(d *schema.ResourceData) bool {
	if _, ok := d.GetOk("quota_cpu_units"); ok {
		return true
	}

	if _, ok := d.GetOk("quota_disk_size"); ok {
		return true
	}

	if _, ok := d.GetOk("quota_ram_size"); ok {
		return true
	}

	if _, ok := d.GetOk("quota_vm_cpu_units"); ok {
		return true
	}

	if _, ok := d.GetOk("quota_vm_volume_size"); ok {
		return true
	}

	if _, ok := d.GetOk("quota_vm_ram_size"); ok {
		return true
	}

	return false
}
