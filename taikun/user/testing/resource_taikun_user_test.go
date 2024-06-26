package testing

import (
	"context"
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

const testAccResourceTaikunUserConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "%s"

  display_name        = "%s"
}
`

func TestAccResourceTaikunUser(t *testing.T) {
	userName := utils.RandomTestName()
	email := utils.RandomEmail()
	role := "User"
	displayName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					userName,
					email,
					role,
					displayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", role),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_csm"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_disabled"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_approved_by_partner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_notification_enabled"),
				),
			},
			{
				ResourceName:      "taikun_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunUserUpdate(t *testing.T) {
	userName := utils.RandomTestName()
	email := utils.RandomEmail()
	role := "Manager"
	displayName := utils.RandomTestName()
	newUserName := utils.RandomTestName()
	newEmail := utils.RandomEmail()
	newRole := "Manager"
	newDisplayName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					userName,
					email,
					role,
					displayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", role),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttr("taikun_user.foo", "is_disabled", "false"),
					resource.TestCheckResourceAttr("taikun_user.foo", "is_approved_by_partner", "true"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_csm"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_notification_enabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					newUserName,
					newEmail,
					newRole,
					newDisplayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", newUserName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", newEmail),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", newRole),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", newDisplayName),
					resource.TestCheckResourceAttr("taikun_user.foo", "is_disabled", "false"),
					resource.TestCheckResourceAttr("taikun_user.foo", "is_approved_by_partner", "true"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_csm"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_notification_enabled"),
				),
			},
		},
	})
}

func testAccCheckTaikunUserExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_user" {
			continue
		}

		response, _, err := client.Client.UsersAPI.UsersList(context.TODO()).Id(rs.Primary.ID).Execute()

		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("user doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunUserDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_user" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			response, _, err := client.Client.UsersAPI.UsersList(context.TODO()).Id(rs.Primary.ID).Execute()

			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("user still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("user still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
