package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/ops_credentials"
	"github.com/itera-io/taikungoclient/models"
	"os"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_billing_credential", &resource.Sweeper{
		Name:         "taikun_billing_credential",
		Dependencies: []string{"taikun_organization", "taikun_billing_rule"},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion)

			var operationCredentialsList []*models.OperationCredentialsListDto
			for {
				response, err := apiClient.client.OpsCredentials.OpsCredentialsList(params, apiClient)
				if err != nil {
					return err
				}
				operationCredentialsList = append(operationCredentialsList, response.GetPayload().Data...)
				if len(operationCredentialsList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(operationCredentialsList))
				params = params.WithOffset(&offset)
			}

			for _, e := range operationCredentialsList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := ops_credentials.NewOpsCredentialsDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, _, err = apiClient.client.OpsCredentials.OpsCredentialsDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunBillingCredentialConfig = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"
  is_locked       = %t

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}
`

func TestAccResourceTaikunBillingCredential(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, false, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "is_locked", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}

func TestAccResourceTaikunBillingCredentialLock(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, false, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "is_locked", "false"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingCredentialConfig, firstName, true, os.Getenv("PROMETHEUS_PASSWORD"), os.Getenv("PROMETHEUS_URL"), os.Getenv("PROMETHEUS_USERNAME")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingCredentialExists,
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_credential.foo", "is_locked", "true"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_password"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_url"),
					resource.TestCheckResourceAttrSet("taikun_billing_credential.foo", "prometheus_username"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingCredentialExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.OpsCredentials.OpsCredentialsList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("billing credential doesn't exist")
		}
	}

	return nil
}

func testAccCheckTaikunBillingCredentialDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := ops_credentials.NewOpsCredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.OpsCredentials.OpsCredentialsList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("billing credential still exists")
		}
	}

	return nil
}
