package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

//TODO We should not use a hardcoded id
const testAccDataSourceAccessProfile = `
data "taikun_access_profile" "foo" {
  id = 333
}
`

func TestAccDataSourceAccessProfile(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessProfile,
				Check: resource.ComposeTestCheckFunc(
					//func(state *terraform.State) error {
					//	fmt.Print(state.String())
					//	return nil
					//},
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_access_profile.foo", "organization_name"),
				),
			},
		},
	})
}
