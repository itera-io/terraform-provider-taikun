package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  enable_auto_upgrade = %t
  enable_monitoring = %t
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_project.foo", "enable_auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("data.taikun_project.foo", "enable_monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("data.taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("data.taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("data.taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("data.taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_project.foo", "organization_id"),
				),
			},
		},
	})
}
