package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"

  lock       = %t
}

data "taikun_cloud_credential_azure" "foo" {
  id = resource.taikun_cloud_credential_azure.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialAzure(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_azure.foo",
					"taikun_cloud_credential_azure.foo",
					map[string]struct{}{
						"client_id":       {},
						"subscription_id": {},
						"client_secret":   {},
					},
				),
			},
		},
	})
}
