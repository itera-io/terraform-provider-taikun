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

const testAccResourceTaikunCloudCredentialZadaraConfig = `
resource "taikun_cloud_credential_zadara" "foo" {
  name = "%s"
  az_count = "%d"

  lock       = %t
}
`

func TestAccResourceTaikunCloudCredentialZadara(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("ZADARA_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialZadaraDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialZadaraConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialZadaraExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "access_key_id", os.Getenv("ZADARA_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "secret_access_key", os.Getenv("ZADARA_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "az_count", os.Getenv("ZADARA_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "region", os.Getenv("ZADARA_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialZadaraLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("ZADARA_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialZadaraDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialZadaraConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialZadaraExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "access_key_id", os.Getenv("ZADARA_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "secret_access_key", os.Getenv("ZADARA_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "az_count", os.Getenv("ZADARA_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "region", os.Getenv("ZADARA_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialZadaraConfig,
					cloudCredentialName,
					azCount,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialZadaraExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "access_key_id", os.Getenv("ZADARA_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "secret_access_key", os.Getenv("ZADARA_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "az_count", os.Getenv("ZADARA_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "region", os.Getenv("ZADARA_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "is_default"),
				),
			},
		},
	})
}

func TestAccResourceTaikunCloudCredentialZadaraRename(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	newCloudCredentialName := utils.RandomTestName()
	azCount, _ := utils.Atoi32(os.Getenv("ZADARA_AZ_COUNT"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckZadara(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunCloudCredentialZadaraDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialZadaraConfig,
					cloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialZadaraExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "name", cloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "access_key_id", os.Getenv("ZADARA_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "secret_access_key", os.Getenv("ZADARA_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "az_count", os.Getenv("ZADARA_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "region", os.Getenv("ZADARA_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "is_default"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunCloudCredentialZadaraConfig,
					newCloudCredentialName,
					azCount,
					false,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunCloudCredentialZadaraExists,
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "name", newCloudCredentialName),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "access_key_id", os.Getenv("ZADARA_ACCESS_KEY_ID")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "secret_access_key", os.Getenv("ZADARA_SECRET_ACCESS_KEY")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "az_count", os.Getenv("ZADARA_AZ_COUNT")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "region", os.Getenv("ZADARA_DEFAULT_REGION")),
					resource.TestCheckResourceAttr("taikun_cloud_credential_zadara.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_cloud_credential_zadara.foo", "is_default"),
				),
			},
		},
	})
}

func testAccCheckTaikunCloudCredentialZadaraExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_zadara" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.ZadaraCloudCredentialAPI.ZadaraList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("zadara cloud credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunCloudCredentialZadaraDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_cloud_credential_zadara" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.ZadaraCloudCredentialAPI.ZadaraList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("zadara cloud credential still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("zadara cloud credential still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
