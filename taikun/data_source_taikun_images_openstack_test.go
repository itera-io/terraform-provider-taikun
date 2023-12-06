package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

data "taikun_images_openstack" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
`

func TestAccDataSourceTaikunImagesOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesOpenStackConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_openstack.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images_openstack.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_openstack.foo", "images.0.id"),
				),
			},
		},
	})
}
