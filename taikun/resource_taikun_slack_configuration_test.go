package taikun

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/slack"
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
	url := "https://www.example.org"
	channel := randomString()
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
	url := "https://www.example.org"
	newUrl := "https://www.example.com"
	channel := randomString()
	newChannel := randomString()
	slackConfigType := []string{"Alert", "General"}[rand.Int()%2]
	newSlackConfigType := []string{"Alert", "General"}[rand.Int()%2]

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
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_slack_configuration" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := slack.NewSlackListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Slack.SlackList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("slack configuration doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunSlackConfigurationDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_slack_configuration" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := slack.NewSlackListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.Slack.SlackList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
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
