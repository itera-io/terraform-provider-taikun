package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccResourceTaikunAppInstanceConfig = testAccAppInstancePrerequisites + `
resource "taikun_app_instance" "foo" {
  name           = "%s"
  namespace      = "%s"
  project_id     = "%s"
  catalog_app_id = local.catalog_app_id
  autosync       = %t
  timeout        = 30

  depends_on = [taikun_catalog_project_binding.foo]
}
`

// TestAccResourceTaikunAppInstance verifies the full install lifecycle:
// bind a catalog app to an existing k8s project, install it, and confirm destroy removes it.
func TestAccResourceTaikunAppInstance(t *testing.T) {
	projectID := os.Getenv("TAIKUN_PROJECT_ID")
	catalogName := utils.RandomTestName()
	appName := testAccAppInstanceName()
	namespace := appName + "-ns"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckAppInstance(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAppInstanceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAppInstanceConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAppInstanceExists(t),
					resource.TestCheckResourceAttrSet("taikun_app_instance.foo", "id"),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "name", appName),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "namespace", namespace),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "project_id", projectID),
					resource.TestCheckResourceAttrSet("taikun_app_instance.foo", "catalog_app_id"),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "autosync", "false"),
				),
			},
		},
	})
}

const testAccResourceTaikunAppInstanceUpdateAutosyncConfig = testAccAppInstancePrerequisites + `
resource "taikun_app_instance" "foo" {
  name           = "%s"
  namespace      = "%s"
  project_id     = "%s"
  catalog_app_id = local.catalog_app_id
  autosync       = %t
  timeout        = 30

  depends_on = [taikun_catalog_project_binding.foo]
}
`

// TestAccResourceTaikunAppInstanceUpdateAutosync verifies in-place autosync toggling
// via POST /api/v1/projectapp/autosync without recreating the app instance.
func TestAccResourceTaikunAppInstanceUpdateAutosync(t *testing.T) {
	projectID := os.Getenv("TAIKUN_PROJECT_ID")
	catalogName := utils.RandomTestName()
	appName := testAccAppInstanceName()
	namespace := appName + "-ns"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckAppInstance(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAppInstanceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAppInstanceUpdateAutosyncConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAppInstanceExists(t),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "autosync", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAppInstanceUpdateAutosyncConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAppInstanceExists(t),
					resource.TestCheckResourceAttr("taikun_app_instance.foo", "autosync", "true"),
				),
			},
		},
	})
}
