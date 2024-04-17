package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunImagesAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

data "taikun_images_aws" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  latest = true
  owners = ["Canonical"]
}
`

func TestAccDataSourceTaikunImagesAWS(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunImagesAWSConfig,
					cloudCredentialName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_images_aws.foo", "images.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_images_aws.foo", "images.0.id"),
				),
			},
		},
	})
}
