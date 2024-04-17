package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This test is removed, because we cannot avoid race conditions.
// This data source lists projects from the default itera organization
// , but they are constantly created and destroyed by other tests run in parallel anyway.
//
//const testAccDataSourceTaikunProjectsConfig = `
//resource "taikun_cloud_credential_openstack" "foo" {
//  name = "%s"
//}
//
//resource "taikun_project" "foo" {
//  name = "%s"
//  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
//}
//
//data "taikun_projects" "all" {
//   depends_on = [
//    taikun_project.foo
//  ]
//}
//`
//
//func TestAccDataSourceTaikunProjects(t *testing.T) {
//	cloudCredentialName := randomTestName()
//	projectName := randomTestName()
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
//		ProviderFactories: testAccProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(testAccDataSourceTaikunProjectsConfig,
//					cloudCredentialName,
//					projectName),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					resource.TestCheckResourceAttr("data.taikun_projects.all", "id", "all"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.#"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.name"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.access_profile_id"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.cloud_credential_id"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.kubernetes_profile_id"),
//					resource.TestCheckResourceAttrSet("data.taikun_projects.all", "projects.0.organization_id"),
//				),
//			},
//		},
//	})
//}

const testAccDataSourceTaikunProjectsConfigWithFilter = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
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
	organizationName := utils.RandomTestName()
	cloudCredentialName := utils.RandomTestName()
	projectCount := 2
	projectName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunProjectsConfigWithFilter,
					organizationName,
					organizationName,
					cloudCredentialName,
					projectCount,
					projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_projects.foo", "projects.#", fmt.Sprint(projectCount)),
				),
			},
		},
	})
}
