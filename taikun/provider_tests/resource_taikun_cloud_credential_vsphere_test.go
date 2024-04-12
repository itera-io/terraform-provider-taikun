package provider_tests

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

const testAccResourceTaikunCloudCredentialVsphereConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialVsphere(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialVsphereConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialVsphereExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "password", os.Getenv("VSPHERE_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialVsphereLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialVsphereConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialVsphereExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "password", os.Getenv("VSPHERE_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialVsphereConfig,
					cloudCredentialName,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialVsphereExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "password", os.Getenv("VSPHERE_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
		},
	})
}

const testAccResourceTaikunCloudCredentialVsphereHypervisorConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
  hypervisors = [%s]
  lock =  %t
}
`

func TestAccResourceTaikunCloudCredentialVsphereUpdate(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()
	hypervisor := os.Getenv("VSPHERE_HYPERVISOR")
	hypervisor2 := os.Getenv("VSPHERE_HYPERVISOR2")
	hypervisors_string := fmt.Sprintf("\"%s\"", hypervisor)
	hypervisors_string_update := fmt.Sprintf("\"%s\", \"%s\"", hypervisor, hypervisor2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialVsphereDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialVsphereHypervisorConfig,
					cloudCredentialName,
					hypervisors_string,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialVsphereExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "password", os.Getenv("VSPHERE_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "hypervisors.0", os.Getenv("VSPHERE_HYPERVISOR")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialVsphereHypervisorConfig,
					newCloudCredentialName,
					hypervisors_string_update,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialVsphereExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "api_host", os.Getenv("VSPHERE_API_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "username", os.Getenv("VSPHERE_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "password", os.Getenv("VSPHERE_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "datacenter", os.Getenv("VSPHERE_DATACENTER")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "resource_pool", os.Getenv("VSPHERE_RESOURCE_POOL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "data_store", os.Getenv("VSPHERE_DATA_STORE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "drs_enabled", os.Getenv("VSPHERE_DRS_ENABLED")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "vm_template_name", os.Getenv("VSPHERE_VM_TEMPLATE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "hypervisors.0", os.Getenv("VSPHERE_HYPERVISOR")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "hypervisors.1", os.Getenv("VSPHERE_HYPERVISOR2")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_vsphere.foo", "is_default"),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_name", os.Getenv("VSPHERE_PUBLIC_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_ip_address", os.Getenv("VSPHERE_PUBLIC_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_net_mask", os.Getenv("VSPHERE_PUBLIC_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_gateway", os.Getenv("VSPHERE_PUBLIC_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_begin_allocation_range", os.Getenv("VSPHERE_PUBLIC_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "public_end_allocation_range", os.Getenv("VSPHERE_PUBLIC_END_RANGE")),

					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_name", os.Getenv("VSPHERE_PRIVATE_NETWORK_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_ip_address", os.Getenv("VSPHERE_PRIVATE_NETWORK_ADDRESS")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_net_mask", os.Getenv("VSPHERE_PRIVATE_NETMASK")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_gateway", os.Getenv("VSPHERE_PRIVATE_GATEWAY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_begin_allocation_range", os.Getenv("VSPHERE_PRIVATE_BEGIN_RANGE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_vsphere.foo", "private_end_allocation_range", os.Getenv("VSPHERE_PRIVATE_END_RANGE")),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialVsphereExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_vsphere" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.VsphereCloudCredentialAPI.VsphereList(context.TODO()).Id(id).Execute()
		if err != nil || len(response.GetData()) != 1 {
			return fmt.Errorf("vsphere cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialVsphereDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_vsphere" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.VsphereCloudCredentialAPI.VsphereList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if len(response.GetData()) != 0 {
				return retry.RetryableError(errors.New("Vsphere cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("Vsphere cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
