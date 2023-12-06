package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunStandaloneProfilesConfig = `
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

data "taikun_standalone_profiles" "all" {
  depends_on = [
    taikun_standalone_profile.foo
  ]
}`

func TestAccDataSourceTaikunStandaloneProfiles(t *testing.T) {
	standaloneProfileName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunStandaloneProfilesConfig, standaloneProfileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_standalone_profiles.all", "id", "all"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.public_key"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunStandaloneProfilesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
    organization_id = resource.taikun_organization.foo.id
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

data "taikun_standalone_profiles" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_standalone_profile.foo
  ]
}`

func TestAccDataSourceTaikunStandaloneProfilesWithFilter(t *testing.T) {
	organizationName := randomTestName()
	organizationFullName := randomTestName()
	standaloneProfileName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunStandaloneProfilesWithFilterConfig, organizationName, organizationFullName, standaloneProfileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_standalone_profiles.all", "standalone_profiles.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.#"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.lock"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.organization_name"),
					resource.TestCheckResourceAttrSet("data.taikun_standalone_profiles.all", "standalone_profiles.0.public_key"),
				),
			},
		},
	})
}
