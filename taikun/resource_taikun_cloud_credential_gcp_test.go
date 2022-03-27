package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourceTaikunCloudCredentialGCPConfig = `
resource "taikun_cloud_credential_gcp" "foo" {
  name = "%s"
  # FIXME
  lock = %t
}
`

func TestAccResourceTaikunCloudCredentialGCP(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					// FIXME
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					// FIXME
					false,
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialGCPLock(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckGCP(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialGCPDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialGCPConfig,
					cloudCredentialName,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialGCPExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_gcp.foo", "lock", "true"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialGCPExists(state *terraform.State) error {
	// FIXME

	return nil
}

func testAccCheckTaikunCloudCredentialGCPDestroy(state *terraform.State) error {
	// FIXME

	return nil
}
