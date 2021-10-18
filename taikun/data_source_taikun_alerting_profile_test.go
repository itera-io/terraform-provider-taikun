package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunAlertingProfileConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url  = "https://www.example.org"
  channel = "any"
  type = "Alert"
}

resource "taikun_alerting_profile" "foo" {
  name = "%s"
  reminder = "%s"
  slack_configuration_id = resource.taikun_slack_configuration.foo.id

  is_locked = %t

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
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := 1
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 4
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)
	numberOfIntegrations := 2

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAlertingProfileConfig,
					slackConfigName,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "is_locked", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("data.taikun_alerting_profile.foo", "integration.#", fmt.Sprint(numberOfIntegrations)),
				),
			},
		},
	})
}
