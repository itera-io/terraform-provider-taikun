package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunBackupCredentialConfig = `
resource "taikun_backup_credential" "foo" {
  name            = "%s"
  lock       = %t

  s3_endpoint = "%s"
  s3_region   = "%s"
}
`

func TestAccResourceTaikunBackupCredential(t *testing.T) {
	backupCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckS3(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("S3_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("S3_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBackupCredentialLock(t *testing.T) {
	backupCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckS3(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("S3_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("S3_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					true,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("S3_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("S3_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBackupCredentialRename(t *testing.T) {
	backupCredentialName := utils.RandomTestName()
	newBackupCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckS3(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBackupCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					backupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", backupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("S3_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("S3_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBackupCredentialConfig,
					newBackupCredentialName,
					false,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBackupCredentialExists,
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "name", newBackupCredentialName),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_access_key_id", os.Getenv("S3_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_secret_access_key", os.Getenv("S3_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_endpoint", os.Getenv("S3_ENDPOINT")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "s3_region", os.Getenv("S3_REGION")),
					resource.TestCheckResourceAttr("taikun_backup_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_backup_credential.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunBackupCredentialExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_credential" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.S3CredentialsAPI.S3credentialsList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("backup credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBackupCredentialDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_backup_credential" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.S3CredentialsAPI.S3credentialsList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("backup credential still exists ()"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("backup credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
