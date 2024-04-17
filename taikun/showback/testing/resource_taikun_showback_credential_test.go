package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunShowbackCredentialConfig = `
resource "taikun_showback_credential" "foo" {
  name            = "%s"
  lock       = %t

  password = "%s"
  url = "%s"
  username = "%s"
}
`

func TestAccResourceTaikunShowbackCredential(t *testing.T) {
	showbackCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					false,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}

func TestAccResourceTaikunShowbackCredentialLock(t *testing.T) {
	showbackCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					false,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					true,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}

func testAccCheckTaikunShowbackCredentialExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("showback credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunShowbackCredentialDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.ShowbackClient.ShowbackCredentialsAPI.ShowbackcredentialsList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("showback credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("showback credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
