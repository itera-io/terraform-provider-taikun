package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesZadaraConfig = `
resource "taikun_cloud_credential_zadara" "foo" {
  name = "%s"
}

data "taikun_images_zadara" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_zadara.foo.id
  latest = false
}
`

func TestAccDataSourceTaikunImagesZadara(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesZadaraConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_zadara.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_zadara.foo", "images.0.id"),
				),
			},
		},
	})
}
