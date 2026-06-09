package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunAccessProfileConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_access_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
}

data "taikun_access_profile" "foo" {
  id = resource.taikun_access_profile.foo.id
}
`

func TestAccDataSourceTaikunAccessProfile(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	accessProfileName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAccessProfileConfig, organizationName, organizationFullName, accessProfileName),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_access_profile.foo",
					"taikun_access_profile.foo",
				),
			},
		},
	})
}
