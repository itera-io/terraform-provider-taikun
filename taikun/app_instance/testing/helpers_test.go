package testing

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

// testAccAppInstancePrerequisites is shared HCL for catalog + project binding.
// App install requires a bound catalog on an existing k8s project (TAIKUN_PROJECT_ID).
const testAccAppInstancePrerequisites = `
resource "taikun_catalog" "foo" {
  name        = "%s"
  description = "Acceptance test catalog for taikun_app_instance"
  projects    = []

  application {
    name       = "ingress-nginx"
    repository = "taikun-managed-apps"
  }
}

resource "taikun_catalog_project_binding" "foo" {
  catalog_name = taikun_catalog.foo.name
  project_id   = "%s"
  is_bound     = true
}

locals {
  catalog_app_id = [
    for app in tolist(taikun_catalog.foo.application) :
    app.id if app.name == "ingress-nginx" && app.repository == "taikun-managed-apps"
  ][0]
}
`

// testAccPreCheckAppInstance ensures API credentials and a k8s project are available.
// Goal: skip early when the platform prerequisite (running project) is missing.
func testAccPreCheckAppInstance(t *testing.T) {
	t.Helper()
	utils_testing.TestAccPreCheck(t)
	if os.Getenv("TAIKUN_PROJECT_ID") == "" {
		t.Skip("TAIKUN_PROJECT_ID must be set to a running k8s project ID")
	}
}

// testAccAppInstanceName returns a schema-valid lowercase app/namespace name.
func testAccAppInstanceName() string {
	return strings.ToLower(utils.ShortRandomTestName())
}

func testAccCheckTaikunAppInstanceExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_app_instance" {
				continue
			}

			id, _ := utils.Atoi32(rs.Primary.ID)
			_, _, err := client.Client.ProjectAppsAPI.ProjectappDetails(t.Context(), id).Execute()
			if err != nil {
				return fmt.Errorf("app instance doesn't exist (id = %s)", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTaikunAppInstanceDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_app_instance" {
				continue
			}

			id, _ := utils.Atoi32(rs.Primary.ID)
			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				response, _, err := client.Client.ProjectAppsAPI.ProjectappList(t.Context()).Id(id).Execute()
				if err != nil {
					return retry.NonRetryableError(err)
				}
				if response.GetTotalCount() != 0 {
					return retry.RetryableError(errors.New("app instance still exists"))
				}
				return nil
			})
			if utils.TimedOut(retryErr) {
				return errors.New("app instance still exists (timed out)")
			}
			if retryErr != nil {
				return retryErr
			}
		}

		return nil
	}
}
