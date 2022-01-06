package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunStandaloneProfileConfig = `
resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"

    security_group {
        name = "http"
        from_port = 80
        to_port = 80
        ip_protocol = "TCP"
        cidr = "0.0.0.0/0"
    }
    security_group {
        name = "https"
        from_port = 443
        to_port = 443
        ip_protocol = "TCP"
        cidr = "0.0.0.0/0"
    }
}

data "taikun_standalone_profile" "foo" {
  id = resource.taikun_standalone_profile.foo.id
}
`

func TestAccDataSourceTaikunStandaloneProfile(t *testing.T) {
	standaloneProfile := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunStandaloneProfileConfig, standaloneProfile),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_standalone_profile.foo",
					"taikun_standalone_profile.foo",
				),
			},
		},
	})
}
