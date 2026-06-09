package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunBillingRuleConfig = `
resource "taikun_organization" "foo" {
  name          = "%s"
  full_name     = "%s"
  discount_rate = 42
}

resource "taikun_billing_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  prometheus_password = "%s"
  prometheus_url      = "%s"
  prometheus_username = "%s"
}

resource "taikun_billing_rule" "foo" {
  name                  = "%s"
  metric_name           = "coredns_forward_request_duration_seconds"
  price                 = 1
  type                  = "Sum"
  billing_credential_id = resource.taikun_billing_credential.foo.id
  label {
    key   = "key"
    value = "value"
  }
}

data "taikun_billing_rule" "foo" {
  id = resource.taikun_billing_rule.foo.id
}
`

// TestAccDataSourceTaikunBillingRule retrieves a single billing rule by ID.
// Skipped because Managing Operation/Billing Credentials requires a Partner role on the Taikun API,
// which is not possessed by the standard robot account credentials in dev-env.sh.
func TestAccDataSourceTaikunBillingRule(t *testing.T) {
	t.Skip("requires Partner API role privileges; re-enable when Partner credentials are available in dev-env.sh")

	orgName := utils.RandomTestName()
	orgFullName := utils.RandomTestName()
	billingCredentialName := utils.RandomTestName()
	billingRuleName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingRuleConfig,
					orgName,
					orgFullName,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					billingRuleName,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_billing_rule.foo",
					"taikun_billing_rule.foo",
				),
			},
		},
	})
}
