package project

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
			ValidateDiagFunc: utils.StringIsInt,
			ForceNew:         true,
		},
		"alerting_profile_id": {
			Description:      "ID of the project's alerting profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         false,
			ValidateDiagFunc: utils.StringIsInt,
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
			ValidateDiagFunc: utils.StringIsInt,
		},
		"cloud_credential_id": {
			Description:      "ID of the cloud credential used to create the project's servers.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: utils.StringIsInt,
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
			ValidateDiagFunc: utils.StringIsDate,
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
			ValidateDiagFunc: utils.StringIsInt,
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
			ValidateDiagFunc: utils.StringIsInt,
			ForceNew:         true,
		},
		"policy_profile_id": {
			Description:      "ID of the Policy profile. If unspecified, Gatekeeper is disabled.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"quota_cpu_units": {
			Description:  "Maximum CPU units.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_disk_size": {
			Description:  "Maximum disk size in GBs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      2048, // 2048 GB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_ram_size": {
			Description:  "Maximum RAM size in GBs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      500, // 500 GB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_cpu_units": {
			Description:  "Maximum CPU units for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_volume_size": {
			Description:  "Maximum volume size in GBs for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      2000, // 2 TB
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_vm_ram_size": {
			Description:  "Maximum RAM size in GBs for standalone VMs.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      500, // 500 GB
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
			Set:          utils.HashAttributes("name", "disk_size", "flavor", "spot_server", "zone", "hypervisor"),
			Elem: &schema.Resource{
				Schema: taikunServerBasicSchema(),
			},
		},
		"server_kubemaster": {
			Description:  "Kubemaster server.",
			Type:         schema.TypeSet,
			Optional:     true,
			RequiredWith: []string{"server_bastion", "server_kubeworker"},
			Set:          utils.HashAttributes("name", "disk_size", "flavor", "kubernetes_node_label", "spot_server", "wasm", "hypervisor"),
			Elem: &schema.Resource{
				Schema: taikunServerSchemaWithKubernetesNodeLabels(),
			},
		},
		"server_kubeworker": {
			Description:  "Kubeworker server.",
			Type:         schema.TypeSet,
			Optional:     true,
			RequiredWith: []string{"server_bastion", "server_kubemaster"},
			Set:          utils.HashAttributes("name", "disk_size", "flavor", "kubernetes_node_label", "spot_server", "wasm", "zone", "hypervisor", "proxmox_extra_disk_size"),
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
		"autoscaler_name": {
			Description:  "Autoscaler group name (specify together with all other autoscaler parameters).",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(3, 10),
			RequiredWith: []string{"autoscaler_flavor", "autoscaler_disk_size", "autoscaler_max_size", "autoscaler_min_size"},
		},
		"autoscaler_flavor": {
			Description:  "Flavor of workers created by autoscaler (specify together with all other autoscaler parameters).",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"autoscaler_name", "autoscaler_disk_size", "autoscaler_max_size", "autoscaler_min_size"},
		},
		"autoscaler_disk_size": {
			Description:  "Disk size of autoscaler in GB (specify together with all other autoscaler parameters).",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(30),
			RequiredWith: []string{"autoscaler_name", "autoscaler_flavor", "autoscaler_max_size", "autoscaler_min_size"},
		},
		"autoscaler_min_size": {
			Description:  "Minimum number of workers created by autoscaler (specify together with all other autoscaler parameters).",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			RequiredWith: []string{"autoscaler_name", "autoscaler_flavor", "autoscaler_disk_size", "autoscaler_max_size"},
		},
		"autoscaler_max_size": {
			Description:  "Maximum number of workers created by autoscaler (specify together with all other autoscaler parameters).",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			RequiredWith: []string{"autoscaler_name", "autoscaler_flavor", "autoscaler_disk_size", "autoscaler_min_size"},
		},
		"autoscaler_spot_enabled": {
			Description:  "When enabled, autoscaler will use spot flavors for autoscaled workers (be sure to enable spot flavors for this project). If not specified, defaults to false.",
			Type:         schema.TypeBool,
			Optional:     true,
			Default:      false,
			RequiredWith: []string{"autoscaler_name", "autoscaler_flavor", "autoscaler_disk_size", "autoscaler_min_size", "autoscaler_max_size"},
		},
		"spot_full": {
			Description:   "When enabled, project will support full spot Kubernetes (controlplane + workers)",
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"spot_worker"},
		},
		"spot_worker": {
			Description:   "When enabled, project will support spot flavors for Kubernetes worker nodes",
			Type:          schema.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"spot_full"},
		},
		"spot_vms": {
			Description: "When enabled, project will support spot flavors of standalone VMs",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"spot_max_price": {
			Description: "Maximum spot price the user can set on servers/standalone VMs.",
			Type:        schema.TypeFloat,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
	}
}

func ResourceTaikunProject() *schema.Resource {
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

	cloudCredentialID, _ := utils.Atoi32(d.Get("cloud_credential_id").(string))
	body.SetCloudCredentialId(cloudCredentialID)
	flavorsData := d.Get("flavors").(*schema.Set).List()
	flavors := make([]string, len(flavorsData))
	for i, flavorData := range flavorsData {
		flavors[i] = flavorData.(string)
	}
	body.SetFlavors(flavors)

	if alertingProfileID, alertingProfileIDIsSet := d.GetOk("alerting_profile_id"); alertingProfileIDIsSet {
		alertingId, _ := utils.Atoi32(alertingProfileID.(string))
		body.SetAlertingProfileId(alertingId)
	}
	if backupCredentialID, backupCredentialIDIsSet := d.GetOk("backup_credential_id"); backupCredentialIDIsSet {
		body.SetIsBackupEnabled(true)
		backupCredential, _ := utils.Atoi32(backupCredentialID.(string))
		body.SetS3CredentialId(backupCredential)
	}
	if enableMonitoring, enableMonitoringIsSet := d.GetOk("monitoring"); enableMonitoringIsSet {
		body.SetIsMonitoringEnabled(enableMonitoring.(bool))
	}
	if deleteOnExpiration, deleteOnExpirationIsSet := d.GetOk("delete_on_expiration"); deleteOnExpirationIsSet {
		body.SetDeleteOnExpiration(deleteOnExpiration.(bool))
	}
	if expirationDate, expirationDateIsSet := d.GetOk("expiration_date"); expirationDateIsSet {
		dateTime := utils.DateToDateTime(expirationDate.(string))
		body.SetExpiredAt(time.Time(dateTime))
	} else {
		body.SetExpiredAtNil()
	}
	if PolicyProfileID, PolicyProfileIDIsSet := d.GetOk("policy_profile_id"); PolicyProfileIDIsSet {
		opaProfileID, _ := utils.Atoi32(PolicyProfileID.(string))
		body.SetOpaProfileId(opaProfileID)
	}

	if organizationID, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		projectOrganizationID, _ := utils.Atoi32(organizationID.(string))
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
		accessProfile, _ := utils.Atoi32(accessProfileID.(string))
		body.SetAccessProfileId(accessProfile)
	}

	if kubernetesProfileID, kubernetesProfileIDIsSet := d.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		kubeId, _ := utils.Atoi32(kubernetesProfileID.(string))
		body.SetKubernetesProfileId(kubeId)
	}

	if kubernetesVersion, kubernetesVersionIsSet := d.GetOk("kubernetes_version"); kubernetesVersionIsSet {
		body.SetKubernetesVersion(kubernetesVersion.(string))
	}

	// Spots
	spotFull, spotFullIsSet := d.GetOk("spot_full")
	spotWorker, spotWorkerIsSet := d.GetOk("spot_worker")
	spotVms, spotVmsIsSet := d.GetOk("spot_vms")
	spotMaxPrice, spotMaxPriceIsSet := d.GetOk("spot_max_price")

	if spotMaxPriceIsSet {
		if !spotFullIsSet && !spotWorkerIsSet && !spotVmsIsSet {
			return diag.Errorf("If you set max spot price, the project must have spots enabled.")
		}
		body.SetMaxSpotPrice(spotMaxPrice.(float64))
	}
	if spotFullIsSet {
		body.SetAllowFullSpotKubernetes(spotFull.(bool))
	}
	if spotWorkerIsSet {
		body.SetAllowSpotWorkers(spotWorker.(bool))
	}
	if spotVmsIsSet {
		body.SetAllowSpotVMs(spotVms.(bool))
	}

	// Autoscaler
	autoscalerName, autoscalerNameIsSet := d.GetOk("autoscaler_name")
	autoscalerFlavor, autoscalerFlavorIsSet := d.GetOk("autoscaler_flavor")
	autoscalerMin, autoscalerMinIsSet := d.GetOk("autoscaler_min_size")
	autoscalerMax, autoscalerMaxIsSet := d.GetOk("autoscaler_max_size")
	autoscalerDisk, autoscalerDiskIsSet := d.GetOk("autoscaler_disk_size")
	autoscalerSpot, autoscalerSpotIsSet := d.GetOk("autoscaler_spot_enabled")

	// If we specified a flavor not bound to project, Taikun would bind it for us
	// - terraform would be very confused by this since it did not bind the flavor.
	// For that reason Taikun TF proivder forbids to specify autoscaler flavor outside of bound flavors.
	if autoscalerFlavor != "" {
		desiredFlavorIsBound := false
		for _, oneFlavor := range flavors {
			if oneFlavor == autoscalerFlavor {
				desiredFlavorIsBound = true
			}
		}
		if !desiredFlavorIsBound {
			return diag.Errorf("Error: Autoscaler's flavor must be present in flavors already bound to project.")
		}
	}

	if autoscalerNameIsSet &&
		autoscalerFlavorIsSet &&
		autoscalerDiskIsSet &&
		autoscalerMinIsSet &&
		autoscalerMaxIsSet {

		body.SetAutoscalingGroupName(autoscalerName.(string))
		body.SetAutoscalingFlavor(autoscalerFlavor.(string))
		body.SetMinSize(int32(autoscalerMin.(int)))
		body.SetMaxSize(int32(autoscalerMax.(int)))
		body.SetDiskSize(utils.GibiByteToByte(autoscalerDisk.(int)))
		body.SetAutoscalingEnabled(true)

		if autoscalerSpotIsSet {
			body.SetAutoscalingSpotEnabled(autoscalerSpot.(bool))
		}
	}

	// Send project creation request
	response, responseBody, err := apiClient.Client.ProjectsAPI.ProjectsCreate(context.TODO()).CreateProjectCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(responseBody, err))
	}

	d.SetId(response.GetId())
	projectID, _ := utils.Atoi32(response.GetId())

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

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, d, meta)
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
		id32, err := utils.Atoi32(id)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, _, err := apiClient.Client.ServersAPI.ServersDetails(ctx, id32).Execute()
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		responseVM, _, err := apiClient.Client.StandaloneAPI.StandaloneDetails(ctx, id32).Execute()
		if err != nil {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
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
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		deleteOnExpiration, err := resourceTaikunProjectGetDeleteOnExpiration(projectDetailsDTO.GetProjectId(), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		projectMap := flattenTaikunProject(&projectDetailsDTO, serverList, vmList, boundFlavorDTOs, boundImageDTOs, &quotaResponse.Data[0], deleteOnExpiration)
		usernames := resourceTaikunProjectGetResourceDataVmUsernames(d)
		if err := utils.SetResourceDataFromMap(d, projectMap); err != nil {
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
	id, err := utils.Atoi32(d.Id())
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
			newAlertingProfileID, _ := utils.Atoi32(newAlertingProfileIDData.(string))
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
			dateTime := utils.DateToDateTime(expirationDate.(string))
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
					kubeWorkerId, _ := utils.Atoi32(kubeWorkerMap["id"].(string))
					serverIds = append(serverIds, kubeWorkerId)
				}

				deleteServerBody := tkcore.ProjectDeploymentDeleteServersCommand{}
				deleteServerBody.SetProjectId(id)
				deleteServerBody.SetServerIds(serverIds)

				_, err = apiClient.Client.ProjectDeploymentAPI.ProjectDeploymentDelete(ctx).ProjectDeploymentDeleteServersCommand(deleteServerBody).Execute()
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
					serverCreateBody.SetDiskSize(utils.GibiByteToByte(kubeWorkerMap["disk_size"].(int)))
					serverCreateBody.SetFlavor(kubeWorkerMap["flavor"].(string))
					serverCreateBody.SetKubernetesNodeLabels(resourceTaikunProjectServerKubernetesLabels(kubeWorkerMap))
					serverCreateBody.SetName(kubeWorkerMap["name"].(string))
					serverCreateBody.SetProjectId(id)
					serverCreateBody.SetRole(tkcore.CLOUDROLE_KUBEWORKER)
					serverCreateBody.SetWasmEnabled(kubeWorkerMap["wasm"].(bool))
					serverCreateBody.SetHypervisor(kubeWorkerMap["hypervisor"].(string))
					serverCreateBody.SetAvailabilityZone(kubeWorkerMap["zone"].(string))

					if kubeWorkerMap["proxmox_extra_disk_size"].(int) != 0 {
						proxmoxStorageString, err1 := utils.GetProxmoxStorageStringForServer(id, apiClient)
						if err1 != nil {
							return diag.FromErr(err1)
						}
						proxmoxRole, err2 := tkcore.NewProxmoxRoleFromValue(proxmoxStorageString)
						if err2 != nil {
							return diag.FromErr(err2)
						}
						proxmoxExtraDiskSize := int32(kubeWorkerMap["proxmox_extra_disk_size"].(int))
						serverCreateBody.SetProxmoxRole(*proxmoxRole)
						serverCreateBody.SetProxmoxExtraDiskSize(proxmoxExtraDiskSize)
					}

					serverCreateBody, err = resourceTaikunProjectSetServerSpots(kubeWorkerMap, serverCreateBody) // Spots
					if err != nil {
						return diag.Errorf("There was an error in server spot configuration")
					}

					serverCreateResponse, response, newErr := apiClient.Client.ServersAPI.ServersCreate(ctx).ServerForCreateDto(serverCreateBody).Execute()
					if newErr != nil {
						return diag.FromErr(tk.CreateError(response, newErr))
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

	// Spots for project
	spotFullChange := d.HasChange("spot_full")
	spotWorkerChange := d.HasChange("spot_worker")
	spotVmsChange := d.HasChange("spot_vms")

	// Vm spots do not collide with anything
	if spotVmsChange {
		if err = resourceTaikunProjectToggleVmsSpot(ctx, d, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}
	// Full and Worker chanege can collide if there was a change on remote
	if spotFullChange && spotWorkerChange {
		diag.Errorf("There has been a change in conflicting parameters spot_full and spot_worker")
	}
	if spotFullChange {
		if err = resourceTaikunProjectToggleFullSpot(ctx, d, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}
	if spotWorkerChange {
		if err = resourceTaikunProjectToggleWorkerSpot(ctx, d, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	// Autoscaler disable and enable with different flavor and autoscaling group name
	// Precedence: high, first we check if we must recreate the whole autoscaler
	iWishToDisable := d.Get("autoscaler_name") == "" || d.Get("autoscaler_flavor") == "" || d.Get("autoscaler_disk_size") == "0"
	iJustRecreated := false
	if d.HasChange("autoscaler_name") || d.HasChange("autoscaler_flavor") || d.HasChange("autoscaler_disk_size") || d.HasChange("autoscaler_spot_enabled") {
		if iWishToDisable {
			// Disable autoscaler
			if err := resourceTaikunProjectDisableAutoscaler(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
		} else {
			// Enable or Recreate autoscaler with changes
			if err := resourceTaikunProjectRecreateAutoscaler(ctx, d, apiClient); err != nil {
				return diag.FromErr(err)
			}
			iJustRecreated = true
		}
	}

	// Autoscaler edit
	// Precedence: medium. Worst case, autoscaler was just recreated
	if (d.HasChange("autoscaler_min_size") || d.HasChange("autoscaler_max_size")) && !iWishToDisable && !iJustRecreated {
		if err := resourceTaikunProjectUpdateAutoscaler(ctx, d, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	// Flavor changes must be checked after autoscaler
	// Precedence: low, autoscaler flavors should not interfere
	if d.HasChange("flavors") {
		if err = resourceTaikunProjectEditFlavors(d, apiClient, id); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunProjectLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, d, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
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

	// Get all autoscaler servers
	autoscalerData, _, err := apiClient.Client.ServersAPI.ServersList(ctx).AutoscalingGroup(d.Get("autoscaler_name").(string)).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	// Add the ids to the list of servers about to be deleted
	var i int32 = 0
	for ; i < autoscalerData.GetTotalCount(); i++ {
		serversToPurge = append(serversToPurge, map[string]interface{}{"id": fmt.Sprint(autoscalerData.GetData()[0].GetId())})
	}
	// errstring := fmt.Sprint(serversToPurge)
	// tflog.Error(ctx, errstring)

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
		body.SetServerRam(utils.GibiByteToByte(ram.(int)))
	}

	if disk, ok := d.GetOk("quota_disk_size"); ok {
		body.SetServerDiskSize(utils.GibiByteToByte(disk.(int)))
	}

	if vmCpu, ok := d.GetOk("quota_vm_cpu_units"); ok {
		body.SetVmCpu(int64(vmCpu.(int)))
	}

	if vmRam, ok := d.GetOk("quota_vm_ram_size"); ok {
		body.SetVmRam(utils.GibiByteToByte(vmRam.(int)))
	}

	if vmVolume, ok := d.GetOk("quota_vm_volume_size"); ok {
		body.SetVmVolumeSize(float64(vmVolume.(int))) // No conversion needed, API takes GBs
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

	// Flatten project Flavors
	flavors := make([]string, len(boundFlavorDTOs))
	for i, boundFlavorDTO := range boundFlavorDTOs {
		flavors[i] = boundFlavorDTO.GetName()
	}

	// Flatten project Images
	images := make([]string, len(boundImageDTOs))
	cloudType := projectDetailsDTO.GetCloudType()
	for i, boundImageDTO := range boundImageDTOs {
		if cloudType == tkcore.CLOUDTYPE_GOOGLE {
			// If GCP - Google uses image.name instead of image.id
			images[i] = boundImageDTO.GetName()
		} else {
			// All other CCs use image.id
			images[i] = boundImageDTO.GetImageId()
		}
	}

	// Flatten project attributes
	projectMap := map[string]interface{}{
		"access_ip":               projectDetailsDTO.GetAccessIp(),
		"access_profile_id":       utils.I32toa(projectDetailsDTO.GetAccessProfileId()),
		"alerting_profile_name":   projectDetailsDTO.GetAlertingProfileName(),
		"cloud_credential_id":     utils.I32toa(projectDetailsDTO.GetCloudId()),
		"monitoring":              projectDetailsDTO.GetIsMonitoringEnabled(),
		"delete_on_expiration":    projectDeleteOnExpiration,
		"expiration_date":         utils.Rfc3339DateTimeToDate(projectDetailsDTO.GetExpiredAt()),
		"flavors":                 flavors,
		"images":                  images,
		"id":                      utils.I32toa(projectDetailsDTO.GetProjectId()),
		"kubernetes_profile_id":   utils.I32toa(projectDetailsDTO.GetKubernetesProfileId()),
		"kubernetes_version":      projectDetailsDTO.GetKubernetesCurrentVersion(),
		"lock":                    projectDetailsDTO.GetIsLocked(),
		"name":                    projectDetailsDTO.GetProjectName(),
		"organization_id":         utils.I32toa(projectDetailsDTO.GetOrganizationId()),
		"quota_cpu_units":         projectQuotaDTO.GetServerCpu(),
		"quota_ram_size":          utils.ByteToGibiByte(projectQuotaDTO.GetServerRam()),
		"quota_disk_size":         utils.ByteToGibiByte(projectQuotaDTO.GetServerDiskSize()),
		"quota_vm_cpu_units":      projectQuotaDTO.GetVmCpu(),
		"quota_vm_ram_size":       utils.ByteToGibiByte(projectQuotaDTO.GetVmRam()),
		"quota_vm_volume_size":    projectQuotaDTO.GetVmVolumeSize(),
		"autoscaler_name":         projectDetailsDTO.GetAutoscalingGroupName(),
		"autoscaler_flavor":       projectDetailsDTO.GetFlavor(),
		"autoscaler_min_size":     projectDetailsDTO.GetMinSize(),
		"autoscaler_max_size":     projectDetailsDTO.GetMaxSize(),
		"autoscaler_disk_size":    utils.ByteToGibiByte(projectDetailsDTO.GetDiskSize()),
		"autoscaler_spot_enabled": projectDetailsDTO.GetIsAutoscalingSpotEnabled(),
		"spot_full":               projectDetailsDTO.GetAllowFullSpotKubernetes(),
		"spot_worker":             projectDetailsDTO.GetAllowSpotWorkers(),
		"spot_vms":                projectDetailsDTO.GetAllowSpotVMs(),
		"spot_max_price":          projectDetailsDTO.GetMaxSpotPrice(),
	}

	// Flatten Kubernetes servers
	bastions := make([]map[string]interface{}, 0)
	kubeMasters := make([]map[string]interface{}, 0)
	kubeWorkers := make([]map[string]interface{}, 0)
	skip_this_server := false
	for _, server := range serverListDTO {
		// Flatten server attributes for every server type
		serverMap := map[string]interface{}{
			"created_by":            server.GetCreatedBy(),
			"disk_size":             utils.ByteToGibiByte(server.GetDiskSize()),
			"id":                    utils.I32toa(server.GetId()),
			"ip":                    server.GetIpAddress(),
			"last_modified":         server.GetLastModified(),
			"last_modified_by":      server.GetLastModifiedBy(),
			"name":                  server.GetName(),
			"status":                server.GetStatus(),
			"spot_server":           server.GetSpotInstance(),
			"spot_server_max_price": server.GetSpotPrice(),
			"zone":                  utils.GetLastCharacter(server.GetAvailabilityZone()), // Get last character of the string
			"hypervisor":            server.GetHypervisor(),
		}

		// Attributes only for Workers
		serverRole := server.GetRole()
		if serverRole == tkcore.CLOUDROLE_KUBEWORKER {
			serverMap["wasm"] = server.GetWasmEnabled()
			serverMap["zone"] = utils.GetLastCharacter(server.GetAvailabilityZone())
			serverMap["proxmox_extra_disk_size"] = server.GetProxmoxExtraDiskSize()
		}
		// Attributes only for Masters
		if serverRole == tkcore.CLOUDROLE_KUBEMASTER {
			serverMap["wasm"] = server.GetWasmEnabled()
			serverMap["zone"] = utils.GetLastCharacter(server.GetAvailabilityZone())

		}

		// Flatten flavor
		switch server.GetCloudType() {
		case tkcore.CLOUDTYPE_AWS:
			serverMap["flavor"] = server.GetAwsInstanceType()
		case tkcore.CLOUDTYPE_AZURE:
			serverMap["flavor"] = server.GetAzureVmSize()
		case tkcore.CLOUDTYPE_OPENSTACK:
			serverMap["flavor"] = server.GetOpenstackFlavor()
		case tkcore.CLOUDTYPE_GOOGLE, "google":
			serverMap["flavor"] = server.GetGoogleMachineType()
		case tkcore.CLOUDTYPE_PROXMOX, "proxmox":
			serverMap["flavor"] = server.GetProxmoxFlavor()
		case tkcore.CLOUDTYPE_VSPHERE, "vsphere":
			serverMap["flavor"] = server.GetVsphereFlavor()
		}

		if serverRole == tkcore.CLOUDROLE_BASTION {
			// Flatten bastion
			bastions = append(bastions, serverMap)
		} else {
			// Flatten masters and workers with labels
			labels := make([]map[string]interface{}, len(server.GetKubernetesNodeLabels()))
			for i, rawLabel := range server.GetKubernetesNodeLabels() {
				labels[i] = map[string]interface{}{
					"key":   *rawLabel.Key.Get(),
					"value": *rawLabel.Value.Get(),
				}

				// Autoscaler If label shows the node is created by autoscaler - ignore it, for TF it is invisible.
				if *rawLabel.Key.Get() == "taikun.cloud/autoscaling-group" {
					skip_this_server = true
				}
			}

			// Add server to state only if its not an autoscaler server
			if !skip_this_server {
				serverMap["kubernetes_node_label"] = labels

				if serverRole == tkcore.CLOUDROLE_KUBEMASTER {
					kubeMasters = append(kubeMasters, serverMap)
				} else {
					kubeWorkers = append(kubeWorkers, serverMap)
				}
			}
		}
	}
	projectMap["server_bastion"] = bastions
	projectMap["server_kubemaster"] = kubeMasters
	projectMap["server_kubeworker"] = kubeWorkers

	// Flatten project VMs
	vms := make([]map[string]interface{}, 0)
	for _, vm := range vmListDTO {
		vmMap := map[string]interface{}{
			"access_ip":             vm.GetPublicIp(),
			"cloud_init":            vm.GetCloudInit(),
			"created_by":            vm.GetCreatedBy(),
			"flavor":                vm.GetTargetFlavor(),
			"id":                    utils.I32toa(vm.GetId()),
			"image_id":              vm.GetImageId(),
			"image_name":            vm.GetImageName(),
			"ip":                    vm.GetIpAddress(),
			"last_modified":         vm.GetLastModified(),
			"last_modified_by":      vm.GetLastModifiedBy(),
			"name":                  vm.GetName(),
			"public_ip":             vm.GetPublicIpEnabled(),
			"standalone_profile_id": utils.I32toa(vm.Profile.GetId()),
			"status":                vm.GetStatus(),
			"volume_size":           vm.GetVolumeSize(),
			"volume_type":           vm.GetVolumeType(),
			"spot_vm":               vm.GetSpotInstance(),
			"spot_vm_max_price":     vm.GetSpotPrice(),
			"hypervisor":            vm.GetHypervisor(),
			"zone":                  utils.GetLastCharacter(vm.GetAvailabilityZone()),
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
				//"device_name": rawDisk.GetDeviceName(),
				"id":          utils.I32toa(rawDisk.GetId()),
				"name":        rawDisk.GetName(),
				"size":        rawDisk.GetCurrentSize(),
				"volume_type": rawDisk.GetVolumeType(),
			}
		}
		vmMap["disk"] = disks

		vms = append(vms, vmMap)
	}
	projectMap["vm"] = vms

	// Flatten alerting profiles
	var nullID int32
	if projectDetailsDTO.GetAlertingProfileId() != nullID {
		projectMap["alerting_profile_id"] = utils.I32toa(projectDetailsDTO.GetAlertingProfileId())
	}

	// Flatten Backups
	if projectDetailsDTO.GetIsBackupEnabled() {
		projectMap["backup_credential_id"] = utils.I32toa(projectDetailsDTO.GetS3CredentialId())
	}

	// Flatten OPA profiles
	if projectDetailsDTO.GetIsOpaEnabled() {
		projectMap["policy_profile_id"] = utils.I32toa(projectDetailsDTO.GetOpaProfileId())
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
	createStateConf := &retry.StateChangeConf{
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
		kubernetesProfileID, _ := utils.Atoi32(kubernetesProfileIDData.(string))
		lbSolution, err := resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID, apiClient)
		if err != nil {
			return err
		}
		if lbSolution == utils.LoadBalancerTaikun {
			cloudCredentialID, _ := utils.Atoi32(d.Get("cloud_credential_id").(string))
			cloudType, err := ResourceTaikunProjectGetCloudType(cloudCredentialID, apiClient)
			if err != nil {
				return err
			}
			if _, taikunLBFlavorIsSet := d.GetOk("taikun_lb_flavor"); !taikunLBFlavorIsSet {
				return fmt.Errorf("if Taikun load balancer is enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor must be set")
			}
			if cloudType != utils.CloudTypeOpenStack {
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
	return utils.GetLoadBalancingSolution(kubernetesProfile.GetOctaviaEnabled(), kubernetesProfile.GetTaikunLBEnabled()), nil
}

func ResourceTaikunProjectGetCloudType(cloudCredentialID int32, apiClient *tk.Client) (string, error) {
	// Check if Cloud credential is Openstack
	responseOS, _, err := apiClient.Client.OpenstackCloudCredentialAPI.OpenstackList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responseOS.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_OPENSTACK), nil
	}

	// Check if CC is AWS
	responseAWS, _, err := apiClient.Client.AWSCloudCredentialAPI.AwsList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responseAWS.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_AWS), nil
	}

	// Check if CC is Azure
	responseAZ, _, err := apiClient.Client.AzureCloudCredentialAPI.AzureList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responseAZ.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_AZURE), nil
	}

	// Check if CC is Google
	responseGCP, _, err := apiClient.Client.GoogleAPI.GooglecloudList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responseGCP.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_GOOGLE), nil
	}

	// Check if CC is Proxmox
	responsePROXMOX, _, err := apiClient.Client.ProxmoxCloudCredentialAPI.ProxmoxList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responsePROXMOX.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_PROXMOX), nil
	}

	// Check if CC is vSphere
	responseVSPHERE, _, err := apiClient.Client.VsphereCloudCredentialAPI.VsphereList(context.TODO()).Id(cloudCredentialID).Execute()
	if err != nil {
		return "", err
	} else if responseVSPHERE.GetTotalCount() == 1 {
		return string(tkcore.CLOUDTYPE_VSPHERE), nil
	}

	// Unknown CC type
	return "", fmt.Errorf("cloud credential with ID %d not found", cloudCredentialID)
}

func resourceTaikunProjectLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.ProjectLockManagerCommand{}
	body.SetId(id)
	body.SetMode(utils.GetLockMode(lock))
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
