package taikun

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTaikunKubernetesProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTaikunKubernetesProfilesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.cni"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.expose_node_port_on_bastion"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.octavia_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.taikun_lb_enabled"),
				),
			},
		},
	})
}

func testAccCheckTaikunKubernetesProfilesConfig() string {
	return `
data "taikun_kubernetes_profiles" "all" {
    #organization_id="638"
}`
}
