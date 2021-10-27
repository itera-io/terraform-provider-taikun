package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunBackupCredentialsConfig = `
resource "taikun_backup_credential" "foo" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

data "taikun_backup_credentials" "all" {
   depends_on = [
    taikun_backup_credential.foo
  ]
}`

func TestAccDataSourceTaikunBackupCredentials(t *testing.T) {
	backupCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBackupCredentialsConfig,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_backup_credentials.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_region"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_endpoint"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_access_key_id"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunBackupCredentialsWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_backup_credential" "foo" {
  name            = "%s"
  organization_id = resource.taikun_organization.foo.id

  s3_endpoint = "%s"
  s3_region   = "%s"
}

data "taikun_backup_credentials" "all" {
  organization_id = resource.taikun_organization.foo.id

   depends_on = [
    taikun_backup_credential.foo
  ]
}`

func TestAccDataSourceTaikunBackupCredentialsWithFilter(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	backupCredentialName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunBackupCredentialsWithFilterConfig,
					organizationName,
					organizationFullName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_backup_credentials.all", "backup_credentials.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.#"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.created_by"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.is_default"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_region"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_endpoint"),
					resource.TestCheckResourceAttrSet("data.taikun_backup_credentials.all", "backup_credentials.0.s3_access_key_id"),
				),
			},
		},
	})
}
