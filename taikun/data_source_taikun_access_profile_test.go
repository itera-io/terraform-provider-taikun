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

func TestAccDataSourceAccessProfile(t *testing.T) {

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
				),
			},
		},
	})
}
