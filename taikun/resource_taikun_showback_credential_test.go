package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/showback"
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
	showbackCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
			{
				ResourceName:      "taikun_showback_credential.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunShowbackCredentialLock(t *testing.T) {
	showbackCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
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
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}

func testAccCheckTaikunShowbackCredentialExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Showback.ShowbackCredentialsList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("showback credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunShowbackCredentialDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.Showback.ShowbackCredentialsList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("showback credential still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("showback credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
