package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunSlackConfigurationsConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url = "https://www.example.org"
  channel = "any"
  type = "General"
}

data "taikun_slack_configurations" "all" {
  depends_on = [
    taikun_slack_configuration.foo
  ]
}
`

func TestAccDataSourceTaikunSlackConfigurations(t *testing.T) {
	slackConfigurationName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationsConfig, slackConfigurationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_slack_configurations.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.#"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.channel"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.url"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunSlackConfigurationsWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_slack_configuration" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  url = "https://www.example.org"
  channel = "any"
  type = "General"
}

data "taikun_slack_configurations" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_slack_configuration.foo
  ]
}
`

func TestAccDataSourceTaikunSlackConfigurationsWithFilter(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	slackConfigurationName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationsWithFilterConfig,
					organizationName,
					organizationFullName,
					slackConfigurationName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_slack_configurations.all", "slack_configurations.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.#"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.channel"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configurations.all", "slack_configurations.0.url"),
				),
			},
		},
	})
}
