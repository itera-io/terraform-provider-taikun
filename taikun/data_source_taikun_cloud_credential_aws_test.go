package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"

  is_locked       = %t
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "availability_zone"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "region"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
		},
	})
}
