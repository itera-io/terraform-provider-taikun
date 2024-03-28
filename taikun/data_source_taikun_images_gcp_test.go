package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  lock = true
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
