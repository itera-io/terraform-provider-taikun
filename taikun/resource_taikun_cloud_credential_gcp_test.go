package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
)

const testAccResourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  billing_account_id = "%s"
  folder_id = "%s"
  region = "%s"
  zone = "%s"
}
`

func TestAccResourceTaikunCloudCredentialGCP(t *testing.T) {
	cloudCredentialName := randomTestName()

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
					os.Getenv("GCP_ZONE"),
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
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "zone", os.Getenv("GCP_ZONE")),
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
  zone = "%s"
  lock = true
}
`

func TestAccResourceTaikunCloudCredentialGCPImportProject(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPImportProjectConfig,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
					os.Getenv("GCP_ZONE"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "config_file", "./gcp.json"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", "true"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "zone", os.Getenv("GCP_ZONE")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialGCPExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
		if err != nil || response.Payload.TotalCountGoogle != 1 {
			return fmt.Errorf("gcp cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialGCPDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_gcp" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCountGoogle != 0 {
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
