package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunBackupCredentialConfig = `
resource "taikun_backup_credential" "foo" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

data "taikun_backup_credential" "foo" {
  id = resource.taikun_backup_credential.foo.id
}
`

func TestAccDataSourceTaikunBackupCredential(t *testing.T) {
	backupCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBackupCredentialConfig,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "is_locked"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "s3_endpoint"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "s3_access_key_id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credential.foo", "s3_region"),
				),
			},
		},
	})
}
