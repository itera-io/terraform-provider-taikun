package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunOPAProfileConfig = `
resource "taikun_opa_profile" "foo" {
  name = "%s"

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t
}

data "taikun_opa_profile" "foo" {
  id = resource.taikun_opa_profile.foo.id
}
`

func TestAccDataSourceTaikunOPAProfile(t *testing.T) {
	OPAProfileName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceTaikunOPAProfileConfig,
					OPAProfileName,
					false,
					false,
					false,
					false,
					false,
				),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_opa_profile.foo",
					"taikun_opa_profile.foo",
				),
			},
		},
	})
}
