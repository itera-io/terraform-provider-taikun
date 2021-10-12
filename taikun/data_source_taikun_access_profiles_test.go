package taikun

import (
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
