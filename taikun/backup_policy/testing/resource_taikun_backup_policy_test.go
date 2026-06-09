package testing

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

// HCL configuration template for taikun_backup_policy.
// Requires: project_id, name, cron_period, retention_period, and at least one included_namespaces.
const testAccResourceTaikunBackupPolicyConfig = `
resource "taikun_backup_policy" "foo" {
  project_id          = "%s"
  name                = "%s"
  cron_period         = "0 0 * * 0"
  retention_period    = "720h"
  included_namespaces = ["default"]
}
`

// TestAccResourceTaikunBackupPolicy verifies the backup policy lifecycle:
// creation, validation of attributes, read, and deletion of backup schedules.
// It is skipped if TAIKUN_PROJECT_ID is not set, as it requires an active k8s project
// with backup/Velero enabled on the platform.
func TestAccResourceTaikunBackupPolicy(t *testing.T) {
	projectID := os.Getenv("TAIKUN_PROJECT_ID")
	if projectID == "" {
		t.Skip("TAIKUN_PROJECT_ID must be set to a running k8s project ID with backups enabled to run this test")
	}

	backupPolicyName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupPolicyConfig,
					projectID,
					backupPolicyName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupPolicyExists,
					resource.TestCheckResourceAttrSet("taikun_backup_policy.foo", "id"),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "name", backupPolicyName),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "project_id", projectID),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "cron_period", "0 0 * * 0"),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "retention_period", "720h"),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "included_namespaces.#", "1"),
					resource.TestCheckResourceAttr("taikun_backup_policy.foo", "included_namespaces.0", "default"),
				),
			},
		},
	})
}

// Check that the backup policy exists in the Taikun API by listing schedules.
func testAccCheckTaikunBackupPolicyExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_policy" {
			continue
		}

		list := strings.Split(rs.Primary.ID, "/")
		if len(list) != 2 {
			return fmt.Errorf("invalid backup policy ID format in state: %s", rs.Primary.ID)
		}
		projectID, _ := utils.Atoi32(list[0])
		backupPolicyName := list[1]

		response, _, err := client.Client.BackupPolicyAPI.BackupListAllSchedules(context.TODO(), projectID).Limit(4000).Execute()
		if err != nil {
			return fmt.Errorf("failed to list backup policies for project %d: %w", projectID, err)
		}

		found := false
		for _, policy := range response.Data {
			if policy.GetMetadataName() == backupPolicyName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("backup policy %s not found on project %d", backupPolicyName, projectID)
		}
	}

	return nil
}

// Verify that the backup policy has been deleted/destroyed correctly in the Taikun API.
func testAccCheckTaikunBackupPolicyDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_policy" {
			continue
		}

		list := strings.Split(rs.Primary.ID, "/")
		if len(list) != 2 {
			return fmt.Errorf("invalid backup policy ID format in state: %s", rs.Primary.ID)
		}
		projectID, _ := utils.Atoi32(list[0])
		backupPolicyName := list[1]

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			response, _, err := client.Client.BackupPolicyAPI.BackupListAllSchedules(context.TODO(), projectID).Limit(4000).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}

			for _, policy := range response.Data {
				if policy.GetMetadataName() == backupPolicyName {
					return retry.RetryableError(errors.New("backup policy still exists"))
				}
			}
			return nil
		})

		if utils.TimedOut(retryErr) {
			return errors.New("backup policy still exists (timed out waiting for deletion)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
