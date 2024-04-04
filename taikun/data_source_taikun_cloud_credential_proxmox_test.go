package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialProxmoxConfig = `
resource "taikun_cloud_credential_proxmox" "foo" {
  name = "%s"
  hypervisors = [%s]
  lock = %t
}

data "taikun_cloud_credential_proxmox" "foo" {
  id = resource.taikun_cloud_credential_proxmox.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialProxmox(t *testing.T) {
	cloudCredentialName := randomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("PROXMOX_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckProxmox(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_proxmox.foo",
					"taikun_cloud_credential_proxmox.foo",
					map[string]struct{}{
						"client_secret": {},
						"api_host":      {},
					},
				),
			},
		},
	})
}
