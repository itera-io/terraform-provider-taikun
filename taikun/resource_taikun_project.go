package taikun

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			Computed:         false,
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
			Description: "If enabled, the project will be deleted on the expiration date and it will not be possible to recover it.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			//ForceNew:     true, // We do not need to force recreate project for just delete on expiration update.
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
	apiClient := meta.(*tk.Client)
	ctx, cancel := context.WithTimeout(ctx, 80*time.Minute)
	defer cancel()

	body := tkcore.CreateProjectCommand{}
	body.SetName(d.Get("name").(string))
	body.SetIsKubernetes(true)

	cloudCredentialID, _ := atoi32(d.Get("cloud_credential_id").(string))
	body.SetCloudCredentialId(cloudCredentialID)
	flavorsData := d.Get("flavors").(*schema.Set).List()
	flavors := make([]string, len(flavorsData))
	for i, flavorData := range flavorsData {
		flavors[i] = flavorData.(string)
	}
	body.SetFlavors(flavors)

	if alertingProfileID, alertingProfileIDIsSet := d.GetOk("alerting_profile_id"); alertingProfileIDIsSet {
		alertingId, _ := atoi32(alertingProfileID.(string))
		body.SetAlertingProfileId(alertingId)
	}
	if backupCredentialID, backupCredentialIDIsSet := d.GetOk("backup_credential_id"); backupCredentialIDIsSet {
		body.SetIsBackupEnabled(true)
		backupCredential, _ := atoi32(backupCredentialID.(string))
		body.SetS3CredentialId(backupCredential)
	}
	if enableAutoUpgrade, enableAutoUpgradeIsSet := d.GetOk("auto_upgrade"); enableAutoUpgradeIsSet {
		body.SetIsAutoUpgrade(enableAutoUpgrade.(bool))
	}
	if enableMonitoring, enableMonitoringIsSet := d.GetOk("monitoring"); enableMonitoringIsSet {
		body.SetIsMonitoringEnabled(enableMonitoring.(bool))
	}
	if deleteOnExpiration, deleteOnExpirationIsSet := d.GetOk("delete_on_expiration"); deleteOnExpirationIsSet {
		body.SetDeleteOnExpiration(deleteOnExpiration.(bool))
	}
	if expirationDate, expirationDateIsSet := d.GetOk("expiration_date"); expirationDateIsSet {
		dateTime := dateToDateTime(expirationDate.(string))
		body.SetExpiredAt(time.Time(dateTime))
	} else {
		body.SetExpiredAtNil()
	}
	if PolicyProfileID, PolicyProfileIDIsSet := d.GetOk("policy_profile_id"); PolicyProfileIDIsSet {
		opaProfileID, _ := atoi32(PolicyProfileID.(string))
		body.SetOpaProfileId(opaProfileID)
	}

	if organizationID, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		projectOrganizationID, _ := atoi32(organizationID.(string))
		body.SetOrganizationId(projectOrganizationID)
	}

	if taikunLBFlavor, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
		body.SetTaikunLBFlavor(taikunLBFlavor.(string))
		body.SetRouterIdStartRange(int32(d.Get("router_id_start_range").(int)))
		body.SetRouterIdEndRange(int32(d.Get("router_id_end_range").(int)))
	}
	if err := resourceTaikunProjectValidateKubernetesProfileLB(d, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if accessProfileID, accessProfileIDIsSet := d.GetOk("access_profile_id"); accessProfileIDIsSet {
		accessProfile, _ := atoi32(accessProfileID.(string))
		body.SetAccessProfileId(accessProfile)
	}

	if kubernetesProfileID, kubernetesProfileIDIsSet := d.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		kubeId, _ := atoi32(kubernetesProfileID.(string))
		body.SetKubernetesProfileId(kubeId)
	}

	if kubernetesVersion, kubernetesVersionIsSet := d.GetOk("kubernetes_version"); kubernetesVersionIsSet {
		body.SetKubernetesVersion(kubernetesVersion.(string))
	}

	// Send project creation request
	response, responseBody, err := apiClient.Client.ProjectsAPI.ProjectsCreate(context.TODO()).CreateProjectCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(responseBody, err))
	}

	d.SetId(response.GetId())
	projectID, _ := atoi32(response.GetId())

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
		apiClient := meta.(*tk.Client)
		id := d.Id()
		id32, err := atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, _, err := apiClient.Client.ServersAPI.ServersDetails(ctx, id32).Execute()
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		responseVM, _, err := apiClient.Client.StandaloneAPI.StandaloneDetails(ctx, id32).Execute()
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		projectDetailsDTO := response.GetProject()
		serverList := response.Data
		vmList := responseVM.Data

		boundFlavorDTOs, err := resourceTaikunProjectGetBoundFlavorDTOs(projectDetailsDTO.GetProjectId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		var boundImageDTOs []tkcore.BoundImagesForProjectsListDto

		boundImageDTOs, err = resourceTaikunProjectGetBoundImageDTOs(projectDetailsDTO.GetProjectId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		quotaResponse, bodyResponse, err := apiClient.Client.ProjectQuotasAPI.ProjectquotasList(context.TODO()).Id(id32).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(bodyResponse, err))
		}
		if len(quotaResponse.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		deleteOnExpiration, err := resourceTaikunProjectGetDeleteOnExpiration(projectDetailsDTO.GetProjectId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		projectMap := flattenTaikunProject(&projectDetailsDTO, serverList, vmList, boundFlavorDTOs, boundImageDTOs, &quotaResponse.Data[0], deleteOnExpiration)
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
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("alerting_profile_id") {
		body := tkcore.AttachDetachAlertingProfileCommand{}
		body.SetProjectId(id)
		bodyResponse, newErr := apiClient.Client.AlertingProfilesAPI.AlertingprofilesDetach(context.TODO()).AttachDetachAlertingProfileCommand(body).Execute()
		if newErr != nil {
			return diag.FromErr(tk.CreateError(bodyResponse, newErr))
		}
		if newAlertingProfileIDData, newAlertingProfileIDProvided := d.GetOk("alerting_profile_id"); newAlertingProfileIDProvided {
			newAlertingProfileID, _ := atoi32(newAlertingProfileIDData.(string))
			body.SetAlertingProfileId(newAlertingProfileID)
			bodyResponse, newErr := apiClient.Client.AlertingProfilesAPI.AlertingprofilesAttach(context.TODO()).AttachDetachAlertingProfileCommand(body).Execute()
			if newErr != nil {
				return diag.FromErr(tk.CreateError(bodyResponse, newErr))
			}
		}
	}

	// expiration_date can exist without delete_on_expiration
	// delete_on_expiration must exist with expiration_date
	if d.HasChange("expiration_date") || d.HasChange("delete_on_expiration") {
		body := tkcore.ProjectExtendLifeTimeCommand{}
		body.SetProjectId(id)

		if expirationDate, expirationDateIsSet := d.GetOk("expiration_date"); expirationDateIsSet {
			dateTime := dateToDateTime(expirationDate.(string))
			body.SetExpireAt(time.Time(dateTime))
		} else {
			body.SetExpireAtNil()
		}

		if deleteOnExpiration, deleteOnExpirationIsSet := d.GetOk("delete_on_expiration"); deleteOnExpirationIsSet {
			body.SetDeleteOnExpiration(deleteOnExpiration.(bool))
		} else {
			body.SetDeleteOnExpiration(false)
		}

		_, err = apiClient.Client.ProjectsAPI.ProjectsExtendLifetime(context.TODO()).ProjectExtendLifeTimeCommand(body).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("flavors") {
		if err = resourceTaikunProjectEditFlavors(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("images") {
		if err = resourceTaikunProjectEditImages(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("quota_cpu_units", "quota_disk_size", "quota_ram_size", "quota_vm_cpu_units", "quota_vm_ram_size", "quota_vm_volume_size") {
		if err = resourceTaikunProjectEditQuotas(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("server_bastion") {
		oldBastions, newBastions := d.GetChange("server_bastion")
		oldSet := oldBastions.(*schema.Set)
		newSet := newBastions.(*schema.Set)

		if oldSet.Len() == 0 {
			// The project was empty before
			if err = resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
			if err = resourceTaikunProjectSetServers(d, apiClient, id); err != nil {
				return diag.FromErr(err)
			}

			if err = resourceTaikunProjectCommit(apiClient, id); err != nil {
				return diag.FromErr(err)
			}

			if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
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
			if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
				return diag.FromErr(err)
			}
			if err = resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if err = resourceTaikunProjectUpdateToggleServices(ctx, d, apiClient); err != nil {
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

				deleteServerBody := tkcore.DeleteServerCommand{}
				deleteServerBody.SetProjectId(id)
				deleteServerBody.SetServerIds(serverIds)

				_, err = apiClient.Client.ServersAPI.ServersDelete(ctx).DeleteServerCommand(deleteServerBody).Execute()
				if err != nil {
					return diag.FromErr(err)
				}

				if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Deleting", "PendingDelete"}, apiClient, id); err != nil {
					return diag.FromErr(err)
				}
			}
			// Create
			if toAdd.Len() != 0 {

				kubeWorkersList := oldSet.Intersection(newSet)

				for _, kubeWorker := range toAdd.List() {
					kubeWorkerMap := kubeWorker.(map[string]interface{})

					serverCreateBody := tkcore.ServerForCreateDto{}
					serverCreateBody.SetCount(1)
					serverCreateBody.SetDiskSize(gibiByteToByte(kubeWorkerMap["disk_size"].(int)))
					serverCreateBody.SetFlavor(kubeWorkerMap["flavor"].(string))
					serverCreateBody.SetKubernetesNodeLabels(resourceTaikunProjectServerKubernetesLabels(kubeWorkerMap))
					serverCreateBody.SetName(kubeWorkerMap["name"].(string))
					serverCreateBody.SetProjectId(id)
					serverCreateBody.SetRole(tkcore.CLOUDROLE_KUBEWORKER)

					serverCreateResponse, _, newErr := apiClient.Client.ServersAPI.ServersCreate(ctx).ServerForCreateDto(serverCreateBody).Execute()
					if newErr != nil {
						return diag.FromErr(newErr)
					}
					kubeWorkerMap["id"] = serverCreateResponse.GetId()

					kubeWorkersList.Add(kubeWorkerMap)
				}

				err = d.Set("server_kubeworker", kubeWorkersList)
				if err != nil {
					return diag.FromErr(err)
				}

				if err = resourceTaikunProjectCommit(apiClient, id); err != nil {
					return diag.FromErr(err)
				}

				if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
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
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
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
		if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}
	if vms := d.Get("vm").([]interface{}); len(vms) != 0 {
		err = resourceTaikunProjectPurgeVMs(vms, apiClient, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}

	// Delete the project
	body := tkcore.DeleteProjectCommand{}
	body.SetProjectId(id)
	body.SetIsForceDelete(false)
	_, err = apiClient.Client.ProjectsAPI.ProjectsDelete(ctx).DeleteProjectCommand(body).Execute()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceTaikunProjectUnlockIfLocked(projectID int32, apiClient *tk.Client) error {
	response, _, err := apiClient.Client.ServersAPI.ServersDetails(context.TODO(), projectID).Execute()
	if err != nil {
		return err
	}

	project := response.GetProject()
	if project.GetIsLocked() {
		if err := resourceTaikunProjectLock(projectID, false, apiClient); err != nil {
			return err
		}
	}

	return nil
}

func resourceTaikunProjectEditQuotas(d *schema.ResourceData, apiClient *tk.Client, projectID int32) (err error) {

	body := tkcore.UpdateQuotaCommand{}
	body.SetQuotaId(projectID)

	if cpu, ok := d.GetOk("quota_cpu_units"); ok {
		body.SetServerCpu(int64(cpu.(int)))
	}

	if ram, ok := d.GetOk("quota_ram_size"); ok {
		body.SetServerRam(gibiByteToByte(ram.(int)))
	}

	if disk, ok := d.GetOk("quota_disk_size"); ok {
		body.SetServerDiskSize(gibiByteToByte(disk.(int)))
	}

	if vmCpu, ok := d.GetOk("quota_vm_cpu_units"); ok {
		body.SetVmCpu(int64(vmCpu.(int)))
	}

	if vmRam, ok := d.GetOk("quota_vm_ram_size"); ok {
		body.SetVmRam(gibiByteToByte(vmRam.(int)))
	}

	if vmVolume, ok := d.GetOk("quota_vm_volume_size"); ok {
		body.SetVmVolumeSize(int64(vmVolume.(int))) // No conversion needed, API takes GBs
	}

	_, err = apiClient.Client.ProjectQuotasAPI.ProjectquotasUpdate(context.TODO()).UpdateQuotaCommand(body).Execute()
	return
}

func flattenTaikunProject(
	projectDetailsDTO *tkcore.ProjectDetailsForServersDto,
	serverListDTO []tkcore.ServerListDto,
	vmListDTO []tkcore.StandaloneVmsListForDetailsDto,
	boundFlavorDTOs []tkcore.BoundFlavorsForProjectsListDto,
	boundImageDTOs []tkcore.BoundImagesForProjectsListDto,
	projectQuotaDTO *tkcore.ProjectQuotaListDto,
	projectDeleteOnExpiration bool,
) map[string]interface{} {

	flavors := make([]string, len(boundFlavorDTOs))
	for i, boundFlavorDTO := range boundFlavorDTOs {
		flavors[i] = boundFlavorDTO.GetName()
	}

	images := make([]string, len(boundImageDTOs))
	for i, boundImageDTO := range boundImageDTOs {
		images[i] = boundImageDTO.GetImageId()
	}

	projectMap := map[string]interface{}{
		"access_ip":             projectDetailsDTO.GetAccessIp(),
		"access_profile_id":     i32toa(projectDetailsDTO.GetAccessProfileId()),
		"alerting_profile_name": projectDetailsDTO.GetAlertingProfileName(),
		"cloud_credential_id":   i32toa(projectDetailsDTO.GetCloudId()),
		"auto_upgrade":          projectDetailsDTO.GetIsAutoUpgrade(),
		"monitoring":            projectDetailsDTO.GetIsMonitoringEnabled(),
		"delete_on_expiration":  projectDeleteOnExpiration,
		"expiration_date":       rfc3339DateTimeToDate(projectDetailsDTO.GetExpiredAt()),
		"flavors":               flavors,
		"images":                images,
		"id":                    i32toa(projectDetailsDTO.GetProjectId()),
		"kubernetes_profile_id": i32toa(projectDetailsDTO.GetKubernetesProfileId()),
		"kubernetes_version":    projectDetailsDTO.GetKubernetesCurrentVersion(),
		"lock":                  projectDetailsDTO.GetIsLocked(),
		"name":                  projectDetailsDTO.GetProjectName(),
		"organization_id":       i32toa(projectDetailsDTO.GetOrganizationId()),
		"quota_cpu_units":       projectQuotaDTO.GetServerCpu(),
		"quota_ram_size":        byteToGibiByte(projectQuotaDTO.GetServerRam()),
		"quota_disk_size":       byteToGibiByte(projectQuotaDTO.GetServerDiskSize()),
		"quota_vm_cpu_units":    projectQuotaDTO.GetVmCpu(),
		"quota_vm_ram_size":     byteToGibiByte(projectQuotaDTO.GetVmRam()),
		"quota_vm_volume_size":  projectQuotaDTO.GetVmVolumeSize(),
	}

	bastions := make([]map[string]interface{}, 0)
	kubeMasters := make([]map[string]interface{}, 0)
	kubeWorkers := make([]map[string]interface{}, 0)
	for _, server := range serverListDTO {
		serverMap := map[string]interface{}{
			"created_by":       server.GetCreatedBy(),
			"disk_size":        byteToGibiByte(server.GetDiskSize()),
			"id":               i32toa(server.GetId()),
			"ip":               server.GetIpAddress(),
			"last_modified":    server.GetLastModified(),
			"last_modified_by": server.GetLastModifiedBy(),
			"name":             server.GetName(),
			"status":           server.GetStatus(),
		}

		switch server.GetCloudType() {
		case tkcore.CLOUDTYPE_AWS:
			serverMap["flavor"] = server.GetAwsInstanceType()
		case tkcore.CLOUDTYPE_AZURE:
			serverMap["flavor"] = server.GetAzureVmSize()
		case tkcore.CLOUDTYPE_OPENSTACK:
			serverMap["flavor"] = server.GetOpenstackFlavor()
		case tkcore.CLOUDTYPE_GOOGLE, "google":
			serverMap["flavor"] = server.GetGoogleMachineType()
		}

		// Bastion
		if server.GetRole() == tkcore.CLOUDROLE_BASTION {
			bastions = append(bastions, serverMap)
		} else {

			labels := make([]map[string]interface{}, len(server.GetKubernetesNodeLabels()))
			for i, rawLabel := range server.GetKubernetesNodeLabels() {
				labels[i] = map[string]interface{}{
					"key":   *rawLabel.Key.Get(),
					"value": *rawLabel.Value.Get(),
				}
			}
			serverMap["kubernetes_node_label"] = labels

			if server.GetRole() == tkcore.CLOUDROLE_KUBEMASTER {
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
			"access_ip":             vm.GetPublicIp(),
			"cloud_init":            vm.GetCloudInit(),
			"created_by":            vm.GetCreatedBy(),
			"flavor":                vm.GetTargetFlavor(),
			"id":                    i32toa(vm.GetId()),
			"image_id":              vm.GetImageId(),
			"image_name":            vm.GetImageName(),
			"ip":                    vm.GetIpAddress(),
			"last_modified":         vm.GetLastModified(),
			"last_modified_by":      vm.GetLastModifiedBy(),
			"name":                  vm.GetName(),
			"public_ip":             vm.GetPublicIpEnabled(),
			"standalone_profile_id": i32toa(vm.Profile.GetId()),
			"status":                vm.GetStatus(),
			"volume_size":           vm.GetVolumeSize(),
			"volume_type":           vm.GetVolumeType(),
		}

		tags := make([]map[string]interface{}, len(vm.GetStandAloneMetaDatas()))
		for i, rawTag := range vm.GetStandAloneMetaDatas() {
			tags[i] = map[string]interface{}{
				"key":   rawTag.GetKey(),
				"value": rawTag.GetValue(),
			}
		}
		vmMap["tag"] = tags

		disks := make([]map[string]interface{}, len(vm.GetDisks()))
		for i, rawDisk := range vm.GetDisks() {
			disks[i] = map[string]interface{}{
				"device_name": rawDisk.GetDeviceName(),
				"id":          i32toa(rawDisk.GetId()),
				"name":        rawDisk.GetName(),
				"size":        rawDisk.GetCurrentSize(),
				"volume_type": rawDisk.GetVolumeType(),
			}
		}
		vmMap["disk"] = disks

		vms = append(vms, vmMap)
	}
	projectMap["vm"] = vms

	var nullID int32
	if projectDetailsDTO.GetAlertingProfileId() != nullID {
		projectMap["alerting_profile_id"] = i32toa(projectDetailsDTO.GetAlertingProfileId())
	}

	if projectDetailsDTO.GetIsBackupEnabled() {
		projectMap["backup_credential_id"] = i32toa(projectDetailsDTO.GetS3CredentialId())
	}

	if projectDetailsDTO.GetIsOpaEnabled() {
		projectMap["policy_profile_id"] = i32toa(projectDetailsDTO.GetOpaProfileId())
	}

	return projectMap
}

func resourceTaikunProjectGetDeleteOnExpiration(projectID int32, apiClient *tk.Client) (bool, error) {
	data, response, err := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(projectID).Execute()
	if err != nil {
		return false, tk.CreateError(response, err)
	}
	return data.GetData()[0].GetDeleteOnExpiration(), nil
}

func resourceTaikunProjectGetBoundFlavorDTOs(projectID int32, apiClient *tk.Client) ([]tkcore.BoundFlavorsForProjectsListDto, error) {
	var boundFlavorDTOs []tkcore.BoundFlavorsForProjectsListDto
	var offset int32 = 0
	params := apiClient.Client.FlavorsAPI.FlavorsSelectedFlavorsForProject(context.TODO()).ProjectId(projectID)
	for {
		response, _, err := params.Offset(offset).Execute()
		if err != nil {
			return nil, err
		}
		boundFlavorDTOs = append(boundFlavorDTOs, response.Data...)
		boundFlavorDTOsCount := int32(len(boundFlavorDTOs))
		if boundFlavorDTOsCount == response.GetTotalCount() {
			break
		}
		offset = boundFlavorDTOsCount
	}
	return boundFlavorDTOs, nil
}

func resourceTaikunProjectGetBoundImageDTOs(projectID int32, apiClient *tk.Client) ([]tkcore.BoundImagesForProjectsListDto, error) {
	var boundImageDTOs []tkcore.BoundImagesForProjectsListDto
	var offset int32 = 0
	for {
		response, _, err := apiClient.Client.ImagesAPI.ImagesSelectedImagesForProject(context.TODO()).ProjectId(projectID).Offset(offset).Execute()
		if err != nil {
			return nil, err
		}
		boundImageDTOs = append(boundImageDTOs, response.GetData()...)
		boundFlavorDTOsCount := int32(len(boundImageDTOs))
		if boundFlavorDTOsCount == response.GetTotalCount() {
			break
		}
		offset = boundFlavorDTOsCount
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

func resourceTaikunProjectWaitForStatus(ctx context.Context, targetList []string, pendingList []string, apiClient *tk.Client, projectID int32) error {
	createStateConf := &resource.StateChangeConf{
		Pending: pendingList,
		Target:  targetList,
		Refresh: func() (interface{}, string, error) {
			resp, _, err := apiClient.Client.ServersAPI.ServersDetails(context.TODO(), projectID).Execute()
			if err != nil {
				return nil, "", err
			}

			project := resp.GetProject()
			status := project.GetProjectStatus()

			return resp, string(status), nil
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

func resourceTaikunProjectValidateKubernetesProfileLB(d *schema.ResourceData, apiClient *tk.Client) error {
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
				return fmt.Errorf("if Taikun load balancer is enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor must be set")
			}
			if cloudType != cloudTypeOpenStack {
				return fmt.Errorf("if Taikun load balancer is enabled, cloud type should be OpenStack; is %s", cloudType)
			}
		} else if _, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
			return fmt.Errorf("if Taikun load balancer is not enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor should not be set")
		}
	}
	return nil
}

func resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID int32, apiClient *tk.Client) (string, error) {
	response, _, err := apiClient.Client.KubernetesProfilesAPI.KubernetesprofilesList(context.TODO()).Id(kubernetesProfileID).Execute()
	if err != nil {
		return "", err
	}
	if len(response.GetData()) == 0 {
		return "", fmt.Errorf("kubernetes profile with ID %d not found", kubernetesProfileID)
	}
	kubernetesProfile := response.GetData()[0]
	return getLoadBalancingSolution(kubernetesProfile.GetOctaviaEnabled(), kubernetesProfile.GetTaikunLBEnabled()), nil
}

func resourceTaikunProjectGetCloudType(cloudCredentialID int32, apiClient *tk.Client) (string, error) {
	response, _, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	}
	if len(response.GetAmazon()) == 1 {
		return string(tkcore.CLOUDTYPE_AWS), nil
	}
	if len(response.GetAzure()) == 1 {
		return string(tkcore.CLOUDTYPE_AZURE), nil
	}
	if len(response.GetOpenstack()) == 1 {
		return string(tkcore.CLOUDTYPE_OPENSTACK), nil
	}
	if len(response.GetGoogle()) == 1 {
		return string(tkcore.CLOUDTYPE_GOOGLE), nil
	}
	return "", fmt.Errorf("cloud credential with ID %d not found", cloudCredentialID)
}

func resourceTaikunProjectLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.ProjectLockManagerCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))
	_, err := apiClient.Client.ProjectsAPI.ProjectsLockManager(context.TODO()).ProjectLockManagerCommand(body).Execute()
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
