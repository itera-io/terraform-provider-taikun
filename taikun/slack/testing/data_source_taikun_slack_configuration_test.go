package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
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
	name := utils.RandomTestName()
	url := os.Getenv("SLACK_WEBHOOK") // Slack webhook is checked if valid in new API
	channel := utils.RandomString()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunSlackConfigurationConfig, name, url, channel, slackConfigType),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_slack_configuration.foo",
					"taikun_slack_configuration.foo",
				),
			},
		},
	})
}
