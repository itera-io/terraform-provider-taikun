package cc_gcp_tests

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

const testAccResourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  lock = false
  az_count = %s
}
`

func TestAccResourceTaikunCloudCredentialGCP(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckGCP(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
					os.Getenv("GCP_AZ_COUNT"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "config_file", "./gcp.json"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", "true"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "az_count", os.Getenv("GCP_AZ_COUNT")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialGCPExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.GoogleAPI.GooglecloudList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("gcp cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialGCPDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.GoogleAPI.GooglecloudList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("gcp cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("gcp cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
