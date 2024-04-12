package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_images_azure" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  publisher = "Canonical"
  offer = "0001-com-ubuntu-server-jammy"
  sku = "22_04-lts"
  latest = true
}
`

func TestAccDataSourceTaikunImagesAzure(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
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
