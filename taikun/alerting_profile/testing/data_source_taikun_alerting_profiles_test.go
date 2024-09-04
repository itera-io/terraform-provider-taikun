package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This test is removed, because we cannot avoid race conditions.
// This data source lists access profiles from the default itera organization
// , but they are constantly created and destroyed by other tests run in parallel anyway.
//
//const testAccDataSourceTaikunAlertingProfilesConfig = `
//data "taikun_alerting_profiles" "all" {
//}
//`
//
//func TestAccDataSourceTaikunAlertingProfiles(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccDataSourceTaikunAlertingProfilesConfig,
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("data.taikun_alerting_profiles.all", "id", "all"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.emails.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.id"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.integration.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.lock"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.name"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_id"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_name"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.reminder"),
//					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.webhook.#"),
//				),
//			},
//		},
//	})
//}

const testAccDataSourceTaikunAlertingProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

data "taikun_alerting_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id
}
`

func TestAccDataSourceTaikunAlertingProfilesWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAlertingProfilesWithFilterConfig, organizationName, organizationFullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.emails.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.reminder"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.webhook.#"),
				),
			},
		},
	})
}
