package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialVsphereConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  hypervisors = [%s]
  lock = %t
}

data "taikun_cloud_credential_vsphere" "foo" {
  id = resource.taikun_cloud_credential_vsphere.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialVsphere(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("VSPHERE_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialVsphereConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_vsphere.foo",
					"taikun_cloud_credential_vsphere.foo",
					map[string]struct{}{
						"password": {},
						"api_host": {},
					},
				),
			},
		},
	})
}
