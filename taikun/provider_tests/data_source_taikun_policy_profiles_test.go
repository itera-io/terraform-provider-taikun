package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunPolicyProfilesConfig = `
resource "taikun_policy_profile" "foo" {
  name = "%s"

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t
}

data "taikun_policy_profiles" "all" {
  depends_on = [
    taikun_policy_profile.foo
  ]
}`

func TestAccDataSourceTaikunPolicyProfiles(t *testing.T) {
	PolicyProfileName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceTaikunPolicyProfilesConfig,
					PolicyProfileName,
					false,
					false,
					false,
					false,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_policy_profiles.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.forbid_node_port"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.forbid_http_ingress"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.require_probe"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.unique_ingress"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.unique_service_selector"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.organization_name"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunPolicyProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_policy_profile" "foo" {
  name = "%s"

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t

  organization_id = resource.taikun_organization.foo.id
}

data "taikun_policy_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id
  depends_on = [
    taikun_policy_profile.foo
  ]
}`

func TestAccDataSourceTaikunPolicyProfilesWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	PolicyProfileName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccDataSourceTaikunPolicyProfilesWithFilterConfig,
					organizationName,
					organizationFullName,
					PolicyProfileName,
					false,
					false,
					false,
					false,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_policy_profiles.all", "policy_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.forbid_node_port"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.forbid_http_ingress"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.require_probe"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.unique_ingress"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.unique_service_selector"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_policy_profiles.all", "policy_profiles.0.organization_name"),
				),
			},
		},
	})
}
