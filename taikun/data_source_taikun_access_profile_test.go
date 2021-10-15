package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunAccessProfileConfig = `
resource "taikun_access_profile" "foo" {
  name = "%s"
}

data "taikun_access_profile" "foo" {
  id = resource.taikun_access_profile.foo.id
}
`

func TestAccDataSourceTaikunAccessProfile(t *testing.T) {
	accessProfileName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAccessProfileConfig, accessProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "dns_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "ntp_server.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "project.#"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "ssh_user.#"),
				),
			},
		},
	})
}
