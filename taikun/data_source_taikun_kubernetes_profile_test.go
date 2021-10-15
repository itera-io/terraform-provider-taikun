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
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunKubernetesProfileConfig, kubernetesProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "bastion_proxy_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profile.foo", "load_balancing_solution"),
				),
			},
		},
	})
}
