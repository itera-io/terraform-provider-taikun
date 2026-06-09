package testing

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunUserConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  global_role      = "%s"

  display_name        = "%s"
}
`

func TestAccResourceTaikunUser(t *testing.T) {
	userName := utils.RandomTestName()
	email := utils.RandomEmail()
	globalRole := "Admin"
	displayName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunUserDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					userName,
					email,
					globalRole,
					displayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists(t),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "global_role", globalRole),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
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
	globalRole := "Admin"
	displayName := utils.RandomTestName()
	newUserName := utils.RandomTestName()
	newEmail := utils.RandomEmail()
	newGlobalRole := "AccountOwner"
	newDisplayName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunUserDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					userName,
					email,
					globalRole,
					displayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists(t),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "global_role", globalRole),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_notification_enabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunUserConfig,
					newUserName,
					newEmail,
					newGlobalRole,
					newDisplayName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunUserExists(t),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", newUserName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", newEmail),
					resource.TestCheckResourceAttr("taikun_user.foo", "global_role", newGlobalRole),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", newDisplayName),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "id"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "is_owner"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_confirmed"),
					resource.TestCheckResourceAttrSet("taikun_user.foo", "email_notification_enabled"),
				),
			},
		},
	})
}

func testAccCheckTaikunUserExists(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_user" {
				continue
			}

			searchBody := tkcore.UsersSearchCommand{}
			searchBody.SetSearchTerm(rs.Primary.ID)
			searchRes, _, err := client.Client.SearchAPI.SearchUsers(t.Context()).UsersSearchCommand(searchBody).Execute()

			if err != nil {
				return fmt.Errorf("user doesn't exist (id = %s)", rs.Primary.ID)
			}

			found := false
			for _, u := range searchRes.GetData() {
				if u.GetId() == rs.Primary.ID {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("user doesn't exist (id = %s)", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckTaikunUserDestroy(t *testing.T) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := utils_testing.TestAccProvider.Meta().(*tk.Client)

		for _, rs := range state.RootModule().Resources {
			if rs.Type != "taikun_user" {
				continue
			}

			retryErr := retry.RetryContext(t.Context(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
				searchBody := tkcore.UsersSearchCommand{}
				searchBody.SetSearchTerm(rs.Primary.ID)
				searchRes, _, err := client.Client.SearchAPI.SearchUsers(t.Context()).UsersSearchCommand(searchBody).Execute()

				if err != nil {
					return retry.NonRetryableError(err)
				}

				for _, u := range searchRes.GetData() {
					if u.GetId() == rs.Primary.ID {
						return retry.RetryableError(errors.New("user still exists"))
					}
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
}
