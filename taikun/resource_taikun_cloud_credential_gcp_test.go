package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  billing_account_id = "%s"
  folder_id = "%s"
  region = "%s"
  az_count = "%d"
}
`

func TestAccResourceTaikunCloudCredentialGCP(t *testing.T) {
	cloudCredentialName := randomTestName()
	azCount, _ := atoi32(os.Getenv("GCP_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					os.Getenv("GCP_BILLING_ACCOUNT"),
					os.Getenv("GCP_FOLDER_ID"),
					os.Getenv("GCP_REGION"),
					azCount,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "billing_account_id", os.Getenv("GCP_BILLING_ACCOUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "config_file", "./gcp.json"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "folder_id", os.Getenv("GCP_FOLDER_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", "false"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "az_count", os.Getenv("GCP_AZ_COUNT")),
				),
			},
		},
	})
}

const testAccResourceTaikunCloudCredentialGCPImportProjectConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  az_count = "%d"
  lock = true
}
`

func TestAccResourceTaikunCloudCredentialGCPImportProject(t *testing.T) {
	cloudCredentialName := randomTestName()
	azCount, _ := atoi32(os.Getenv("GCP_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPImportProjectConfig,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
					azCount,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "config_file", "./gcp.json"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", "true"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "az_count", os.Getenv("GCP_AZ_COUNT")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialGCPExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := client.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCountGoogle() != 1 {
			return fmt.Errorf("gcp cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialGCPDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := client.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCountGoogle() != 0 {
				return resource.RetryableError(errors.New("gcp cloud credential still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("gcp cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
