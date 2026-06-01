package testing

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccResourceTaikunRobotConfig = `
resource "taikun_account" "test" {
  name  = "%s"
  email = "%s@example.com"
}

resource "taikun_robot" "foo" {
  name        = "%s"
  account_id  = taikun_account.test.id
  description = "Test robot user"
  expires_at  = "2027-01-01T00:00:00Z"
  scopes      = ["read:project"]
}
`

func TestAccResourceTaikunRobot(t *testing.T) {
	accountName := utils.RandomTestName()
	robotName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunRobotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunRobotConfig, accountName, accountName, robotName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRobotExists,
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "id"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "name", robotName),
					resource.TestCheckResourceAttr("taikun_robot.foo", "description", "Test robot user"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "access_key"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "secret_key"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "user_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunRobotExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_robot" {
			continue
		}

		response, _, err := client.Client.RobotAPI.RobotList(context.TODO()).SearchId(rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("error looking up robot: %s", err)
		}

		found := false
		for _, r := range response.GetData() {
			if r.GetUserId() == rs.Primary.ID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("robot doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunRobotDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_robot" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			response, _, err := client.Client.RobotAPI.RobotList(context.TODO()).SearchId(rs.Primary.ID).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}

			for _, r := range response.GetData() {
				if r.GetUserId() == rs.Primary.ID {
					return retry.RetryableError(errors.New("robot still exists"))
				}
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("robot still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
