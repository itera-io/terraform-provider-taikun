package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"

  lock       = %t
}

data "taikun_cloud_credential_openstack" "foo" {
  id = resource.taikun_cloud_credential_openstack.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_openstack.foo",
					"taikun_cloud_credential_openstack.foo",
					map[string]struct{}{
						"password": {},
						"url":      {},
					},
				),
			},
		},
	})
}
