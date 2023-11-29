package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	userName := randomTestName()
	email := randomEmail()
	role := "User"
	displayName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	userName := randomTestName()
	email := randomEmail()
	role := "Manager"
	displayName := randomTestName()
	newUserName := randomTestName()
	newEmail := randomEmail()
	newRole := "Manager"
	newDisplayName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	client := testAccProvider.Meta().(*tk.Client)

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
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_user" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			response, _, err := client.Client.UsersAPI.UsersList(context.TODO()).Id(rs.Primary.ID).Execute()

			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return resource.RetryableError(errors.New("user still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("user still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
