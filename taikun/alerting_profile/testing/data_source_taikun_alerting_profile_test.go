package testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"
)

const testAccDataSourceTaikunAlertingProfileConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url  = "%s"
  channel = "any"
  type = "Alert"
}

resource "taikun_alerting_profile" "foo" {
  name = "%s"
  reminder = "%s"
  slack_configuration_id = resource.taikun_slack_configuration.foo.id

  lock = %t

  # emails:
  %s

  # webhooks:
  %s

  # integrations
  integration {
    type = "Pagerduty"
    url = "https://www.pagerduty.example"
    token = "secret_token"
  }
  integration {
    type = "MicrosoftTeams"
    url = "https://www.teams.example"
  }
}

data "taikun_alerting_profile" "foo" {
  id = resource.taikun_alerting_profile.foo.id
}
`

func TestAccDataSourceTaikunAlertingProfile(t *testing.T) {
	slackConfigName := utils.RandomTestName()
	slackUrl := os.Getenv("SLACK_WEBHOOK")
	alertingProfileName := utils.RandomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[utils.RandomInt(3)]
	isLocked := utils.RandomBool()
	numberOfEmails := 1
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 4
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAlertingProfileConfig,
					slackConfigName,
					slackUrl,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_alerting_profile.foo",
					"taikun_alerting_profile.foo",
				),
			},
		},
	})
}
