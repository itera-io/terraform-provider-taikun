package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunAccessProfilesConfig = `
data "taikun_access_profiles" "all" {
}`

func TestAccDataSourceTaikunAccessProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunAccessProfilesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_access_profiles.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.project.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ssh_user.#"),
				),
			},
		},
	})
}

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
	organizationName := randomTestName()
	organizationFullName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAccessProfilesWithFilterConfig, organizationName, organizationFullName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_access_profiles.all", "access_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.project.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ssh_user.#"),
				),
			},
		},
	})
}
