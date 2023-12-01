package taikun

import (
	"fmt"
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
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_gcp.foo",
					"taikun_cloud_credential_gcp.foo",
					map[string]struct{}{
						"config_file":    {},
						"import_project": {},
					},
				),
			},
		},
	})
}
