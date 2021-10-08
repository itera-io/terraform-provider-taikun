package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunBillingRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaikunBillingRulesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.labels.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.labels.0.label"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.labels.0.value"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.labels.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.price"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rules.all", "billing_rules.0.billing_credential_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingRulesConfig() string {
	return fmt.Sprintln(`
data "taikun_billing_rules" "all" {
}`)
}
