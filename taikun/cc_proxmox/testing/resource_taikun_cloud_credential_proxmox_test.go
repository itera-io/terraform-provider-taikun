package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunCloudCredentialProxmoxConfig = `
resource "taikun_cloud_credential_proxmox" "foo" {
  name = "%s"
  hypervisors = [%s]

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialProxmox(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisor := os.Getenv("PROXMOX_HYPERVISOR")
	hypervisors_string := fmt.Sprintf("\"%s\"", hypervisor)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialProxmoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialProxmoxLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	hypervisor := os.Getenv("PROXMOX_HYPERVISOR")
	hypervisors_string := fmt.Sprintf("\"%s\"", hypervisor)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialProxmoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialProxmoxUpdate(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()
	hypervisor := os.Getenv("PROXMOX_HYPERVISOR")
	hypervisor2 := os.Getenv("PROXMOX_HYPERVISOR2")
	hypervisors_string := fmt.Sprintf("\"%s\"", hypervisor)
	hypervisors_string_update := fmt.Sprintf("\"%s\", \"%s\"", hypervisor, hypervisor2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckProxmox(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialProxmoxDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "false"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "hypervisors.0", os.Getenv("PROXMOX_HYPERVISOR")),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					newCloudCredentialName,
					hypervisors_string_update,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "hypervisors.0", os.Getenv("PROXMOX_HYPERVISOR")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "hypervisors.1", os.Getenv("PROXMOX_HYPERVISOR2")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialProxmoxConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialProxmoxExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "api_host", os.Getenv("PROXMOX_API_HOST")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_id", os.Getenv("PROXMOX_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "client_secret", os.Getenv("PROXMOX_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "storage", os.Getenv("PROXMOX_STORAGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "vm_template_name", os.Getenv("PROXMOX_VM_TEMPLATE_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "lock", "false"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "hypervisors.0", os.Getenv("PROXMOX_HYPERVISOR")),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_proxmox.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_ip_address", os.Getenv("PROXMOX_PRIVATE_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_net_mask", os.Getenv("PROXMOX_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_gateway", os.Getenv("PROXMOX_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_begin_allocation_range", os.Getenv("PROXMOX_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_end_allocation_range", os.Getenv("PROXMOX_PRIVATE_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "private_bridge", os.Getenv("PROXMOX_PRIVATE_BRIDGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_ip_address", os.Getenv("PROXMOX_PUBLIC_NETWORK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_net_mask", os.Getenv("PROXMOX_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_gateway", os.Getenv("PROXMOX_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_begin_allocation_range", os.Getenv("PROXMOX_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_end_allocation_range", os.Getenv("PROXMOX_PUBLIC_END_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_proxmox.foo", "public_bridge", os.Getenv("PROXMOX_PUBLIC_BRIDGE")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialProxmoxExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_proxmox" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.ProxmoxCloudCredentialAPI.ProxmoxList(context.TODO()).Id(id).Execute()
		if err != nil || len(response.GetData()) != 1 {
			return fmt.Errorf("proxmox cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialProxmoxDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_proxmox" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.ProxmoxCloudCredentialAPI.ProxmoxList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if len(response.GetData()) != 0 {
				return retry.RetryableError(errors.New("Proxmox cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("Proxmox cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
