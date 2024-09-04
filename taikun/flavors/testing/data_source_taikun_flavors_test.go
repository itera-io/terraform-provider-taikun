package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunFlavorsAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  min_cpu = %d
  max_cpu = %d
  min_ram = %d
  max_ram = %d
}
`

func TestAccDataSourceTaikunFlavorsAWS(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	cpu := 16
	ram := 64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunFlavorsAWSConfig,
					cloudCredentialName,
					cpu, cpu,
					ram, ram,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.#"),
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.0.name"),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.cpu", fmt.Sprint(cpu)),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.ram", fmt.Sprint(ram)),
				),
			},
		},
	})
}

const testAccDataSourceTaikunFlavorsAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id

  min_cpu = %d
  max_cpu = %d
  min_ram = %d
  max_ram = %d
}
`

func TestAccDataSourceTaikunFlavorsAzure(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	cpu := 12
	ram := 48

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunFlavorsAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					cpu, cpu,
					ram, ram,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.#"),
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.0.name"),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.cpu", fmt.Sprint(cpu)),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.ram", fmt.Sprint(ram)),
				),
			},
		},
	})
}

const testAccDataSourceTaikunFlavorsOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id

  min_cpu = %d
  max_cpu = %d
  min_ram = %d
  max_ram = %d
}
`

func TestAccDataSourceTaikunFlavorsOpenStack(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	cpu := 8
	ram := 32

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunFlavorsOpenStackConfig,
					cloudCredentialName,
					cpu, cpu,
					ram, ram,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.#"),
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.0.name"),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.cpu", fmt.Sprint(cpu)),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.ram", fmt.Sprint(ram)),
				),
			},
		},
	})
}

const testAccDataSourceTaikunFlavorsVsphereConfig = `
resource "taikun_cloud_credential_vsphere" "foo" {
  name = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_vsphere.foo.id

  min_cpu = %d
  max_cpu = %d
  min_ram = %d
  max_ram = %d
}
`

func TestAccDataSourceTaikunFlavorsVsphere(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	cpu := 2
	ram := 4

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckVsphere(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunFlavorsVsphereConfig,
					cloudCredentialName,
					cpu, cpu,
					ram, ram,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.#"),
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.0.name"),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.cpu", fmt.Sprint(cpu)),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.ram", fmt.Sprint(ram)),
				),
			},
		},
	})
}
