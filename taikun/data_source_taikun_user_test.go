package taikun

import (
	"fmt"
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
	userName := randomTestName()
	email := randomEmail()
	role := "Manager"
	displayName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunUserConfig,
					userName,
					email,
					role,
					displayName,
				),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_user.foo",
					"taikun_user.foo",
				),
			},
		},
	})
}
