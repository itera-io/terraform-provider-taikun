package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	email := randomString() + "@" + randomString() + ".fr"
	role := "Manager"
	displayName := randomTestName()

	resource.Test(t, resource.TestCase{
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("data.taikun_user.foo", "role", role),
					resource.TestCheckResourceAttr("data.taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("data.taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "user_disabled"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "approved_by_partner"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "is_csm"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("data.taikun_user.foo", "email_notification_enabled"),
				),
			},
		},
	})
}
