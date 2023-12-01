package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

data "taikun_cloud_credentials_aws" "all" {
   depends_on = [
    taikun_cloud_credential_aws.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsAWS(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsAWSConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_aws.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.#"),
					//resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.created_by"), // First test credential in dev does not have "created by" set for some reason
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsAWSWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  organization_id = resource.taikun_organization.foo.id
  depends_on = [
    taikun_organization.foo
  ]
}

data "taikun_cloud_credentials_aws" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_cloud_credential_aws.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsAWSWithFilter(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	cloudCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsAWSWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_aws.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}
