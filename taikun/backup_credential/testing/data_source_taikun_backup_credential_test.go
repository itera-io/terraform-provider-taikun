package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunBackupCredentialConfig = `
resource "taikun_organization" "foo" {
  name          = "%s"
  full_name     = "%s"
  discount_rate = 42
}

resource "taikun_backup_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  s3_endpoint = "%s"
  s3_region   = "%s"
}

data "taikun_backup_credential" "foo" {
  id = resource.taikun_backup_credential.foo.id
}
`

// TestAccDataSourceTaikunBackupCredential verifies GET /api/v1/s3credentials/list by id
// matches the managed resource (secret key omitted from API read).
func TestAccDataSourceTaikunBackupCredential(t *testing.T) {
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	backupCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckS3(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBackupCredentialConfig,
					organizationName,
					organizationFullName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
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
