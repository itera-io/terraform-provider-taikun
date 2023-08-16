package taikun

import (
	"context"
	"fmt"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	apiClient := meta.(*tk.Client)

	projectId, _ := atoi32(d.Get("project_id").(string))
	body := tkcore.CreateBackupPolicyCommand{}
	body.SetCronPeriod(d.Get("cron_period").(string))
	body.SetIncludeNamespaces(resourceGetStringList(d.Get("included_namespaces")))
	body.SetName(d.Get("name").(string))
	body.SetProjectId(projectId)
	body.SetRetentionPeriod(d.Get("retention_period").(string))

	res, err := apiClient.Client.BackupPolicyApi.BackupCreate(ctx).CreateBackupPolicyCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
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
		apiClient := meta.(*tk.Client)
		projectId, backupPolicyName, err := parseBackupPolicyId(d.Id())
		if err != nil {
			return diag.Errorf("Error while reading taikun_backup_policy : %s", err)
		}
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		// maybe add search?
		response, res, err := apiClient.Client.BackupPolicyApi.BackupListAllSchedules(context.TODO(), projectId).Limit(4000).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		for _, policy := range response.Data {
			if policy.GetMetadataName() == backupPolicyName {

				err = setResourceDataFromMap(d, flattenTaikunBackupPolicy(&policy))
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
	apiClient := meta.(*tk.Client)
	projectId, backupPolicyName, err := parseBackupPolicyId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_backup_policy : %s", err)
	}

	deleteBody := tkcore.DeleteScheduleCommand{}
	deleteBody.SetName(backupPolicyName)
	deleteBody.SetProjectId(projectId)

	res, err := apiClient.Client.BackupPolicyApi.BackupDeleteSchedule(context.TODO()).DeleteScheduleCommand(deleteBody).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunBackupPolicy(rawBackupPolicy *tkcore.CScheduleDto) map[string]interface{} {

	return map[string]interface{}{
		"cron_period":         rawBackupPolicy.GetSchedule(),
		"included_namespaces": rawBackupPolicy.GetIncludedNamespaces(),
		"name":                rawBackupPolicy.GetMetadataName(),
		"phase":               rawBackupPolicy.GetPhase(),
		"retention_period":    rawBackupPolicy.GetTtl(),
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
