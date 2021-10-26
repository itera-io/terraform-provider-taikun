package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  availability_zone = "%s"
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
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
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
