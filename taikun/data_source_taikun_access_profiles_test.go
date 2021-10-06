package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO add test with organization_ID set
// with additional check:
// resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "organization_id")

func TestAccDataSourceTaikunAccessProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaikunAccessProfilesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.last_modified"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.last_modified_by"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profiles.all", "access_profiles.0.project.#"),
				),
			},
		},
	})
}

func testAccCheckTaikunAccessProfilesConfig() string {
	return fmt.Sprintln(`
data "taikun_access_profiles" "all" {
  # organization_id="441"
}`)
}
