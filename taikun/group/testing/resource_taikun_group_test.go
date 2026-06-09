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

const testAccResourceTaikunGroupConfig = `
resource "taikun_account" "test" {
  name  = "%s"
  email = "%s@example.com"
}

resource "taikun_group" "foo" {
  name        = "%s"
  account_id  = taikun_account.test.id
  claim_value = "%s"
}
`

func TestAccResourceTaikunGroup(t *testing.T) {
	accountName := utils.RandomTestName()
	groupName := utils.RandomTestName()
	claimValue := "test-claim"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunGroupDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunGroupConfig, accountName, accountName, groupName, claimValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunGroupExists(t),
					resource.TestCheckResourceAttrSet("taikun_group.foo", "id"),
					resource.TestCheckResourceAttr("taikun_group.foo", "name", groupName),
					resource.TestCheckResourceAttr("taikun_group.foo", "claim_value", claimValue),
					resource.TestCheckResourceAttrSet("taikun_group.foo", "account_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunGroupExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_group" {
				continue
			}

			id, _ := utils.Atoi32(rs.Primary.ID)
			response, _, err := client.Client.GroupsAPI.GroupsList(t.Context()).Execute()
			if err != nil {
				return fmt.Errorf("error listing groups: %s", err)
			}

			found := false
			for _, g := range response.GetData() {
				if g.GetId() == id {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("group doesn't exist (id = %s)", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTaikunGroupDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_group" {
				continue
			}

			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				id, _ := utils.Atoi32(rs.Primary.ID)
				response, _, err := client.Client.GroupsAPI.GroupsList(t.Context()).Execute()
				if err != nil {
					return retry.NonRetryableError(err)
				}

				for _, g := range response.GetData() {
					if g.GetId() == id {
						return retry.RetryableError(errors.New("group still exists"))
					}
				}
				return nil
			})
			if utils.TimedOut(retryErr) {
				return errors.New("group still exists (timed out)")
			}
			if retryErr != nil {
				return retryErr
			}
		}

		return nil
	}
}
