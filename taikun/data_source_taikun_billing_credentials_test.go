package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunBillingCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaikunBillingCredentialConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_password"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_username"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingCredentialConfig() string {
	return fmt.Sprintln(`
data "taikun_billing_credentials" "all" {
    organization_id="638"
}`)
}
