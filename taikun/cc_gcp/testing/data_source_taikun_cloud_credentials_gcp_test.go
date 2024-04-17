package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsGCPConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_gcp" "foo" {
  organization_id = resource.taikun_organization.foo.id
  name = "%s"
  config_file = "./gcp.json"
  import_project = true
  region = "%s"
  lock = true
}

data "taikun_cloud_credentials_gcp" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_cloud_credential_gcp.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsGCP(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckGCP(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsGCPConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
					os.Getenv("GCP_REGION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.name", cloudCredentialName),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.lock", "true"),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.region", os.Getenv("GCP_REGION")),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_gcp.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_gcp.all", "cloud_credentials.0.organization_id"),
				),
			},
		},
	})
}
