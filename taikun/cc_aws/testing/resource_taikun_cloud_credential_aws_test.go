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

const testAccResourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  az_count = "%d"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialAWS(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialAWSLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					azCount,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialAWSRename(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					newCloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "az_count", os.Getenv("AWS_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "region", os.Getenv("AWS_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_aws.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialAWSExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.AWSCloudCredentialAPI.AwsList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("aws cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialAWSDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.AWSCloudCredentialAPI.AwsList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("aws cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("aws cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
