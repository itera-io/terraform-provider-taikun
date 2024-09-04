package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
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
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
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
