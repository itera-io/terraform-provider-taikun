package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunProjectUserAttachmentConfig = `
resource "taikun_user" "foo" {
  user_name = "%s"
  email     = "%s"
  role      = "User"
}

resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
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
					projectName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
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
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project_user_attachment" {
			continue
		}

		projectId, userId, err := parseProjectUserAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		response, _, err := apiClient.Client.UsersAPI.UsersList(context.TODO()).Id(userId).Execute()
		if err != nil {
			return err
		}
		if len(response.GetData()) != 1 {
			return fmt.Errorf("user with ID %s not found", userId)
		}

		rawUser := response.GetData()[0]

		for _, e := range rawUser.BoundProjects {
			if e.GetProjectId() == projectId {
				return nil
			}
		}

		return fmt.Errorf("project_user_attachment doesn't exist (id = %s)", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTaikunProjectUserAttachmentDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project_user_attachment" {
			continue
		}

		projectId, userId, err := parseProjectUserAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			response, _, err := apiClient.Client.UsersAPI.UsersList(context.TODO()).Id(userId).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 1 {
				return nil
			}

			rawUser := response.GetData()[0]

			for _, e := range rawUser.BoundProjects {
				if e.GetProjectId() == projectId {
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
