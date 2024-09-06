package testing

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunCloudCredentialZadaraConfig = `
resource "taikun_cloud_credential_zadara" "foo" {
  name = "%s"
  lock = %t
}

data "taikun_cloud_credential_zadara" "foo" {
  id = resource.taikun_cloud_credential_zadara.foo.id
}
`

func TestAccDataSourceTaikunCloudCredentialZadara(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunCloudCredentialZadaraConfig,
					cloudCredentialName,
					false,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceStateWithIgnores(
					"data.taikun_cloud_credential_zadara.foo",
					"taikun_cloud_credential_zadara.foo",
					map[string]struct{}{
						"access_key_id":     {},
						"secret_access_key": {},
					},
				),
			},
		},
	})
}
