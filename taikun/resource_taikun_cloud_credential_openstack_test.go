package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_cloud_credential_openstack", &resource.Sweeper{
		Name:         "taikun_cloud_credential_openstack",
		Dependencies: []string{"taikun_project"},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion)

			var cloudCredentialsList []*models.OpenstackCredentialsListDto
			for {
				response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
				if err != nil {
					return err
				}
				cloudCredentialsList = append(cloudCredentialsList, response.GetPayload().Openstack...)
				if len(cloudCredentialsList) == int(response.GetPayload().TotalCountOpenstack) {
					break
				}
				offset := int32(len(cloudCredentialsList))
				params = params.WithOffset(&offset)
			}

			for _, e := range cloudCredentialsList {
				if shouldSweep(e.Name) {
					params := cloud_credentials.NewCloudCredentialsDeleteParams().WithV(ApiVersion).WithCloudID(e.ID)
					_, _, err = apiClient.client.CloudCredentials.CloudCredentialsDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunCloudCredentialOpenStackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialOpenStack(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
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
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
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
	cloudCredentialName := randomTestName()
	newCloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
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
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_openstack" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
		if err != nil || response.Payload.TotalCountOpenstack != 1 {
			return fmt.Errorf("openstack cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialOpenStackDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_openstack" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCountOpenstack != 0 {
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
