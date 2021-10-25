package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunProjectsConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
}

data "taikun_projects" "all" {
   depends_on = [
    taikun_project.foo
  ]
}
`

func TestAccDataSourceTaikunProjects(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectsConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_projects.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.#"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.access_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.alerting_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.cloud_credential_id"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.organization_id"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunProjectsConfigWithFilter = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
  organization_id = resource.taikun_organization.foo.id
}

resource "taikun_project" "foo" {
  count = %d
  name = "%s-${count.index}"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_projects" "foo" {
  depends_on = [
    resource.taikun_project.foo
  ]
  organization_id = resource.taikun_organization.foo.id
}
`

func TestAccDataSourceTaikunProjectsWithFilter(t *testing.T) {
	organizationName := randomTestName()
	cloudCredentialName := randomTestName()
	projectCount := 3
	projectName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectsConfigWithFilter,
					organizationName,
					organizationName,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectCount,
					projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_projects.foo", "projects.#", fmt.Sprint(projectCount)),
				),
			},
		},
	})
}
