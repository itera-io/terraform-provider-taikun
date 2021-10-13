package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunKubernetesProfiles(t *testing.T) {
	kubernetesProfileName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckTaikunKubernetesProfilesConfig(), kubernetesProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.bastion_proxy_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.cni"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.load_balancing_solution"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_name"),
				),
			},
		},
	})
}

func testAccCheckTaikunKubernetesProfilesConfig() string {
	return `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
}

data "taikun_kubernetes_profiles" "all" {
   depends_on = [
    taikun_kubernetes_profile.foo
  ]
}`
}
