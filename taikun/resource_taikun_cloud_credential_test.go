package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunCloudCredentialConfig = `
resource "taikun_cloud_credential" "foo" {
  type = "%s"
  name = "%s"
  region = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentials(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialConfig, "openstack",
					cloudCredentialName,
					os.Getenv("OS_REGION_NAME"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialsLock(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialConfig, "openstack",
					cloudCredentialName,
					os.Getenv("OS_REGION_NAME"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialsExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "type", "openstack"),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialConfig, "openstack",
					cloudCredentialName,
					os.Getenv("OS_REGION_NAME"),
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialsExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "type", "openstack"),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialsRename(t *testing.T) {
	cloudCredentialName := randomTestName()
	newCloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialConfig, "openstack",
					cloudCredentialName,
					os.Getenv("OS_REGION_NAME"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialsExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "type", "openstack"),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialConfig, "openstack",
					newCloudCredentialName,
					os.Getenv("OS_REGION_NAME"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialsExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "type", "openstack"),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialsExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := client.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCountOpenstack() != 1 {
			return fmt.Errorf("openstack cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTaikunCloudCredentialsDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := client.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCountOpenstack() != 0 {
				return resource.RetryableError(errors.New("openstack cloud credential still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("openstack cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
