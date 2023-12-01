package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunPolicyProfileConfig = `
resource "taikun_policy_profile" "foo" {
  name = "%s"

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t
}

data "taikun_policy_profile" "foo" {
  id = resource.taikun_policy_profile.foo.id
}
`

func TestAccDataSourceTaikunPolicyProfile(t *testing.T) {
	PolicyProfileName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceTaikunPolicyProfileConfig,
					PolicyProfileName,
					false,
					false,
					false,
					false,
					false,
				),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_policy_profile.foo",
					"taikun_policy_profile.foo",
				),
			},
		},
	})
}
