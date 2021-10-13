package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunKubernetesProfilesConfig = `
data "taikun_kubernetes_profiles" "all" {
}`

func TestAccDataSourceTaikunKubernetesProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunKubernetesProfilesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_kubernetes_profiles.all", "id", "all"),
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
