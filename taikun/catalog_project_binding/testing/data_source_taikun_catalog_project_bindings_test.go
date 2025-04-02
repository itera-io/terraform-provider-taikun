package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceCatalogProjectBindingsConfig = `
resource "taikun_catalog" "foo" {
	name="%s"
	description="%s"

	application {
		name="wordpress"
		repository="taikun-managed-apps"
	}
	application {
		name="apache"
		repository="taikun-managed-apps"
	}
}

resource "taikun_catalog" "bar" {
	name="%s"
	description="%s"

	application {
		name="apache"
		repository="taikun-managed-apps"
	}
}

data "taikun_catalog_project_binding" "foo" {
	catalog_name = taikun_catalog.foo.name
	project_id = 42
}

data "taikun_catalog_project_bindings" "all" {
	catalog_name = taikun_catalog.bar.name
	depends_on = [
		taikun_catalog.foo,
		taikun_catalog.bar
	]
}
`

func TestAccDataSourceCatalogProjectBindings(t *testing.T) {
	cat1Name := utils.RandomTestName()
	cat2Name := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceCatalogProjectBindingsConfig,
					cat1Name,
					cat1Name,
					cat2Name,
					cat2Name,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_catalog_project_binding.foo", "catalog_name", cat1Name),
					resource.TestCheckResourceAttr("data.taikun_catalog_project_binding.foo", "project_id", "42"),
					resource.TestCheckResourceAttr("data.taikun_catalog_project_binding.foo", "is_bound", "false"),

					resource.TestCheckResourceAttr("data.taikun_catalog_project_bindings.all", "catalog_name", cat2Name),
					resource.TestCheckResourceAttr("data.taikun_catalog_project_bindings.all", "catalog_project_bindings.#", "0"),
				),
			},
		},
	})
}
