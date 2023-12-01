package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunOrganizationsConfig = `
data "taikun_organizations" "all" {
}
`

func TestAccDataSourceTaikunOrganizations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunOrganizationsConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_organizations.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.#"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.discount_rate"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.full_name"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.projects"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.servers"),
					resource.TestCheckResourceAttrSet("data.taikun_organizations.all", "organizations.0.users"),
				),
			},
		},
	})
}
