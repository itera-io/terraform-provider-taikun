package taikun

import (
	"fmt"
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
	cloudCredentialName := randomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("VSPHERE_HYPERVISOR"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckVsphere(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialVsphereConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
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
