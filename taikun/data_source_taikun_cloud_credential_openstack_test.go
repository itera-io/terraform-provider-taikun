package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"

  is_locked       = %t
}

data "taikun_cloud_credential_openstack" "foo" {
  id = resource.taikun_cloud_credential_openstack.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "user"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "domain"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "project_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "public_network_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "region"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
		},
	})
}
