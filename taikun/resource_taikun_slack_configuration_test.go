package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunSlackConfigurationConfig = `
resource "taikun_slack_configuration" "foo" {
  name = "%s"
  url  = "%s"
  channel = "%s"
  type = "%s"
}
`

func TestAccResourceTaikunSlackConfiguration(t *testing.T) {
	name := randomTestName()
	//url := "https://www.example.org"
	url := os.Getenv("SLACK_WEBHOOK")
	channel := randomTestName()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunSlackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunSlackConfigurationConfig, name, url, channel, slackConfigType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunSlackConfigurationExists,
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "channel", channel),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "id"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "name", name),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "type", slackConfigType),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "url", url),
				),
			},
			{
				ResourceName:      "taikun_slack_configuration.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunSlackConfigurationModify(t *testing.T) {
	name := randomTestName()
	newName := randomTestName()
	url := os.Getenv("SLACK_WEBHOOK")    // Slack webhook is checked if valid in new API
	newUrl := os.Getenv("SLACK_WEBHOOK") // I do not have a second slack WEBHOOK URL, #TODO room for improvement
	channel := randomTestName()
	newChannel := randomTestName()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]
	newSlackConfigType := []string{"Alert", "General"}[rand.Int()%2]

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunSlackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunSlackConfigurationConfig, name, url, channel, slackConfigType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunSlackConfigurationExists,
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "channel", channel),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "id"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "name", name),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "type", slackConfigType),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "url", url),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunSlackConfigurationConfig, newName, newUrl, newChannel, newSlackConfigType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunSlackConfigurationExists,
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "channel", newChannel),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "id"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "name", newName),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_slack_configuration.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "type", newSlackConfigType),
					resource.TestCheckResourceAttr("taikun_slack_configuration.foo", "url", newUrl),
				),
			},
		},
	})
}

func testAccCheckTaikunSlackConfigurationExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_slack_configuration" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := client.Client.SlackAPI.SlackList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("slack configuration doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunSlackConfigurationDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_slack_configuration" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := client.Client.SlackAPI.SlackList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return resource.RetryableError(errors.New("slack configuration still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("slack configuration still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
