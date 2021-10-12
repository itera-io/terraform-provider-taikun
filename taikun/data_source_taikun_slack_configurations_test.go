package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceSlackConfigurations = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url = "https://www.example.org"
  channel = "any"
  type = "General"
}

data "taikun_slack_configurations" "all" {
}
`

func TestAccDataSourceTaikunSlackConfigurations(t *testing.T) {
	slackConfigurationName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceSlackConfigurations, slackConfigurationName),
				Check: resource.ComposeTestCheckFunc(
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
