package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialsProxmoxConfig = `
resource "taikun_cloud_credential_proxmox" "foo" {
  name = "%s"
  hypervisors = [%s]
}

data "taikun_cloud_credentials_proxmox" "all" {
   depends_on = [
    taikun_cloud_credential_proxmox.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsProxmox(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("PROXMOX_HYPERVISOR"))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.api_host"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.client_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.storage"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.vm_template_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.is_default"),

					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_ip_address"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_net_mask"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_gateway"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_begin_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_end_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_bridge"),

					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_ip_address"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_net_mask"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_gateway"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_begin_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_end_allocation_range"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_bridge"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunCloudCredentialsProxmoxWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_proxmox" "foo" {
  name = "%s"
  hypervisors = [%s]
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_cloud_credentials_proxmox" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_cloud_credential_proxmox.foo
  ]
}`

func TestAccDataSourceTaikunCloudCredentialsProxmoxWithFilter(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()
	hypervisors_string := fmt.Sprintf("\"%s\"", os.Getenv("PROXMOX_HYPERVISOR"))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialsProxmoxWithFilterConfig,
					organizationName,
					organizationFullName,
					cloudCredentialName,
					hypervisors_string,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_id"),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.lock", "false"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.is_default"),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("data.taikun_cloud_credentials_proxmox.all", "cloud_credentials.0.public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
		},
	})
}
