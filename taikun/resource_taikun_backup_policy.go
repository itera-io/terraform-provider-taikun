package taikun

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/backup"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBackupPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cron_period": {
			Description:      "Frequency of backups.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsCron,
		},
		"excluded_namespaces": {
			Description: "Namespaces excluded from the backups.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			DefaultFunc: func() (interface{}, error) {
				return []string{}, nil
			},
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"included_namespaces": {
			Description: "Namespaces included in the backups.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			DefaultFunc: func() (interface{}, error) {
				return []string{}, nil
			},
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"name": {
			Description:  "The name of the backup policy.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"phase": {
			Description: "The phase of the backup policy.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"project_id": {
			Description:      "The ID of the project.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"retention_period": {
			Description: "How long to store the backups.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "720h",
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^(((0*[1-9][0-9]*)h)?((0*[1-9][0-9]*)m)?((0*[1-9][0-9]*)s)|((0*[1-9][0-9]*)h)?((0*[1-9][0-9]*)m)((\\d+)s)?|((0*[1-9][0-9]*)h)((\\d+)m)?((\\d+)s)?)$`),
				"The retention period must follow the HMS format, for example: `10h30m15s`, `48h5s` or `360h`.",
			),
		},
	}
}

func resourceTaikunBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Backup Policy",
		CreateContext: resourceTaikunBackupPolicyCreate,
		ReadContext:   generateResourceTaikunBackupPolicyReadWithoutRetries(),
		DeleteContext: resourceTaikunBackupPolicyDelete,
		Schema:        resourceTaikunBackupPolicySchema(),
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {

			if len(diff.Get("included_namespaces").([]interface{})) == 0 && len(diff.Get("excluded_namespaces").([]interface{})) == 0 {
				return fmt.Errorf("please specify include or exclude namespace to create backup")
			}

			return nil
		},
	}
}

func resourceTaikunBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	projectId, _ := atoi32(d.Get("project_id").(string))

	body := &models.CreateBackupPolicyCommand{
		CronPeriod:        stringAddress(d.Get("cron_period")),
		IncludeNamespaces: resourceGetStringList(d.Get("included_namespaces")),
		Name:              stringAddress(d.Get("name")),
		ProjectID:         int32Address(projectId),
		RetentionPeriod:   stringAddress(d.Get("retention_period")),
	}

	params := backup.NewBackupCreateParams().WithV(ApiVersion).WithBody(body)
	_, err := apiClient.Client.Backup.BackupCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d/%s", projectId, d.Get("name").(string)))

	return readAfterCreateWithRetries(generateResourceTaikunBackupPolicyReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunBackupPolicyReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunBackupPolicyRead(true)
}
func generateResourceTaikunBackupPolicyReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunBackupPolicyRead(false)
}
func generateResourceTaikunBackupPolicyRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		projectId, backupPolicyName, err := parseBackupPolicyId(d.Id())
		if err != nil {
			return diag.Errorf("Error while reading taikun_backup_policy : %s", err)
		}
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.Backup.BackupListAllSchedules(backup.NewBackupListAllSchedulesParams().WithV(ApiVersion).WithProjectID(projectId), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, policy := range response.Payload.Data {
			if policy.MetadataName == backupPolicyName {

				err = setResourceDataFromMap(d, flattenTaikunBackupPolicy(policy))
				if err != nil {
					return diag.FromErr(err)
				}

				d.SetId(fmt.Sprintf("%d/%s", projectId, backupPolicyName))

				return nil
			}
		}

		if withRetries {
			d.SetId(fmt.Sprintf("%d/%s", projectId, backupPolicyName))
			return diag.Errorf(notFoundAfterCreateOrUpdateError)
		}
		return nil
	}
}

func resourceTaikunBackupPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	projectId, backupPolicyName, err := parseBackupPolicyId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_backup_policy : %s", err)
	}

	deleteBody := &models.DeleteScheduleCommand{
		Name:      backupPolicyName,
		ProjectID: projectId,
	}
	params := backup.NewBackupDeleteScheduleParams().WithV(ApiVersion).WithBody(deleteBody)
	_, err = apiClient.Client.Backup.BackupDeleteSchedule(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenTaikunBackupPolicy(rawBackupPolicy *models.CScheduleDto) map[string]interface{} {

	return map[string]interface{}{
		"cron_period":         rawBackupPolicy.Schedule,
		"excluded_namespaces": rawBackupPolicy.ExcludedNamespaces,
		"included_namespaces": rawBackupPolicy.IncludedNamespaces,
		"name":                rawBackupPolicy.MetadataName,
		"phase":               rawBackupPolicy.Phase,
		"retention_period":    rawBackupPolicy.TTL,
	}
}

func parseBackupPolicyId(id string) (int32, string, error) {
	list := strings.Split(id, "/")
	if len(list) != 2 {
		return 0, "", fmt.Errorf("unable to determine taikun_backup_policy ID")
	}

	projectId, err := atoi32(list[0])
	if err != nil {
		return 0, "", fmt.Errorf("unable to determine taikun_backup_policy ID")
	}

	backupPolicyName := list[1]

	return projectId, backupPolicyName, nil
}
