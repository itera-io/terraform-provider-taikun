package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_cloud_credentials_azure" "all" {
   depends_on = [
    taikun_cloud_credential_azure.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsAzure(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_azure.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.tenant_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.location"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsAzureWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_cloud_credentials_azure" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_cloud_credential_azure.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsAzureWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsAzureWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.location"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_azure.all", "cloud_credentials.0.tenant_id"),
				),
			},
		},
	})
}
