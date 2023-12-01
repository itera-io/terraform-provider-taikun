package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

	resource.ParallelTest(t, resource.TestCase{
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
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_showback_credential.foo",
					"taikun_showback_credential.foo",
					map[string]struct{}{
						"password": {},
					},
				),
			},
		},
	})
}
