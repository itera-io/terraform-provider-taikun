package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
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
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, false, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
			{
				ResourceName:      "taikun_billing_credential.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunBillingCredentialLock(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, false, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
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
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingCredentialExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		resource, err := resourceTaikunBillingCredentialFind(id, client)
		if err != nil || resource == nil {
			return fmt.Errorf("billing credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBillingCredentialDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			billingCredential, err := resourceTaikunBillingCredentialFind(id, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if billingCredential != nil {
				return resource.RetryableError(errors.New("billing credential still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("billing credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
