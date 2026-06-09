package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccDataSourceTaikunAppInstanceConfig = testAccAppInstancePrerequisites + `
resource "taikun_app_instance" "foo" {
  name           = "%s"
  namespace      = "%s"
  project_id     = "%s"
  catalog_app_id = local.catalog_app_id
  timeout        = 30

  depends_on = [taikun_catalog_project_binding.foo]
}

data "taikun_app_instance" "foo" {
  id = taikun_app_instance.foo.id
}
`

// TestAccDataSourceTaikunAppInstance verifies GET /api/v1/projectapp/{id} via the data source
// returns the same attributes as the managed resource.
func TestAccDataSourceTaikunAppInstance(t *testing.T) {
	projectID := os.Getenv("TAIKUN_PROJECT_ID")
	catalogName := utils.RandomTestName()
	appName := testAccAppInstanceName()
	namespace := appName + "-ns"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckAppInstance(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAppInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunAppInstanceConfig,
					catalogName,
					projectID,
					appName,
					namespace,
					projectID,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_app_instance.foo",
					"taikun_app_instance.foo",
				),
			},
		},
	})
}
