package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
  lock = %t
}

data "taikun_cloud_credential_openstack" "foo" {
  id = resource.taikun_cloud_credential_openstack.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialOpenStack(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
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
