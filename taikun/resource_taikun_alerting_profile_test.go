package taikun

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccResourceTaikunAlertingProfileRandomEmails(n int) string {
	if n == 0 {
		return "emails = []"
	}
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(fmt.Sprintf("emails = [ \"%s\"", randomEmail()))
	for i := 0; i < n-1; i++ {
		stringBuilder.WriteString(fmt.Sprintf(", \"%s\"", randomEmail()))
	}
	stringBuilder.WriteString(" ]")
	return stringBuilder.String()
}

func testAccResourceTaikunAlertingProfileRandomWebhook() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(fmt.Sprintf(`
webhook {
  url = "%s"
`, randomURL()))
	for i := 0; i < rand.Int()%10; i++ {
		stringBuilder.WriteString(fmt.Sprintf(`
  header {
    key = "%s"
    value = "%s"
  }
`, randomString(), randomString()))
	}
	stringBuilder.WriteString("}\n")
	return stringBuilder.String()
}

func testAccResourceTaikunAlertingProfileRandomWebhooks(n int) string {
	stringBuilder := strings.Builder{}
	for i := 0; i < n; i++ {
		stringBuilder.WriteString(testAccResourceTaikunAlertingProfileRandomWebhook())
	}
	return stringBuilder.String()
}

const testAccResourceTaikunAlertingProfileConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_slack_configuration" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  url  = "%s"
  channel = "any"
  type = "Alert"
}

resource "taikun_alerting_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  reminder = "%s"
  slack_configuration_id = resource.taikun_slack_configuration.foo.id

  lock = %t

  # emails:
  %s

  # webhooks:
  %s

  # integrations:
  %s
}
`

func TestAccResourceTaikunAlertingProfile(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	slackConfigName := randomTestName()
	slackUrl := os.Getenv("SLACK_WEBHOOK")
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := 5
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 3
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)
	integrations := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAlertingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					organizationFullName,
					slackConfigName,
					slackUrl,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks,
					integrations),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "lock", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
				),
			},
			{
				ResourceName:      "taikun_alerting_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunAlertingProfileModify(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	slackConfigName := randomTestName()
	slackUrl := os.Getenv("SLACK_WEBHOOK")
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := 5
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 3
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)
	integrations := ""
	newAlertingProfileName := randomTestName()
	newReminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	newIsLocked := randomBool()
	newNumberOfEmails := 2
	newEmails := testAccResourceTaikunAlertingProfileRandomEmails(newNumberOfEmails)
	newNumberOfWebhooks := 4
	newWebhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(newNumberOfWebhooks)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAlertingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					organizationFullName,
					slackConfigName,
					slackUrl,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks,
					integrations),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "lock", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					organizationFullName,
					slackConfigName,
					slackUrl,
					newAlertingProfileName,
					newReminder,
					newIsLocked,
					newEmails,
					newWebhooks,
					integrations),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", newAlertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", newReminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "lock", fmt.Sprint(newIsLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(newNumberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(newNumberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
				),
			},
		},
	})
}

func TestAccResourceTaikunAlertingProfileModifyIntegrations(t *testing.T) {
	//fmt.Println("Random URLs")
	//for i := 0; i < 20; i++ {
	//	fmt.Println(randomURL())
	//}
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	slackConfigName := randomTestName()
	slackUrl := os.Getenv("SLACK_WEBHOOK")
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := 5
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := 3
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)
	integrations := `
integration {
  type = "Opsgenie"
  url = "https://www.opsgenie.example"
  token = "secret_token"
}`
	numberOfIntegrations := 1
	newIntegrations := `
integration {
  type = "Opsgenie"
  url = "https://www.opsgenie.example"
  token = "secret_token"
}
integration {
  type = "Pagerduty"
  url = "https://www.pagerduty.example"
  token = "secret_token"
}
integration {
  type = "MicrosoftTeams"
  url = "https://www.teams.example"
}
integration {
  type = "Splunk"
  url = "https://www.splunk.example"
  token = "secret_token"
}`
	newNumberOfIntegrations := 4

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAlertingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					organizationFullName,
					slackConfigName,
					slackUrl,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks,
					integrations),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "lock", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "integration.#", fmt.Sprint(numberOfIntegrations)),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					organizationFullName,
					slackConfigName,
					slackUrl,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks,
					newIntegrations),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "lock", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "integration.#", fmt.Sprint(newNumberOfIntegrations)),
				),
			},
		},
	})
}

func testAccCheckTaikunAlertingProfileExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_alerting_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := apiClient.Client.AlertingProfilesAPI.AlertingprofilesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("alerting profile with ID %d doesn't exist", id)
		}
	}

	return nil
}

func testAccCheckTaikunAlertingProfileDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_alerting_profile" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := apiClient.Client.AlertingProfilesAPI.AlertingprofilesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("alerting profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("alerting profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
