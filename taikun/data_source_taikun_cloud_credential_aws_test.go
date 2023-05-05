package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  lock = %t
  az_count = 1
}

data "taikun_cloud_credential_aws" "foo" {
  id = resource.taikun_cloud_credential_aws.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialAWS(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					false,
				),
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
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
