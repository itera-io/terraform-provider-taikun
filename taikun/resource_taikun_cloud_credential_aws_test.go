package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/chnyda/taikungoclient"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  az_count = "%d"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialAWS(t *testing.T) {
	cloudCredentialName := randomTestName()
	azCount, _ := atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
	cloudCredentialName := randomTestName()
	azCount, _ := atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
		},
	})
}

func TestAccResourceTaikunCloudCredentialAWSRename(t *testing.T) {
	cloudCredentialName := randomTestName()
	newCloudCredentialName := randomTestName()
	azCount, _ := atoi32(os.Getenv("AWS_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := client.Client.CloudCredentialApi.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCountAws() != 1 {
			return fmt.Errorf("aws cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialAWSDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := client.Client.CloudCredentialApi.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCountAws() != 0 {
				return resource.RetryableError(errors.New("aws cloud credential still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("aws cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
