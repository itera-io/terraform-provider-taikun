package taikun

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/itera-io/taikungoclient/client/users"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourceTaikunProjectUserAttachmentConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "User"
}

resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
}

resource "taikun_project_user_attachment" "foo" {
  project_id = resource.taikun_project.foo.id
  user_id    = resource.taikun_user.foo.id
}
`

func TestAccResourceTaikunProjectUserAttachment(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	userName := randomTestName()
	userEmail := randomEmail()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectUserAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectUserAttachmentConfig,
					userName,
					userEmail,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunProjectUserAttachmentExists,
					resource.TestCheckResourceAttrSet("taikun_project_user_attachment.foo", "project_id"),
					resource.TestCheckResourceAttr("taikun_project_user_attachment.foo", "project_name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project_user_attachment.foo", "user_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunProjectUserAttachmentExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project_user_attachment" {
			continue
		}

		projectId, userId, err := parseProjectUserAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		params := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
		response, err := apiClient.client.Users.UsersList(params, apiClient)
		if err != nil {
			return err
		}
		if len(response.Payload.Data) != 1 {
			return fmt.Errorf("user with ID %s not found", userId)
		}

		rawUser := response.GetPayload().Data[0]

		for _, e := range rawUser.BoundProjects {
			if e.ProjectID == projectId {
				return nil
			}
		}

		return fmt.Errorf("project_user_attachment doesn't exist (id = %s)", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTaikunProjectUserAttachmentDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project_user_attachment" {
			continue
		}

		projectId, userId, err := parseProjectUserAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		retryErr := resource.Retry(getReadAfterOpTimeout(false), func() *resource.RetryError {
			params := users.NewUsersListParams().WithV(ApiVersion).WithID(&userId)
			response, err := apiClient.client.Users.UsersList(params, apiClient)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 1 {
				return resource.NonRetryableError(errors.New(fmt.Sprintf("user with ID %s not found", userId)))
			}

			rawUser := response.GetPayload().Data[0]

			for _, e := range rawUser.BoundProjects {
				if e.ProjectID == projectId {
					return resource.RetryableError(errors.New("project_user_attachment still exists"))
				}
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("project_user_attachment still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
