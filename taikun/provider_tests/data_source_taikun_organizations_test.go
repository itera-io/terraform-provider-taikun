package provider_tests

import (
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunOrganizationsConfig = `
data "taikun_organizations" "all" {
}
`

func TestAccDataSourceTaikunOrganizations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
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
