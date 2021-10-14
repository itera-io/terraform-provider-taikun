package taikun

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/alerting_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_alerting_profile", &resource.Sweeper{
		Name: "taikun_alerting_profile",
		F: func(r string) error {
			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion)

			var alertingProfilesList []*models.AlertingProfilesListDto

			for {
				response, err := apiClient.client.AlertingProfiles.AlertingProfilesList(params, apiClient)
				if err != nil {
					return err
				}
				alertingProfilesList = append(alertingProfilesList, response.Payload.Data...)
				if len(alertingProfilesList) == int(response.Payload.TotalCount) {
					break
				}
				offset := int32(len(alertingProfilesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range alertingProfilesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					body := models.DeleteAlertingProfilesCommand{ID: e.ID}
					params := alerting_profiles.NewAlertingProfilesDeleteParams().WithV(ApiVersion).WithBody(&body)
					_, _, err = apiClient.client.AlertingProfiles.AlertingProfilesDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

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
  full_name = "Foo"
  discount_rate = 42
}

resource "taikun_slack_configuration" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  url  = "https://www.example.org"
  channel = "any"
  type = "Alert"
}

resource "taikun_alerting_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

  name = "%s"
  reminder = "%s"
  slack_configuration_id = resource.taikun_slack_configuration.foo.id

  is_locked = %t

  # emails:
  %s

  # webhooks:
  %s
}
`

func TestAccResourceTaikunAlertingProfileConfiguration(t *testing.T) {
	organizationName := randomTestName()
	slackConfigName := randomTestName()
	alertingProfileName := randomTestName()
	reminder := []string{"HalfHour", "Hourly", "Daily"}[randomInt(3)]
	isLocked := randomBool()
	numberOfEmails := randomInt(10)
	emails := testAccResourceTaikunAlertingProfileRandomEmails(numberOfEmails)
	numberOfWebhooks := randomInt(10)
	webhooks := testAccResourceTaikunAlertingProfileRandomWebhooks(numberOfWebhooks)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAlertingProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAlertingProfileConfig,
					organizationName,
					slackConfigName,
					alertingProfileName,
					reminder,
					isLocked,
					emails,
					webhooks),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAlertingProfileExists,
					resource.TestCheckResourceAttrSet("taikun_alerting_profile.foo", "id"),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "slack_configuration_name", slackConfigName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "reminder", reminder),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "is_locked", fmt.Sprint(isLocked)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "emails.#", fmt.Sprint(numberOfEmails)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "webhook.#", fmt.Sprint(numberOfWebhooks)),
					resource.TestCheckResourceAttr("taikun_alerting_profile.foo", "organization_name", organizationName),
				),
			},
		},
	})

}

func testAccCheckTaikunAlertingProfileExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_alerting_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.AlertingProfiles.AlertingProfilesList(params, apiClient)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("alerting profile with ID %d doesn't exist", id)
		}
	}

	return nil
}

func testAccCheckTaikunAlertingProfileDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_alerting_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := alerting_profiles.NewAlertingProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.AlertingProfiles.AlertingProfilesList(params, apiClient)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("alerting profile with ID %d still exists", id)
		}
	}

	return nil
}
