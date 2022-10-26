package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  auto_upgrade = %t
  monitoring = %t
  expiration_date = "%s"
}

data "taikun_project" "foo" {
  id = resource.taikun_project.foo.id
}
`

func TestAccDataSourceTaikunProject(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_project.foo",
					"taikun_project.foo",
				),
			},
		},
	})
}
