package taikun

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/itera-io/taikungoclient/client/opa_profiles"

	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/client/project_quotas"
	"github.com/itera-io/taikungoclient/client/users"

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
		"kubernetes_profile_id": {
			Description:      "ID of the project's Kubernetes profile. Defaults to the default Kubernetes profile of the project's organization.",
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
		"policy_profile_id": {
			Description:      "ID of the Policy profile. If unspecified, Gatekeeper is disabled.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: stringIsInt,
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
			Description:  "Maximum CPU units. Unlimited if unspecified.",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_disk_size": {
			Description:  "Maximum disk size in GBs. Unlimited if unspecified.",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"quota_id": {
			Description: "ID of the project quota.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"quota_ram_size": {
			Description:  "Maximum RAM size in GBs. Unlimited if unspecified.",
			Type:         schema.TypeInt,
			Optional:     true,
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
	}
}

func taikunServerKubeworkerSchema() map[string]*schema.Schema {
	kubeworkerSchema := taikunServerSchemaWithKubernetesNodeLabels()
	removeForceNewsFromSchema(kubeworkerSchema)
	return kubeworkerSchema
}

func taikunServerSchemaWithKubernetesNodeLabels() map[string]*schema.Schema {
	serverSchema := taikunServerBasicSchema()
	serverSchema["kubernetes_node_label"] = &schema.Schema{
		Description: "Attach Kubernetes node labels.",
		Type:        schema.TypeList,
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

	var projectOrganizationID int32 = -1

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
	if PolicyProfileID, PolicyProfileIDIsSet := data.GetOk("policy_profile_id"); PolicyProfileIDIsSet {
		body.OpaProfileID, _ = atoi32(PolicyProfileID.(string))
	}

	if organizationID, organizationIDIsSet := data.GetOk("organization_id"); organizationIDIsSet {
		projectOrganizationID, _ = atoi32(organizationID.(string))
		body.OrganizationID = projectOrganizationID
	}

	if taikunLBFlavor, taikunLBFlavorIsSet := data.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
		body.TaikunLBFlavor = taikunLBFlavor.(string)
		body.RouterIDStartRange = int32(data.Get("router_id_start_range").(int))
		body.RouterIDEndRange = int32(data.Get("router_id_end_range").(int))
	}
	if err := resourceTaikunProjectValidateKubernetesProfileLB(data, apiClient); err != nil {
		return diag.FromErr(err)
	}

	if accessProfileID, accessProfileIDIsSet := data.GetOk("access_profile_id"); accessProfileIDIsSet {
		body.AccessProfileID, _ = atoi32(accessProfileID.(string))
	} else {
		if projectOrganizationID == -1 {
			if err := resourceTaikunProjectGetDefaultOrganization(&projectOrganizationID, apiClient); err != nil {
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
	if kubernetesProfileID, kubernetesProfileIDIsSet := data.GetOk("kubernetes_profile_id"); kubernetesProfileIDIsSet {
		body.KubernetesProfileID, _ = atoi32(kubernetesProfileID.(string))
	} else {
		if projectOrganizationID == -1 {
			if err := resourceTaikunProjectGetDefaultOrganization(&projectOrganizationID, apiClient); err != nil {
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

	params := projects.NewProjectsCreateParams().WithV(ApiVersion).WithBody(&body)
	response, err := apiClient.client.Projects.ProjectsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.Payload.ID)
	projectID, _ := atoi32(response.Payload.ID)

	_, quotaCPUIsSet := data.GetOk("quota_cpu_units")
	_, quotaDiskIsSet := data.GetOk("quota_disk_size")
	_, quotaRAMIsSet := data.GetOk("quota_ram_size")
	if quotaCPUIsSet || quotaDiskIsSet || quotaRAMIsSet {

		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
		response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = resourceTaikunProjectEditQuotas(data, apiClient, response.Payload.Project.QuotaID); err != nil {
			return diag.FromErr(err)
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

		if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, projectID); err != nil {
			return diag.FromErr(err)
		}
	}

	if data.Get("lock").(bool) {
		if err := resourceTaikunProjectLock(projectID, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, data, meta)
}
func generateResourceTaikunProjectReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectRead(true)
}
func generateResourceTaikunProjectReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunProjectRead(false)
}
func generateResourceTaikunProjectRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id := data.Id()
		id32, err := atoi32(id)
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		params := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(id32)
		response, err := apiClient.client.Servers.ServersDetails(params, apiClient)
		if err != nil {
			if withRetries {
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
			if withRetries {
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

	if err := resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
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
	if data.HasChanges("quota_cpu_units", "quota_disk_size", "quota_ram_size") {
		quotaId, _ := atoi32(data.Get("quota_id").(string))

		if err := resourceTaikunProjectEditQuotas(data, apiClient, quotaId); err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("server_bastion") {
		oldBastions, newBastions := data.GetChange("server_bastion")
		oldSet := oldBastions.(*schema.Set)
		newSet := newBastions.(*schema.Set)

		if oldSet.Len() == 0 {
			// The project was empty before
			if err := resourceTaikunProjectUpdateToggleServices(ctx, data, apiClient); err != nil {
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
			serversToPurge := resourceTaikunProjectFlattenServersData(oldBastions, oldKubeMasters, oldKubeWorkers)
			err = resourceTaikunProjectPurgeServers(serversToPurge, apiClient, id)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := resourceTaikunProjectUpdateToggleServices(ctx, data, apiClient); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		if err := resourceTaikunProjectUpdateToggleServices(ctx, data, apiClient); err != nil {
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

				if err := resourceTaikunProjectCommit(apiClient, id); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if err := resourceTaikunProjectWaitForStatus(ctx, []string{"Ready"}, []string{"Updating", "Pending"}, apiClient, id); err != nil {
		return diag.FromErr(err)
	}

	if data.Get("lock").(bool) {
		if err := resourceTaikunProjectLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunProjectReadWithRetries(), ctx, data, meta)
}

func resourceTaikunProjectDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceTaikunProjectUnlockIfLocked(id, apiClient); err != nil {
		return diag.FromErr(err)
	}

	serversToPurge := resourceTaikunProjectFlattenServersData(
		data.Get("server_bastion"),
		data.Get("server_kubemaster"),
		data.Get("server_kubeworker"),
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

	// Delete the project
	body := models.DeleteProjectCommand{ProjectID: id, IsForceDelete: false}
	params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&body)
	if _, _, err := apiClient.client.Projects.ProjectsDelete(params, apiClient); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceTaikunProjectUnlockIfLocked(projectID int32, apiClient *apiClient) error {
	readParams := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(projectID)
	response, err := apiClient.client.Servers.ServersDetails(readParams, apiClient)
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

func resourceTaikunProjectUpdateToggleServices(ctx context.Context, data *schema.ResourceData, apiClient *apiClient) error {
	if err := resourceTaikunProjectUpdateToggleMonitoring(ctx, data, apiClient); err != nil {
		return err
	}
	if err := resourceTaikunProjectUpdateToggleBackup(ctx, data, apiClient); err != nil {
		return err
	}
	if err := resourceTaikunProjectUpdateToggleOPA(ctx, data, apiClient); err != nil {
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
				return fmt.Errorf("error waiting for project (%s) to disable backup: %s", data.Id(), err)
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

func resourceTaikunProjectUpdateToggleOPA(ctx context.Context, data *schema.ResourceData, apiClient *apiClient) error {
	if data.HasChange("policy_profile_id") {
		projectID, _ := atoi32(data.Id())
		oldOPAProfile, _ := data.GetChange("policy_profile_id")

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

		newOPAProfile, newOPAProfileIsSet := data.GetOk("policy_profile_id")

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
				return fmt.Errorf("error waiting for project (%s) to disable OPA: %s", data.Id(), err)
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

func resourceTaikunProjectEditQuotas(data *schema.ResourceData, apiClient *apiClient, quotaID int32) error {

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
		quotaEditBody.DiskSize = gibiByteToByte(quotaDisk.(int))
		quotaEditBody.IsDiskSizeUnlimited = false
	}

	if quotaRAM, quotaRAMIsSet := data.GetOk("quota_ram_size"); quotaRAMIsSet {
		quotaEditBody.RAM = gibiByteToByte(quotaRAM.(int))
		quotaEditBody.IsRAMUnlimited = false
	}

	quotaEditParams := project_quotas.NewProjectQuotasEditParams().WithV(ApiVersion).WithQuotaID(quotaID).WithBody(quotaEditBody)
	_, err := apiClient.client.ProjectQuotas.ProjectQuotasEdit(quotaEditParams, apiClient)
	if err != nil {
		return err
	}
	return nil
}

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

	if !projectQuotaDTO.IsCPUUnlimited {
		projectMap["quota_cpu_units"] = projectQuotaDTO.CPU
	}

	if !projectQuotaDTO.IsDiskSizeUnlimited {
		projectMap["quota_disk_size"] = byteToGibiByte(projectQuotaDTO.DiskSize)
	}

	if !projectQuotaDTO.IsRAMUnlimited {
		projectMap["quota_ram_size"] = byteToGibiByte(projectQuotaDTO.RAM)
	}

	return projectMap
}

func resourceTaikunProjectGetBoundFlavorDTOs(projectID int32, apiClient *apiClient) ([]*models.BoundFlavorsForProjectsListDto, error) {
	var boundFlavorDTOs []*models.BoundFlavorsForProjectsListDto
	boundFlavorsParams := flavors.NewFlavorsGetSelectedFlavorsForProjectParams().WithV(ApiVersion).WithProjectID(&projectID)
	for {
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

func resourceTaikunProjectFlattenServersData(bastionsData interface{}, kubeMastersData interface{}, kubeWorkersData interface{}) []interface{} {
	var servers []interface{}
	servers = append(servers, bastionsData.(*schema.Set).List()...)
	servers = append(servers, kubeMastersData.(*schema.Set).List()...)
	servers = append(servers, kubeWorkersData.(*schema.Set).List()...)
	return servers
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

func resourceTaikunProjectValidateKubernetesProfileLB(data *schema.ResourceData, apiClient *apiClient) error {
	if kubernetesProfileIDData, kubernetesProfileIsSet := data.GetOk("kubernetes_profile_id"); kubernetesProfileIsSet {
		kubernetesProfileID, _ := atoi32(kubernetesProfileIDData.(string))
		lbSolution, err := resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID, apiClient)
		if err != nil {
			return err
		}
		if lbSolution == loadBalancerTaikun {
			cloudCredentialID, _ := atoi32(data.Get("cloud_credential_id").(string))
			cloudType, err := resourceTaikunProjectGetCloudType(cloudCredentialID, apiClient)
			if err != nil {
				return err
			}
			if _, taikunLBFlavorIsSet := data.GetOk("taikun_lb_flavor"); !taikunLBFlavorIsSet {
				return fmt.Errorf("If Taikun load balancer is enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor must be set")
			}
			if cloudType != cloudTypeOpenStack {
				return fmt.Errorf("If Taikun load balancer is enabled, cloud type should be OpenStack; is %s", cloudType)
			}
		} else if _, taikunLBFlavorIsSet := data.GetOk("taikun_lb_flavor"); taikunLBFlavorIsSet {
			return fmt.Errorf("If Taikun load balancer is not enabled, router_id_start_range, router_id_end_range and taikun_lb_flavor should not be set")
		}
	}
	return nil
}

func resourceTaikunProjectGetKubernetesLBSolution(kubernetesProfileID int32, apiClient *apiClient) (string, error) {
	params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&kubernetesProfileID)
	response, err := apiClient.client.KubernetesProfiles.KubernetesProfilesList(params, apiClient)
	if err != nil {
		return "", err
	}
	if len(response.Payload.Data) == 0 {
		return "", fmt.Errorf("kubernetes profile with ID %d not found", kubernetesProfileID)
	}
	kubernetesProfile := response.Payload.Data[0]
	return getLoadBalancingSolution(kubernetesProfile.OctaviaEnabled, kubernetesProfile.TaikunLBEnabled), nil
}

func resourceTaikunProjectGetCloudType(cloudCredentialID int32, apiClient *apiClient) (string, error) {
	params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&cloudCredentialID)
	response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
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

func resourceTaikunProjectGetDefaultOrganization(defaultOrganizationID *int32, apiClient *apiClient) error {
	params := users.NewUsersDetailsParams().WithV(ApiVersion)
	response, err := apiClient.client.Users.UsersDetails(params, apiClient)
	if err != nil {
		return err
	}
	*defaultOrganizationID = response.Payload.Data.OrganizationID
	return nil
}

const defaultAccessProfileName = "default"

func resourceTaikunProjectGetDefaultAccessProfile(organizationID int32, apiClient *apiClient) (accessProfileID int32, found bool, err error) {
	params := access_profiles.NewAccessProfilesAccessProfilesForOrganizationListParams().WithV(ApiVersion).WithOrganizationID(&organizationID)
	response, err := apiClient.client.AccessProfiles.AccessProfilesAccessProfilesForOrganizationList(params, apiClient)
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

func resourceTaikunProjectGetDefaultKubernetesProfile(organizationID int32, apiClient *apiClient) (kubernetesProfileID int32, found bool, err error) {
	params := kubernetes_profiles.NewKubernetesProfilesKubernetesProfilesForOrganizationListParams().WithV(ApiVersion).WithOrganizationID(&organizationID)
	response, err := apiClient.client.KubernetesProfiles.KubernetesProfilesKubernetesProfilesForOrganizationList(params, apiClient)
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

func resourceTaikunProjectLock(id int32, lock bool, apiClient *apiClient) error {
	lockMode := getLockMode(lock)
	params := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&id).WithMode(&lockMode)
	_, err := apiClient.client.Projects.ProjectsLockManager(params, apiClient)
	return err
}
