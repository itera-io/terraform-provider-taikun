package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunProjectConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id

  monitoring = %t
  expiration_date = "%s"
}

data "taikun_project" "foo" {
  id = resource.taikun_project.foo.id
}
`

func TestAccDataSourceTaikunProject(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableMonitoring,
					expirationDate),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_project.foo",
					"taikun_project.foo",
				),
			},
		},
	})
}
