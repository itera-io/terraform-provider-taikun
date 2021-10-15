package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceOrganizationConfig = `
data "taikun_organization" "foo" {
}
`

func TestAccDataSourceTaikunOrganization(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "discount_rate"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "full_name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "is_read_only"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "projects"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "servers"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "users"),
				),
			},
		},
	})
}
