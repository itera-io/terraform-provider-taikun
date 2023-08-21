package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunImagesAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_images_azure" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  publisher = "Canonical"
  offer = "UbuntuServer"
  sku = "19.04"
  latest = true
}
`

func TestAccDataSourceTaikunImagesAzure(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesAzureConfig,
					cloudCredentialName,
					os.Getenv("ARM_LOCATION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_azure.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images_azure.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_azure.foo", "images.0.id"),
				),
			},
		},
	})
}
