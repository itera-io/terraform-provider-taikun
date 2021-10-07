package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//TODO We should not use a hardcoded id
const testAccDataSourceAccessProfile = `
data "taikun_access_profile" "foo" {
  id = "333"
}
`

func TestAccDataSourceTaikunAccessProfile(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessProfile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "last_modified"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "last_modified_by"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "project.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "ssh_user.#"),
				),
			},
		},
	})
}
