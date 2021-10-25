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
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_access_profile.foo",
					"taikun_access_profile.foo",
				),
			},
		},
	})
}
