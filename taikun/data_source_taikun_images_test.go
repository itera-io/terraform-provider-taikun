package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunImagesDeprecatedAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  aws_limit = 3
}
`

func TestAccDataSourceTaikunImagesDeprecatedAWS(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesDeprecatedAWSConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_images.foo", "images.#", "3"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.id"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunImagesDeprecatedAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  azure_publisher = "Canonical"
  azure_offer = "UbuntuServer"
  azure_sku = "19.04"
}
`

func TestAccDataSourceTaikunImagesDeprecatedAzure(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesDeprecatedAzureConfig,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.id"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunImagesDeprecatedOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
`

func TestAccDataSourceTaikunImagesDeprecatedOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesDeprecatedOpenStackConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.#"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images.foo", "images.0.id"),
				),
			},
		},
	})
}
