package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunKubernetesProfileConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
}

data "taikun_kubernetes_profile" "foo" {
  id = resource.taikun_kubernetes_profile.foo.id
}
`

func TestAccDataSourceTaikunKubernetesProfile(t *testing.T) {
	kubernetesProfileName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunKubernetesProfileConfig, kubernetesProfileName),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_kubernetes_profile.foo",
					"taikun_kubernetes_profile.foo",
				),
			},
		},
	})
}
