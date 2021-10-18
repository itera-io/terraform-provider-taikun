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

  is_locked       = %t
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "tenant_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "availability_zone"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "location"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
		},
	})
}
