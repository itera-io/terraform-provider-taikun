package testing_test

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This test is removed, because we cannot avoid race conditions.
// This data source lists access profiles from the default itera organization
// , but they are constantly created and destroyed by other tests run in parallel anyway.
//
//const testAccDataSourceTaikunAccessProfilesConfig = `
//data "taikun_access_profiles" "all" {
//}`
//
//func TestAccDataSourceTaikunAccessProfiles(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccDataSourceTaikunAccessProfilesConfig,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("data.taikun_access_profiles.all", "id", "all"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.lock"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.name"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ntp_server.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_id"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_name"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ssh_user.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.allowed_host.#"),
//				),
//			},
//		},
//	})
//}

const testAccDataSourceTaikunAccessProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

data "taikun_access_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id
}`

func TestAccDataSourceTaikunAccessProfilesWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAccessProfilesWithFilterConfig, organizationName, organizationFullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_access_profiles.all", "access_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ssh_user.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.allowed_host.#"),
				),
			},
		},
	})
}
