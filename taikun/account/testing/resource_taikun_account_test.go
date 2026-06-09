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

const testAccResourceTaikunAccountConfig = `
resource "taikun_account" "foo" {
  name  = "%s"
  email = "%s"
}
`

func TestAccResourceTaikunAccount(t *testing.T) {
	// Account create/delete requires Administrator privileges on the Taikun API
	// (POST /api/v1/accounts/create returns HTTP 403 with standard acceptance-test credentials).
	t.Skip("requires Administrator API privileges; re-enable when admin credentials are available in dev-env.sh")

	name := utils.RandomTestName()
	email := fmt.Sprintf("%s@example.com", name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccountDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccountConfig, name, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccountExists(t),
					resource.TestCheckResourceAttrSet("taikun_account.foo", "id"),
					resource.TestCheckResourceAttr("taikun_account.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_account.foo", "email", email),
				),
			},
		},
	})
}

func testAccCheckTaikunAccountExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_account" {
				continue
			}

			id, _ := utils.Atoi32(rs.Primary.ID)
			_, _, err := client.Client.AccountsAPI.AccountsDetails(t.Context(), id).Execute()
			if err != nil {
				return fmt.Errorf("account doesn't exist (id = %s)", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTaikunAccountDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_account" {
				continue
			}

			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				id, _ := utils.Atoi32(rs.Primary.ID)
				_, res, err := client.Client.AccountsAPI.AccountsDetails(t.Context(), id).Execute()
				if err != nil {
					if res != nil && res.StatusCode == 404 {
						return nil
					}
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(errors.New("account still exists"))
			})
			if utils.TimedOut(retryErr) {
				return errors.New("account still exists (timed out)")
			}
			if retryErr != nil {
				return retryErr
			}
		}

		return nil
	}
}
