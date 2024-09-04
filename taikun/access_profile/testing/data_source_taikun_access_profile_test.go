package testing_test

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunAccessProfileConfig = `
resource "taikun_access_profile" "foo" {
  name = "%s"
}

data "taikun_access_profile" "foo" {
  id = resource.taikun_access_profile.foo.id
}
`

func TestAccDataSourceTaikunAccessProfile(t *testing.T) {
	accessProfileName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAccessProfileConfig, accessProfileName),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_access_profile.foo",
					"taikun_access_profile.foo",
				),
			},
		},
	})
}
