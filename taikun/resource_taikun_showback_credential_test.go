package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
	"os"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_showback_credential", &resource.Sweeper{
		Name: "taikun_showback_credential",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion)

			var showbackCredentialsList []*models.ShowbackCredentialsListDto
			for {
				response, err := apiClient.client.Showback.ShowbackCredentialsList(params, apiClient)
				if err != nil {
					return err
				}
				showbackCredentialsList = append(showbackCredentialsList, response.GetPayload().Data...)
				if len(showbackCredentialsList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(showbackCredentialsList))
				params = params.WithOffset(&offset)
			}

			for _, e := range showbackCredentialsList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := showback.NewShowbackDeleteShowbackCredentialParams().WithV(ApiVersion).WithBody(&models.DeleteShowbackCredentialCommand{ID: e.ID})
					_, err = apiClient.client.Showback.ShowbackDeleteShowbackCredential(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunShowbackCredentialConfig = `
resource "taikun_showback_credential" "foo" {
  name            = "%s"
  is_locked       = %t

  password = "%s"
  url = "%s"
  username = "%s"
}
`

func TestAccResourceTaikunShowbackCredential(t *testing.T) {
	showbackCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					false,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "is_locked", "false"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}

func TestAccResourceTaikunShowbackCredentialLock(t *testing.T) {
	showbackCredentialName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					false,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "is_locked", "false"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackCredentialConfig,
					showbackCredentialName,
					true,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunShowbackCredentialExists,
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "name", showbackCredentialName),
					resource.TestCheckResourceAttr("taikun_showback_credential.foo", "is_locked", "true"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "password"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "url"),
					resource.TestCheckResourceAttrSet("taikun_showback_credential.foo", "username"),
				),
			},
		},
	})
}

func testAccCheckTaikunShowbackCredentialExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Showback.ShowbackCredentialsList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("showback credential doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunShowbackCredentialDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_credential" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Showback.ShowbackCredentialsList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("showback credential still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
