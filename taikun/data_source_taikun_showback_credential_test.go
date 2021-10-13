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

data "taikun_showback_credential" "foo" {
  id = resource.taikun_showback_credential.foo.id
}
`

func TestAccDataSourceTaikunShowbackCredential(t *testing.T) {
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
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}
