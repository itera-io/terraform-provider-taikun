package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
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
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("PROXMOX_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
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
