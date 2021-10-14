package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/users"
	"github.com/itera-io/taikungoclient/models"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_user", &resource.Sweeper{
		Name: "taikun_user",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := users.NewUsersListParams().WithV(ApiVersion)

			var userList []*models.UserForListDto
			for {
				response, err := apiClient.client.Users.UsersList(params, apiClient)
				if err != nil {
					return err
				}
				userList = append(userList, response.GetPayload().Data...)
				if len(userList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(userList))
				params = params.WithOffset(&offset)
			}

			for _, e := range userList {
				if strings.HasPrefix(e.Username, testNamePrefix) {
					params := users.NewUsersDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, _, err = apiClient.client.Users.UsersDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunUserConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "%s"

  display_name        = "%s"
  user_disabled       = %t
  approved_by_partner = %t
}
`

func TestAccResourceTaikunUser(t *testing.T) {
	userName := randomTestName()
	email := randomString() + "@" + randomString() + ".fr"
	role := "User"
	displayName := randomTestName()
	userDisabled := false
	approvedByPartner := false

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
					userDisabled,
					approvedByPartner,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", role),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_disabled", fmt.Sprint(userDisabled)),
					resource.TestCheckResourceAttr("taikun_user.foo", "approved_by_partner", fmt.Sprint(approvedByPartner)),
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

func TestAccResourceTaikunUserUpdate(t *testing.T) {
	userName := randomTestName()
	email := randomString() + "@" + randomString() + ".fr"
	role := "Manager"
	displayName := randomTestName()
	userDisabled := false
	approvedByPartner := true
	newUserName := randomTestName()
	newEmail := randomString() + "@" + randomString() + ".fr"
	newRole := "Manager"
	newDisplayName := randomTestName()
	newUserDisabled := true
	newApprovedByPartner := false

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
					userDisabled,
					approvedByPartner,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", userName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", role),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", displayName),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_disabled", fmt.Sprint(userDisabled)),
					resource.TestCheckResourceAttr("taikun_user.foo", "approved_by_partner", fmt.Sprint(approvedByPartner)),
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
					newUserDisabled,
					newApprovedByPartner,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunUserExists,
					resource.TestCheckResourceAttr("taikun_user.foo", "user_name", newUserName),
					resource.TestCheckResourceAttr("taikun_user.foo", "email", newEmail),
					resource.TestCheckResourceAttr("taikun_user.foo", "role", newRole),
					resource.TestCheckResourceAttr("taikun_user.foo", "display_name", newDisplayName),
					resource.TestCheckResourceAttr("taikun_user.foo", "user_disabled", fmt.Sprint(newUserDisabled)),
					resource.TestCheckResourceAttr("taikun_user.foo", "approved_by_partner", fmt.Sprint(newApprovedByPartner)),
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
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_user" {
			continue
		}

		params := users.NewUsersListParams().WithV(ApiVersion).WithID(&rs.Primary.ID)

		response, err := client.client.Users.UsersList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("user doesn't exist")
		}
	}

	return nil
}

func testAccCheckTaikunUserDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_user" {
			continue
		}

		params := users.NewUsersListParams().WithV(ApiVersion).WithID(&rs.Primary.ID)

		response, err := client.client.Users.UsersList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("user still exists")
		}
	}

	return nil
}
