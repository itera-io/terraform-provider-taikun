package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunShowbackCredentials(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaikunShowbackCredentialConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.prometheus_password"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_credentials.all", "showback_credentials.0.prometheus_username"),
				),
			},
		},
	})
}

func testAccCheckTaikunShowbackCredentialConfig() string {
	return fmt.Sprintln(`
data "taikun_showback_credentials" "all" {
}`)
}
