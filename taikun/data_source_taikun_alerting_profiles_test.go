package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunAlertingProfilesConfig = `
data "taikun_alerting_profiles" "all" {
}
`

func TestAccDataSourceTaikunAlertingProfiles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaikunAlertingProfilesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_alerting_profiles.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.emails.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.reminder"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.webhook.#"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunAlertingProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "foo"
  discount_rate = 42
}

data "taikun_alerting_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id
}
`

func TestAccDataSourceTaikunAlertingProfilesWithFilter(t *testing.T) {
	organizationName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAlertingProfilesWithFilterConfig, organizationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.emails.#"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.reminder"),
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profiles.all", "alerting_profiles.0.webhook.#"),
				),
			},
		},
	})
}
