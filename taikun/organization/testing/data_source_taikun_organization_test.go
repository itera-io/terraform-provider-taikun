package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceOrganizationConfig = `
data "taikun_organization" "foo" {
}
`

func TestAccDataSourceTaikunOrganization(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "full_name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "projects"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "servers"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "cloud_credentials"),
				),
			},
		},
	})
}

const testAccDataSourceOrganizationNewConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
}

data "taikun_organization" "foo" {
  id = resource.taikun_organization.foo.id
}
`

func TestAccDataSourceTaikunOrganizationNew(t *testing.T) {
	name := utils.RandomTestName()
	fullName := utils.RandomString()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceOrganizationNewConfig,
					name,
					fullName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "id"),
					resource.TestCheckResourceAttr("data.taikun_organization.foo", "name", fmt.Sprint(name)),
					resource.TestCheckResourceAttr("data.taikun_organization.foo", "full_name", fmt.Sprint(fullName)),
				),
			},
		},
	})
}
