package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunKubernetesProfilesConfig = `
data "taikun_kubernetes_profiles" "all" {
}`

func TestAccDataSourceTaikunKubernetesProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunKubernetesProfilesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_kubernetes_profiles.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.bastion_proxy"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.cni"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.load_balancing_solution"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.schedule_on_master"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunKubernetesProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

data "taikun_kubernetes_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id
}`

func TestAccDataSourceTaikunKubernetesProfilesWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunKubernetesProfilesWithFilterConfig, organizationName, organizationFullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.bastion_proxy"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.cni"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.load_balancing_solution"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_kubernetes_profiles.all", "kubernetes_profiles.0.schedule_on_master"),
				),
			},
		},
	})
}
