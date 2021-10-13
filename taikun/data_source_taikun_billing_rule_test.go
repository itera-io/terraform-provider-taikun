package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunBillingRuleConfig = `
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

data "taikun_billing_rule" "foo" {
  id = resource.taikun_billing_rule.foo.id
}
`

func TestAccDataSourceTaikunBillingRule(t *testing.T) {
	billingCredentialName := randomTestName()
	billingRuleName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingRuleConfig,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					billingRuleName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.0.key"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.0.value"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "type"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "price"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
		},
	})
}
