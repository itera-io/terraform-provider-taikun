package taikun

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunSlackConfigurationConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url = "%s"
  channel = "%s"
  type = "%s"
}

data "taikun_slack_configuration" "foo" {
  id = resource.taikun_slack_configuration.foo.id
}
`

func TestAccDataSourceTaikunSlackConfiguration(t *testing.T) {
	name := randomTestName()
	url := "https://www.example.org"
	channel := randomString()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationConfig, name, url, channel, slackConfigType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_slack_configuration.foo", "name", name),
					resource.TestCheckResourceAttr("data.taikun_slack_configuration.foo", "url", url),
					resource.TestCheckResourceAttr("data.taikun_slack_configuration.foo", "channel", channel),
					resource.TestCheckResourceAttr("data.taikun_slack_configuration.foo", "type", slackConfigType),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configuration.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configuration.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_slack_configuration.foo", "organization_name"),
				),
			},
		},
	})
}
