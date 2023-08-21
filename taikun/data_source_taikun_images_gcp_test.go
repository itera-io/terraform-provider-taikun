package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunImagesGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  billing_account_id = "%s"
  folder_id = "%s"
  region = "%s"
}

data "taikun_images_gcp" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_gcp.foo.id
  type = "windows"
}
`

func TestAccDataSourceTaikunImagesGCP(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesGCPConfig,
					cloudCredentialName,
					os.Getenv("GCP_BILLING_ACCOUNT"),
					os.Getenv("GCP_FOLDER_ID"),
					os.Getenv("GCP_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_gcp.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images_gcp.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_gcp.foo", "images.0.id"),
				),
			},
		},
	})
}
