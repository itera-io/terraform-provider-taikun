package taikun

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	url := os.Getenv("SLACK_WEBHOOK") // Slack webhook is checked if valid in new API
	channel := randomString()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationConfig, name, url, channel, slackConfigType),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_slack_configuration.foo",
					"taikun_slack_configuration.foo",
				),
			},
		},
	})
}
