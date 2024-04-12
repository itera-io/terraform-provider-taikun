package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunUsersConfig = `
data "taikun_users" "all" {
}`

func TestAccDataSourceTaikunUsers(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunUsersConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_users.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.#"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.user_name"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.role"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.email"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.email_confirmed"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.email_notification_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_csm"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_owner"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_disabled"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_approved_by_partner"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunUsersWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "%s"

  display_name        = "%s"
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_users" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_user.foo
  ]
}`

func TestAccDataSourceTaikunUsersWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	userName := utils.RandomTestName()
	email := utils.RandomEmail()
	role := "User"
	displayName := utils.RandomString()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunUsersWithFilterConfig,
					organizationName,
					organizationFullName,
					userName,
					email,
					role,
					displayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_users.all", "users.0.user_name", userName),
					resource.TestCheckResourceAttr("data.taikun_users.all", "users.0.role", role),
					resource.TestCheckResourceAttr("data.taikun_users.all", "users.0.email", email),
					resource.TestCheckResourceAttr("data.taikun_users.all", "users.0.display_name", displayName),
					resource.TestCheckResourceAttr("data.taikun_users.all", "users.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.#"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.email_confirmed"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.email_notification_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_csm"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_owner"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_disabled"),
					resource.TestCheckResourceAttrSet("data.taikun_users.all", "users.0.is_approved_by_partner"),
				),
			},
		},
	})
}
