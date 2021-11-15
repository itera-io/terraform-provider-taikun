package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

data "taikun_cloud_credentials_openstack" "all" {
   depends_on = [
    taikun_cloud_credential_openstack.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsOpenStackConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_openstack.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.user"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.project_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.project_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.public_network_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.domain"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsOpenStackWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_cloud_credentials_openstack" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_cloud_credential_openstack.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsOpenStackWithFilter(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	cloudCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsOpenStackWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.user"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.project_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.project_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.public_network_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.domain"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_openstack.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}
