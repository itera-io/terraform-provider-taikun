package testing

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccDataSourceTaikunAppInstancesConfig = testAccAppInstancePrerequisites + `
resource "taikun_app_instance" "foo" {
  name           = "%s"
  namespace      = "%s"
  project_id     = "%s"
  catalog_app_id = local.catalog_app_id
  timeout        = 30

  depends_on = [taikun_catalog_project_binding.foo]
}

data "taikun_app_instances" "all" {
  depends_on = [taikun_app_instance.foo]
}
`

// TestAccDataSourceTaikunAppInstances verifies GET /api/v1/projectapp/list without filters
// includes the app instance created in the same config.
func TestAccDataSourceTaikunAppInstances(t *testing.T) {
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
				Config: fmt.Sprintf(testAccDataSourceTaikunAppInstancesConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_app_instances.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_app_instances.all", "application_instances.#"),
					testAccCheckTaikunAppInstanceListed("data.taikun_app_instances.all", appName),
				),
			},
		},
	})
}

const testAccDataSourceTaikunAppInstancesWithFilterConfig = testAccAppInstancePrerequisites + `
resource "taikun_app_instance" "foo" {
  name           = "%s"
  namespace      = "%s"
  project_id     = "%s"
  catalog_app_id = local.catalog_app_id
  timeout        = 30

  depends_on = [taikun_catalog_project_binding.foo]
}

data "taikun_app_instances" "all" {
  organization_id = "%s"

  depends_on = [taikun_app_instance.foo]
}
`

// TestAccDataSourceTaikunAppInstancesWithFilter verifies the OrganizationId query parameter
// on GET /api/v1/projectapp/list. Requires TAIKUN_ORGANIZATION_ID for the project's org.
func TestAccDataSourceTaikunAppInstancesWithFilter(t *testing.T) {
	organizationID := os.Getenv("TAIKUN_ORGANIZATION_ID")
	if organizationID == "" {
		t.Skip("TAIKUN_ORGANIZATION_ID must be set to the organization that owns TAIKUN_PROJECT_ID")
	}

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
				Config: fmt.Sprintf(testAccDataSourceTaikunAppInstancesWithFilterConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
					organizationID,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_app_instances.all", "id", organizationID),
					resource.TestCheckResourceAttrSet("data.taikun_app_instances.all", "application_instances.#"),
					testAccCheckTaikunAppInstanceListed("data.taikun_app_instances.all", appName),
				),
			},
		},
	})
}

func testAccCheckTaikunAppInstanceListed(dataSourceName, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["application_instances.#"])
		if err != nil {
			return fmt.Errorf("invalid application_instances.# on %s: %w", dataSourceName, err)
		}

		for i := 0; i < count; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("application_instances.%d.name", i)] == expectedName {
				return nil
			}
		}

		return fmt.Errorf("app instance %q not found in %s list (%d entries)", expectedName, dataSourceName, count)
	}
}
