package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
)

const testAccResourceTaikunCloudCredentialAWSConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialAWS(t *testing.T) {
	cloudCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "availability_zone", os.Getenv("AWS_AVAILABILITY_ZONE")),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "availability_zone", os.Getenv("AWS_AVAILABILITY_ZONE")),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "availability_zone", os.Getenv("AWS_AVAILABILITY_ZONE")),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialAWSDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialAWSConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "availability_zone", os.Getenv("AWS_AVAILABILITY_ZONE")),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialAWSExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "access_key_id", os.Getenv("AWS_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "secret_access_key", os.Getenv("AWS_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_aws.foo", "availability_zone", os.Getenv("AWS_AVAILABILITY_ZONE")),
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
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.Client.CloudCredentials.CloudCredentialsDashboardList(params, client)
		if err != nil || response.Payload.TotalCountAws != 1 {
			return fmt.Errorf("aws cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialAWSDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_aws" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.Client.CloudCredentials.CloudCredentialsDashboardList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCountAws != 0 {
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
