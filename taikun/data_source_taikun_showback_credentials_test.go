package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunShowbackCredentialConfig = `
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
	showbackCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.password"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.url"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.username"),
				),
			},
		},
	})
}
