package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	slackConfigName := randomTestName()
	slackUrl := os.Getenv("SLACK_WEBHOOK")
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := 1
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 4
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_alerting_profile.foo",
					"taikun_alerting_profile.foo",
				),
			},
		},
	})
}
