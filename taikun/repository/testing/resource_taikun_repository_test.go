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

const testAccResourceTaikunRepositoryConfig = `
resource "taikun_repository" "foo" {
  name              = "%s"
  organization_name = "%s"
  private           = false
  enabled           = %s
}
`

func TestAccResourceTaikunRepository(t *testing.T) {
	repositoryName := "taikun-managed-apps"
	repositoryOrgName := "taikun"
	repositoryEnabled := "true"
	//repositoryDisabled := "false"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunRepositoryConfig,
					repositoryName,
					repositoryOrgName,
					repositoryEnabled,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRepositoryExists,
					resource.TestCheckResourceAttr("taikun_repository.foo", "name", repositoryName),
					resource.TestCheckResourceAttr("taikun_repository.foo", "organization_name", repositoryOrgName),
					resource.TestCheckResourceAttr("taikun_repository.foo", "private", "false"),
					resource.TestCheckResourceAttr("taikun_repository.foo", "enabled", repositoryEnabled),
				),
			},
			// We cannot guarantee that argoproj project will be present in staging and dev and we cannot gurantee that taikun-managed-apps will be disableable.
			//{
			//	Config: fmt.Sprintf(testAccResourceTaikunRepositoryConfig,
			//		repositoryName,
			//		repositoryOrgName,
			//		repositoryDisabled,
			//	),
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		testAccCheckTaikunRepositoryExists,
			//		resource.TestCheckResourceAttr("taikun_repository.foo", "name", repositoryName),
			//		resource.TestCheckResourceAttr("taikun_repository.foo", "organization_name", repositoryOrgName),
			//		resource.TestCheckResourceAttr("taikun_repository.foo", "private", "false"),
			//		resource.TestCheckResourceAttr("taikun_repository.foo", "enabled", repositoryDisabled),
			//	),
			//},
		},
	})
}

func testAccCheckTaikunRepositoryExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_repository" {
			continue
		}

		id := rs.Primary.ID

		response, _, err := client.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("repository doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunRepositoryDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_repository" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.CatalogAPI.CatalogList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("repository still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("repository still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
