package testing

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

const testAccResourceTaikunCloudCredentialAzureConfig = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  az_count = "%d"
  location = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialAzure(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AZURE_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
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
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AZURE_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialAzureRename(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AZURE_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAzure(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAzureDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					cloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_azure.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAzureConfig,
					newCloudCredentialName,
					azCount,
					os.Getenv("AZURE_LOCATION"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAzureExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_id", os.Getenv("AZURE_CLIENT_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "client_secret", os.Getenv("AZURE_SECRET")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "tenant_id", os.Getenv("AZURE_TENANT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "subscription_id", os.Getenv("AZURE_SUBSCRIPTION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "az_count", os.Getenv("AZURE_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_azure.foo", "location", os.Getenv("AZURE_LOCATION")),
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
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_azure" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.AzureCloudCredentialAPI.AzureList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("azure cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialAzureDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_azure" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.AzureCloudCredentialAPI.AzureList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("azure cloud credential still exists ()"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("azure cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
