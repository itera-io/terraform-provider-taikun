package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsVsphereConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  hypervisors = [%s]
}

data "taikun_cloud_credentials_vsphere" "all" {
   depends_on = [
    taikun_cloud_credential_vsphere.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsVsphere(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("VSPHERE_HYPERVISOR"))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsVsphereConfig,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.api_host"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.username"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.datacenter"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.resource_pool"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.data_store"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.drs_enabled"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.vm_template_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.continent"),

					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_ip_address"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_net_mask"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_gateway"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_begin_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_end_allocation_range"),

					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_ip_address"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_net_mask"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_gateway"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_begin_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_end_allocation_range"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsVsphereWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  hypervisors = [%s]
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_cloud_credentials_vsphere" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_cloud_credential_vsphere.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsVsphereWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("VSPHERE_HYPERVISOR"))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsVsphereWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.lock", "false"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.organization_name"),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.hypervisors.0", os.Getenv("VSPHERE_HYPERVISOR")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_vsphere.all", "cloud_credentials.0.private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
		},
	})
}
