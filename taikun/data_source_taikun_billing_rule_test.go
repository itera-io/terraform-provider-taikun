package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//TODO We should not use a hardcoded id
const testAccDataSourceBillingRule = `
data "taikun_billing_rule" "foo" {
  id = "162"
}
`

func TestAccDataSourceTaikunBillingRule(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBillingRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_rule.foo", "label.0.label"),
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
