package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesProxmoxConfig = `
resource "taikun_cloud_credential_proxmox" "foo" {
  name = "%s"
  hypervisors = [%s]
}

data "taikun_images_proxmox" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_proxmox.foo.id
}
`

func TestAccDataSourceTaikunImagesProxmox(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("PROXMOX_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_proxmox.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("data.taikun_images_proxmox.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_images_proxmox.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images_proxmox.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_proxmox.foo", "images.0.id"),
				),
			},
		},
	})
}
