package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  lock = %t
}

data "taikun_cloud_credential_aws" "foo" {
  id = resource.taikun_cloud_credential_aws.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialAWS(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_aws.foo",
					"taikun_cloud_credential_aws.foo",
					map[string]struct{}{
						"access_key_id":     {},
						"secret_access_key": {},
					},
				),
			},
		},
	})
}
