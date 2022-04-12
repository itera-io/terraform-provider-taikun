package taikun

import (
	"fmt"
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
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", false),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", false),
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
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "import_project", true),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", true),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "zone", os.Getenv("GCP_ZONE")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialGCPExists(state *terraform.State) error {
	// FIXME

	return nil
}

func testAccCheckTaikunCloudCredentialGCPDestroy(state *terraform.State) error {
	// FIXME

	return nil
}
