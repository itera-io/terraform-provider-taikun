package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	kubernetesProfileName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunKubernetesProfileConfig, kubernetesProfileName),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_kubernetes_profile.foo",
					"taikun_kubernetes_profile.foo",
				),
			},
		},
	})
}
