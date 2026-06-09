package testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"
)

const testAccDataSourceTaikunAlertingProfileConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_alerting_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  reminder = "%s"

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
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
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
					organizationName,
					organizationFullName,
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

// --- Slack integration tests (disabled) ---
// Slack is not available on the current platform: POST /api/v1/slack/create returns
// HTTP 403 (Forbidden Access Exception). Active test above omits slack_configuration_id.
// When Slack is re-enabled, uncomment the block below, add `import "os"`, and set
// SLACK_WEBHOOK to a valid webhook URL.
//
// TestAccDataSourceTaikunAlertingProfileSlack is the Slack-linked variant of the
// data source test (resource + data source state parity with slack_configuration_id).
//
// const testAccDataSourceTaikunAlertingProfileSlackConfig = `
// resource "taikun_organization" "foo" {
//   name = "%s"
//   full_name = "%s"
//   discount_rate = 42
// }
//
// resource "taikun_slack_configuration" "foo" {
//   organization_id = resource.taikun_organization.foo.id
//
//   name = "%s"
//   url  = "%s"
//   channel = "any"
//   type = "Alert"
// }
//
// resource "taikun_alerting_profile" "foo" {
//   organization_id = resource.taikun_organization.foo.id
//
//   name = "%s"
//   reminder = "%s"
//   slack_configuration_id = resource.taikun_slack_configuration.foo.id
//
//   lock = %t
//
//   # emails:
//   %s
//
//   # webhooks:
//   %s
//
//   # integrations
//   integration {
//     type = "Pagerduty"
//     url = "https://www.pagerduty.example"
//     token = "secret_token"
//   }
//   integration {
//     type = "MicrosoftTeams"
//     url = "https://www.teams.example"
//   }
// }
//
// data "taikun_alerting_profile" "foo" {
//   id = resource.taikun_alerting_profile.foo.id
// }
// `
//
// func TestAccDataSourceTaikunAlertingProfileSlack(t *testing.T) {
// 	organizationName := utils.RandomTestName()
// 	organizationFullName := utils.RandomTestName()
// 	slackConfigName := utils.RandomTestName()
// 	slackUrl := os.Getenv("SLACK_WEBHOOK")
// 	alertingProfileName := utils.RandomTestName()
// 	reminder := []string{"HalfHour", "Hourly", "Daily"}[utils.RandomInt(3)]
// 	isLocked := utils.RandomBool()
// 	numberOfEmails := 1
// 	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
// 	numberOfWebhooks := 4
// 	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
// 		ProviderFactories: utils_testing.TestAccProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: fmt.Sprintf(testAccDataSourceTaikunAlertingProfileSlackConfig,
// 					organizationName,
// 					organizationFullName,
// 					slackConfigName,
// 					slackUrl,
// 					alertingProfileName,
// 					reminder,
// 					isLocked,
// 					emails,
// 					webhooks),
// 				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
// 					"data.taikun_alerting_profile.foo",
// 					"taikun_alerting_profile.foo",
// 				),
// 			},
// 		},
// 	})
// }
