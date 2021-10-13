package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunAccessProfiles(t *testing.T) {
	accessProfileName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckTaikunAccessProfilesConfig(), accessProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.is_locked"),
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

func testAccCheckTaikunAccessProfilesConfig() string {
	return `
resource "taikun_access_profile" "foo" {
  name = "%s"
}

data "taikun_access_profiles" "all" {
   depends_on = [
    taikun_access_profile.foo
  ]
}`
}
