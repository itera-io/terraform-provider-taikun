package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunFlavorsAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

data "taikun_flavors" "foo" {
  cloud_type = "AWS"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  min_cpu = %d
  max_cpu = %d
  min_ram = %d
  max_ram = %d
}
`

func TestAccDataSourceTaikunFlavorsAWS(t *testing.T) {
	cloudCredentialName := randomTestName()
	cpu := 16
	ram := 64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunFlavorsAWSConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					cpu, cpu,
					ram, ram,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.#"),
					resource.TestCheckResourceAttrSet("data.taikun_flavors.foo", "flavors.0.name"),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.cpu", fmt.Sprint(cpu)),
					resource.TestCheckResourceAttr("data.taikun_flavors.foo", "flavors.0.ram", fmt.Sprint(ram)),
				),
			},
		},
	})
}
