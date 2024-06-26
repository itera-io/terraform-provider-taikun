package testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"
)

const testAccDataSourceTaikunRepositoryConfig = `
data "taikun_repository" "foo" {
  name = "%s"
  organization_name = "%s"
  private="false"
}
`

func TestAccDataSourceTaikunRepository(t *testing.T) {
	repositoryName := "argo"
	repositoryOrgName := "argoproj"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunRepositoryConfig,
					repositoryName,
					repositoryOrgName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunRepositoryExists,
					resource.TestCheckResourceAttr("data.taikun_repository.foo", "name", repositoryName),
					resource.TestCheckResourceAttr("data.taikun_repository.foo", "organization_name", repositoryOrgName),
					resource.TestCheckResourceAttr("data.taikun_repository.foo", "private", "false"),
					resource.TestCheckResourceAttrSet("data.taikun_repository.foo", "enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_repository.foo", "id"),
				),
			},
		},
	})
}
