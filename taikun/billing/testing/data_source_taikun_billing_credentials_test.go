package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunBillingCredentialsConfig = `
resource "taikun_organization" "foo" {
  name          = "%s"
  full_name     = "%s"
  discount_rate = 42
}

resource "taikun_billing_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  prometheus_password = "%s"
  prometheus_url      = "%s"
  prometheus_username = "%s"
}

data "taikun_billing_credentials" "all" {
   depends_on = [
    taikun_billing_credential.foo
  ]
}`

// TestAccDataSourceTaikunBillingCredentials lists all billing credentials.
// Skipped because Managing Operation/Billing Credentials requires a Partner role on the Taikun API,
// which is not possessed by the standard robot account credentials in dev-env.sh.
func TestAccDataSourceTaikunBillingCredentials(t *testing.T) {
	t.Skip("requires Partner API role privileges; re-enable when Partner credentials are available in dev-env.sh")

	orgName := utils.RandomTestName()
	orgFullName := utils.RandomTestName()
	billingCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingCredentialsConfig,
					orgName,
					orgFullName,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_billing_credentials.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_username"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunBillingCredentialsWithFilterConfig = `
resource "taikun_organization" "foo" {
  name          = "%s"
  full_name     = "%s"
  discount_rate = 42
}

resource "taikun_billing_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  prometheus_password = "%s"
  prometheus_url      = "%s"
  prometheus_username = "%s"
}

data "taikun_billing_credentials" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_billing_credential.foo
  ]
}`

// TestAccDataSourceTaikunBillingCredentialsWithFilter lists billing credentials filtered by OrganizationId.
// Skipped because Managing Operation/Billing Credentials requires a Partner role on the Taikun API,
// which is not possessed by the standard robot account credentials in dev-env.sh.
func TestAccDataSourceTaikunBillingCredentialsWithFilter(t *testing.T) {
	t.Skip("requires Partner API role privileges; re-enable when Partner credentials are available in dev-env.sh")

	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	billingCredentialName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBillingCredentialsWithFilterConfig,
					organizationName,
					organizationFullName,
					billingCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_billing_credentials.all", "billing_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_url"),
					resource.TestCheckResourceAttrSet("data.taikun_billing_credentials.all", "billing_credentials.0.prometheus_username"),
				),
			},
		},
	})
}
