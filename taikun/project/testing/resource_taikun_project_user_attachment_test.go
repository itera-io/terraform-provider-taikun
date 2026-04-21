package testing

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
)

const testAccResourceTaikunProjectUserAttachmentConfig = `
resource "taikun_project_user_attachment" "foo" {
  project_id = "dummy-project-id"
  user_id    = "dummy-user-id"
}
`

func TestAccResourceTaikunProjectUserAttachment_Deprecated(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceTaikunProjectUserAttachmentConfig,
				ExpectError: regexp.MustCompile("The taikun_project_user_attachment resource is deprecated and no longer supported by the Taikun API."),
			},
		},
	})
}
