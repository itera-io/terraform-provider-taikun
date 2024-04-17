package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunBillingRulesConfig = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}

resource "taikun_billing_rule" "foo" {
  name            = "%s"
  metric_name     =  "coredns_forward_request_duration_seconds"
  price = 1
  type = "Sum"
  billing_credential_id = resource.taikun_billing_credential.foo.id
  label {
    key = "key"
    value = "value"
  }
}

data "taikun_billing_rules" "all" {
   depends_on = [
    taikun_billing_rule.foo
  ]
}`

func TestAccDataSourceTaikunBillingRules(t *testing.T) {
	billingCredentialName := utils.RandomTestName()
	billingRuleName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingRulesConfig,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					billingRuleName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.label.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.label.0.key"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.label.0.value"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.label.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.price"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.billing_credential_id"),
				),
			},
		},
	})
}
