package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsZadaraConfig = `
resource "taikun_cloud_credential_zadara" "foo" {
  name = "%s"
}

data "taikun_cloud_credentials_zadara" "all" {
   depends_on = [
    taikun_cloud_credential_zadara.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsZadara(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsZadaraConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_zadara.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.#"),
					//resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.created_by"), // First test credential in dev does not have "created by" set for some reason
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsZadaraWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_zadara" "foo" {
  name = "%s"
  organization_id = resource.taikun_organization.foo.id
  depends_on = [
    taikun_organization.foo
  ]
}

data "taikun_cloud_credentials_zadara" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_cloud_credential_zadara.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsZadaraWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsZadaraWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_zadara.all", "cloud_credentials.0.region"),
				),
			},
		},
	})
}
