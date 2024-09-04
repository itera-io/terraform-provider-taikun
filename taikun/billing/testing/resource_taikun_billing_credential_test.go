package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/billing"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunBillingCredentialConfig = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"
  lock       = %t

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}
`

func TestAccResourceTaikunBillingCredential(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig,
					firstName,
					false,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBillingCredentialLock(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, false, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, true, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingCredentialExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)
		foundBillingCredId, err := billing.ResourceTaikunBillingCredentialFind(id, client)
		if err != nil || foundBillingCredId == nil {
			return fmt.Errorf("billing credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBillingCredentialDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			billingCredential, err := billing.ResourceTaikunBillingCredentialFind(id, client)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if billingCredential != nil {
				return retry.RetryableError(errors.New("billing credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("billing credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
