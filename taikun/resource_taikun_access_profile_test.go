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

const testAccResourceTaikunAccessProfileConfig = `
resource "taikun_access_profile" "foo" {
  name            = "%s"
  lock       = %t

  ssh_user {
    name       = "oui_oui"
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

  allowed_host {
    description = "Host A"
    address = "10.0.0.0"
    mask_bits = 24
  }

  allowed_host {
    description = "Host B"
    address = "10.1.0.0"
    mask_bits = 24
  }

  allowed_host {
    description = "Host C"
    address = "172.19.42.0"
    mask_bits = 24
  }
}
`

func TestAccResourceTaikunAccessProfile(t *testing.T) {
	name := randomTestName()
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, name, unlocked),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui_oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.#", "3"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.description", "Host A"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.address", "10.0.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.description", "Host B"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.address", "10.1.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.description", "Host C"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.address", "172.19.42.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.mask_bits", "24"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
			{
				ResourceName:      "taikun_access_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunAccessProfileLock(t *testing.T) {
	name := randomTestName()
	const locked = true
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, name, unlocked),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui_oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.#", "3"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.description", "Host A"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.address", "10.0.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.description", "Host B"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.address", "10.1.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.description", "Host C"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.address", "172.19.42.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.mask_bits", "24"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, name, locked),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui_oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.#", "3"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.description", "Host A"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.address", "10.0.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.description", "Host B"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.address", "10.1.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.description", "Host C"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.address", "172.19.42.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.mask_bits", "24"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
		},
	})
}

const testAccResourceTaikunAccessProfileConfigUpdate = `
resource "taikun_access_profile" "foo" {
  name            = "%s"
  lock       = %t

  ssh_user {
    name       = "oui_oui"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
  }

  ssh_user {
    name       = "non_non"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :non_non:"
  }


  ntp_server {
    address = "time.apple.com"
  }

  dns_server {
    address = "1.1.1.1"
  }

  allowed_host {
    description = "Host A"
    address = "192.168.1.0"
    mask_bits = 24
  }
}
`

func TestAccResourceTaikunAccessProfileUpdate(t *testing.T) {
	name := randomTestName()
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, name, unlocked),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui_oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.#", "3"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.description", "Host A"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.address", "10.0.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.description", "Host B"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.address", "10.1.0.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.1.mask_bits", "24"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.description", "Host C"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.address", "172.19.42.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.2.mask_bits", "24"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfigUpdate, name, unlocked),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.apple.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui_oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.1.name", "non_non"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.1.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :non_non:"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.description", "Host A"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.address", "192.168.1.0"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "allowed_host.0.mask_bits", "24"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunAccessProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := client.Client.AccessProfilesAPI.AccessprofilesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("access profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunAccessProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := client.Client.AccessProfilesAPI.AccessprofilesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return resource.RetryableError(errors.New("access profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("access profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
