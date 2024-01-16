package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Check: checkDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_backup_credential.foo",
					"taikun_backup_credential.foo",
					map[string]struct{}{
						"s3_secret_access_key": {},
					},
				),
			},
		},
	})
}
