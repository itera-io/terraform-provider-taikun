package testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"
)

const testAccDataSourceTaikunCatalogConfig = `
resource "taikun_catalog" "foo" {
  name="%s"
  description="%s"
  projects=[]

  application {
    name="wordpress"
    repository="taikun-managed-apps"
  }
}

data "taikun_catalog" "foo" {
  name = resource.taikun_catalog.foo.name
}
`

func TestAccDataSourceTaikunCatalog(t *testing.T) {
	catalogName := utils.RandomTestName()
	catalogDescription := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCatalogConfig,
					catalogName,
					catalogDescription,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_catalog.foo",
					"taikun_catalog.foo",
				),
			},
		},
	})
}
