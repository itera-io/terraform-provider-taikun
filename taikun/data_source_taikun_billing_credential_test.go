package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//TODO We should not use a hardcoded id
const testAccDataSourceBillingCredential = `
data "taikun_billing_credential" "foo" {
  id = "89"
}
`

func TestAccDataSourceTaikunBillingCredential(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBillingCredential,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}
