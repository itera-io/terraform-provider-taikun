package testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccDataSourceTaikunProjectSubnetsConfig = `
data "taikun_project_subnets" "foo" {
  project_id = "%s"
}
`

func TestAccDataSourceTaikunProjectSubnets(t *testing.T) {
	projectId := os.Getenv("TAIKUN_PROJECT_ID")
	if projectId == "" {
		t.Skip("TAIKUN_PROJECT_ID must be set for this acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectSubnetsConfig, projectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_project_subnets.foo", "project_id"),
					resource.TestCheckResourceAttrSet("data.taikun_project_subnets.foo", "subnets.#"),
				),
			},
		},
	})
}
