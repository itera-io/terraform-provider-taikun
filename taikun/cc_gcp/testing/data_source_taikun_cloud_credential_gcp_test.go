package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  lock = true
}

data "taikun_cloud_credential_gcp" "foo" {
  id = resource.taikun_cloud_credential_gcp.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialGCP(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckGCP(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_gcp.foo",
					"taikun_cloud_credential_gcp.foo",
					map[string]struct{}{
						"config_file":    {},
						"import_project": {},
						"az_count":       {},
					},
				),
			},
		},
	})
}
