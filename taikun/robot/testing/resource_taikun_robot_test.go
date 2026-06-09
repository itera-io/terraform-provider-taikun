package testing

import (
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

func preCheckRobot(t *testing.T) {
	utils_testing.TestAccPreCheck(t)
	t.Skip("POST /api/v1/robot/create requires Administrator privileges (HTTP 403)")
}

const testAccResourceTaikunRobotConfig = `
resource "taikun_robot" "foo" {
  name        = "%s"
  description = "Test robot user"
  expires_at  = "2027-01-01T00:00:00Z"
  scopes      = ["read:project"]
}
`

func TestAccResourceTaikunRobot(t *testing.T) {
	robotName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheckRobot(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunRobotDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunRobotConfig, robotName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRobotExists(t),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "id"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "name", robotName),
					resource.TestCheckResourceAttr("taikun_robot.foo", "description", "Test robot user"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "expires_at", "2027-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "scopes.#", "1"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "scopes.0", "read:project"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "access_key"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "secret_key"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "user_id"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "created_at"),
					resource.TestCheckResourceAttrSet("taikun_robot.foo", "account_id"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "is_active", "true"),
				),
			},
			{
				ResourceName:      "taikun_robot.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"secret_key",
				},
			},
		},
	})
}

const testAccResourceTaikunRobotUpdateConfig = `
resource "taikun_robot" "foo" {
  name        = "%s"
  description = "%s"
  expires_at  = "%s"
  scopes      = [%s]
  ips         = [%s]
}
`

func TestAccResourceTaikunRobotUpdate(t *testing.T) {
	robotName := utils.RandomTestName()
	newRobotName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheckRobot(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunRobotDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunRobotUpdateConfig,
					robotName,
					"Initial description",
					"2027-01-01T00:00:00Z",
					`"read:project"`,
					""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRobotExists(t),
					resource.TestCheckResourceAttr("taikun_robot.foo", "name", robotName),
					resource.TestCheckResourceAttr("taikun_robot.foo", "description", "Initial description"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "scopes.#", "1"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "scopes.0", "read:project"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "ips.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunRobotUpdateConfig,
					newRobotName,
					"Updated description",
					"2027-06-01T00:00:00Z",
					`"read:project", "write:project"`,
					`"192.168.1.0/24"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRobotExists(t),
					resource.TestCheckResourceAttr("taikun_robot.foo", "name", newRobotName),
					resource.TestCheckResourceAttr("taikun_robot.foo", "description", "Updated description"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "expires_at", "2027-06-01T00:00:00Z"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "scopes.#", "2"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "ips.#", "1"),
					resource.TestCheckResourceAttr("taikun_robot.foo", "ips.0", "192.168.1.0/24"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunRobotsConfig = `
resource "taikun_robot" "foo" {
  name        = "%s"
  description = "Data source test robot"
  expires_at  = "2027-01-01T00:00:00Z"
  scopes      = ["read:project"]
}

data "taikun_robots" "all" {
  account_id = taikun_robot.foo.account_id
}
`

func TestAccDataSourceTaikunRobots(t *testing.T) {
	robotName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheckRobot(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunRobotDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunRobotsConfig, robotName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_robots.all", "robots.#"),
				),
			},
		},
	})
}

func testAccCheckTaikunRobotExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_robot" {
				continue
			}

			response, _, err := client.Client.RobotAPI.RobotList(t.Context()).SearchId(rs.Primary.ID).Execute()
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
}

func testAccCheckTaikunRobotDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_robot" {
				continue
			}

			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				response, _, err := client.Client.RobotAPI.RobotList(t.Context()).SearchId(rs.Primary.ID).Execute()
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
}
