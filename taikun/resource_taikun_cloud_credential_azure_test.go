package taikun

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_cloud_credential_azure", &resource.Sweeper{
		Name:         "taikun_cloud_credential_azure",
		Dependencies: []string{"taikun_project"},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion)

			var cloudCredentialsList []*models.AzureCredentialsListDto
			for {
				response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(params, apiClient)
				if err != nil {
					return err
				}
				cloudCredentialsList = append(cloudCredentialsList, response.GetPayload().Azure...)
				if len(cloudCredentialsList) == int(response.GetPayload().TotalCountAzure) {
					break
				}
				offset := int32(len(cloudCredentialsList))
				params = params.WithOffset(&offset)
			}

			for _, e := range cloudCredentialsList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
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

const testAccResourceTaikunCloudCredentialAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  availability_zone = "%s"
  location = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialAzure(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("ARM_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("ARM_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("ARM_TENANT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "availability_zone", os.Getenv("ARM_AVAILABILITY_ZONE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("ARM_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialAzureLock(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialOpenStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("ARM_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("ARM_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("ARM_TENANT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "availability_zone", os.Getenv("ARM_AVAILABILITY_ZONE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("ARM_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("ARM_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("ARM_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("ARM_TENANT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "availability_zone", os.Getenv("ARM_AVAILABILITY_ZONE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("ARM_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialAzureRename(t *testing.T) {
	cloudCredentialName := randomTestName()
	newCloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialOpenStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("ARM_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("ARM_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("ARM_TENANT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "availability_zone", os.Getenv("ARM_AVAILABILITY_ZONE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("ARM_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					newCloudCredentialName,
					os.Getenv("ARM_AVAILABILITY_ZONE"),
					os.Getenv("ARM_LOCATION"),
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("ARM_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("ARM_CLIENT_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("ARM_TENANT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("ARM_SUBSCRIPTION_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "availability_zone", os.Getenv("ARM_AVAILABILITY_ZONE")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("ARM_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialAzureExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_azure" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
		if err != nil || response.Payload.TotalCountAzure != 1 {
			return fmt.Errorf("azure cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialAzureDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_azure" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.CloudCredentials.CloudCredentialsDashboardList(params, client)
		if err == nil && response.Payload.TotalCountAzure != 0 {
			return fmt.Errorf("azure cloud credential still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
