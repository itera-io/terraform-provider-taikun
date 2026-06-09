package testing

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunOrganizationConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
}
`

func TestAccResourceTaikunOrganization(t *testing.T) {
	name := utils.RandomTestName()
	fullName := utils.RandomString()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOrganizationDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
					name,
					fullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists(t),
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fullName),
					resource.TestCheckResourceAttrSet("taikun_organization.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_organization.foo", "created_at"),
				),
			},
			{
				ResourceName:      "taikun_organization.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunOrganizationUpdate(t *testing.T) {
	name := utils.RandomTestName()
	newName := utils.RandomTestName()
	fullName := utils.RandomString()
	newFullName := utils.RandomString()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOrganizationDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
					name,
					fullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists(t),
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fullName),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
					newName,
					newFullName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists(t),
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", newFullName),
				),
			},
		},
	})
}

func testAccCheckTaikunOrganizationExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_organization" {
				continue
			}

			id, _ := utils.Atoi32(rs.Primary.ID)
			response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(t.Context()).Id(id).Execute()
			if err != nil || len(response.GetData()) != 1 {
				return fmt.Errorf("organization doesn't exist (id = %s)", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTaikunOrganizationDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_organization" {
				continue
			}

			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				id, _ := utils.Atoi32(rs.Primary.ID)
				response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(t.Context()).Id(id).Execute()

				if err != nil {
					return retry.NonRetryableError(err)
				}
				if response.GetTotalCount() != 0 {
					return retry.RetryableError(errors.New("organization still exists"))
				}
				return nil
			})
			if utils.TimedOut(retryErr) {
				return errors.New("organization still exists (timed out)")
			}
			if retryErr != nil {
				return retryErr
			}
		}

		return nil
	}
}
