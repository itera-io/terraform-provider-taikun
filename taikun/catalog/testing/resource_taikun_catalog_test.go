package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"
)

const testAccResourceTaikunCatalogConfig = `
resource "taikun_catalog" "foo" {
  name="%s"
  description="%s"
  projects=[]

  %s
}
`

var testAccResourceWordpress string = `
  application {
    name="wordpress"
    repository="taikun-managed-apps"
  }
`

var testAccResourceNginx string = `
  application {
    name="nginx"
    repository="taikun-managed-apps"
  }
`

func TestAccResourceTaikunCatalogAppBinding(t *testing.T) {
	catalogName := utils.RandomTestName()
	catalogDescription := utils.RandomTestName()
	catalogDescriptionChanged := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCatalogConfig,
					catalogName,
					catalogDescription,
					testAccResourceWordpress,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCatalogExists,
					resource.TestCheckResourceAttr("taikun_catalog.foo", "name", catalogName),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "description", catalogDescription),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "application.#", "1"),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "application.0.name", "wordpress"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCatalogConfig,
					catalogName,
					catalogDescriptionChanged,
					testAccResourceNginx,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCatalogExists,
					resource.TestCheckResourceAttr("taikun_catalog.foo", "name", catalogName),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "description", catalogDescriptionChanged),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "application.#", "1"),
					resource.TestCheckResourceAttr("taikun_catalog.foo", "application.0.name", "nginx"),
				),
			},
		},
	})
}

func testAccCheckTaikunCatalogExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_catalog" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.CatalogAPI.CatalogList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("catalog doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCatalogDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_catalog" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.CatalogAPI.CatalogList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("catalog still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("catalog still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
