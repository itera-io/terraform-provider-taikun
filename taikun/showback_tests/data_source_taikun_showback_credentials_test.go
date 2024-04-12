package showback_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunShowbackCredentialsConfig = `
resource "taikun_showback_credential" "foo" {
  name            = "%s"

  password = "%s"
  url = "%s"
  username = "%s"
}

data "taikun_showback_credentials" "all" {
   depends_on = [
    taikun_showback_credential.foo
  ]
}`

func TestAccDataSourceTaikunShowbackCredentials(t *testing.T) {
	showbackCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackCredentialsConfig,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.url"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.username"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunShowbackCredentialsWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_showback_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  password = "%s"
  url = "%s"
  username = "%s"
}

data "taikun_showback_credentials" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_showback_credential.foo
  ]
}
`

func TestAccDataSourceTaikunShowbackCredentialsWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	showbackCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackCredentialsWithFilterConfig,
					organizationName,
					organizationFullName,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("data.taikun_showback_credentials.all", "showback_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.url"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.username"),
				),
			},
		},
	})
}
