package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunBillingCredentialConfig = `
resource "taikun_billing_credential" "foo" {
  name = "%s"

  prometheus_password = "%s"
  prometheus_url      = "%s"
  prometheus_username = "%s"
}

data "taikun_billing_credential" "foo" {
  id = resource.taikun_billing_credential.foo.id
}
`

func TestAccDataSourceTaikunBillingCredential(t *testing.T) {
	billingCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingCredentialConfig,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}
