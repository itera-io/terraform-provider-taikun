package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"testing"
)

const testAccResourceTaikunAccessProfile = `
resource "taikun_access_profile" "foo" {
  name            = "%s"
  organization_id = "441"
  #is_locked       = %t

  ssh_user {
    name       = "oui oui"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
  }

  ntp_server {
    address = "time.windows.com"
  }

  ntp_server {
    address = "ntp.pool.org"
  }

  dns_server {
    address = "8.8.8.8"
  }

  dns_server {
    address = "8.8.4.4"
  }
}
`

func TestAccResourceTaikunAccessProfile(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfile, firstName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", firstName),
					//resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "true"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "organization_id", "441"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
				),
			},
		},
	})
}

func TestAccResourceTaikunAccessProfileRename(t *testing.T) {
	firstName := randomTestName()
	secondName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfile, firstName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", firstName),
					//resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "true"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "organization_id", "441"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfile, secondName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", secondName),
					//resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "true"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "organization_id", "441"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
				),
			},
		},
	})
}

func testAccCheckTaikunAccessProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.AccessProfiles.AccessProfilesList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("access profile doesn't exists")
		}
	}

	return nil
}

func testAccCheckTaikunAccessProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.AccessProfiles.AccessProfilesList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("access profile still exists")
		}
	}

	return nil
}
