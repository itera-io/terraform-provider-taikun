package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunUserConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "%s"

  display_name = "%s"
}

data "taikun_user" "foo" {
  id = resource.taikun_user.foo.id
}
`

func TestAccDataSourceTaikunUser(t *testing.T) {
	userName := utils.RandomTestName()
	email := utils.RandomEmail()
	role := "Manager"
	displayName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunUserConfig,
					userName,
					email,
					role,
					displayName,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_user.foo",
					"taikun_user.foo",
				),
			},
		},
	})
}
