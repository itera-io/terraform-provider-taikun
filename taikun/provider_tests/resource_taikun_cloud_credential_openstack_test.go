package provider_tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunCloudCredentialOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialOpenStack(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialOpenStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialOpenStackLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialOpenStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialOpenStackRename(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckOpenStack(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialOpenStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialOpenStackConfig,
					cloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialOpenStackConfig,
					newCloudCredentialName,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialOpenStackExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "user", os.Getenv("OS_USERNAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "password", os.Getenv("OS_PASSWORD")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "url", os.Getenv("OS_AUTH_URL")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "domain", os.Getenv("OS_USER_DOMAIN_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "project_name", os.Getenv("OS_PROJECT_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "public_network_name", os.Getenv("OS_INTERFACE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "region", os.Getenv("OS_REGION_NAME")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "continent", os.Getenv("OS_CONTINENT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_openstack.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "project_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_openstack.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialOpenStackExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_openstack" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.OpenstackCloudCredentialAPI.OpenstackList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("openstack cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialOpenStackDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_openstack" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.OpenstackCloudCredentialAPI.OpenstackList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("openstack cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("openstack cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
