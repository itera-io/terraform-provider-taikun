package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesVsphereConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  hypervisors = [%s]
}

data "taikun_images_vsphere" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_vsphere.foo.id
}
`

func TestAccDataSourceTaikunImagesVsphere(t *testing.T) {
	cloudCredentialName := randomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("VSPHERE_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckVsphere(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesVsphereConfig,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_vsphere.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("data.taikun_images_vsphere.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_images_vsphere.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images_vsphere.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_vsphere.foo", "images.0.id"),
				),
			},
		},
	})
}
