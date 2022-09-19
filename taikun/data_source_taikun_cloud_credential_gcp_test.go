package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  zone = "%s"
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
					os.Getenv("GCP_ZONE"),
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
