package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunSlackConfigurationsConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url = "%s"
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
	slackConfigurationName := utils.RandomTestName()
	url := os.Getenv("SLACK_WEBHOOK") // Slack webhook is checked if valid in new API

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationsConfig, slackConfigurationName, url),
				Check: resource.ComposeAggregateTestCheckFunc(
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
  url = "%s"
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
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	slackConfigurationName := utils.RandomTestName()
	url := os.Getenv("SLACK_WEBHOOK") // Slack webhook is checked if valid in new API

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationsWithFilterConfig,
					organizationName,
					organizationFullName,
					slackConfigurationName,
					url,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
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
