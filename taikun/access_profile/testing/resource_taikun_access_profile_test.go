package testing

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunAccessProfileConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_access_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

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
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	name := utils.RandomTestName()
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, organizationName, organizationFullName, name, unlocked),
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
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	name := utils.RandomTestName()
	const locked = true
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, organizationName, organizationFullName, name, unlocked),
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
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, organizationName, organizationFullName, name, locked),
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
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = 42
}

resource "taikun_access_profile" "foo" {
  organization_id = resource.taikun_organization.foo.id

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
	organizationName := utils.RandomTestName()
	organizationFullName := utils.RandomTestName()
	name := utils.RandomTestName()
	const unlocked = false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, organizationName, organizationFullName, name, unlocked),
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
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfigUpdate, organizationName, organizationFullName, name, unlocked),
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

// --- Trusted registry tests (disabled) ---
// TestAccResourceTaikunAccessProfileTrustedRegistries is kept commented out pending
// platform hotfix / API availability for trusted_registry. Uncomment when re-enabled.
//
// TestAccResourceTaikunAccessProfileTrustedRegistries checks the trusted_registry feature.
//const testAccResourceTaikunAccessProfileTrustedRegistriesConfig = `
//resource "taikun_organization" "bar" {
//  name = "%s"
//  full_name = "%s"
//  discount_rate = 42
//}
//
//resource "taikun_access_profile" "bar" {
//  organization_id = resource.taikun_organization.bar.id
//
//  name = "%s"
//  lock = false
//
//  trusted_registry {
//    registry = "ghcr.io"
//  }
//
//  trusted_registry {
//    registry = "quay.io"
//  }
//}
//`
//
//func TestAccResourceTaikunAccessProfileTrustedRegistries(t *testing.T) {
//	organizationName := utils.RandomTestName()
//	organizationFullName := utils.RandomTestName()
//	name := utils.RandomTestName()
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
//		ProviderFactories: utils_testing.TestAccProviderFactories,
//		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileTrustedRegistriesConfig, organizationName, organizationFullName, name),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					testAccCheckTaikunAccessProfileExists,
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "name", name),
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "lock", "false"),
//					testAccCheckTaikunAccessProfileTrustedRegistries("taikun_access_profile.bar", "ghcr.io", "quay.io"),
//					// Assert that other optional lists are empty by default
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "dns_server.#", "0"),
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "ntp_server.#", "0"),
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "ssh_user.#", "0"),
//					resource.TestCheckResourceAttr("taikun_access_profile.bar", "allowed_host.#", "0"),
//					resource.TestCheckResourceAttrSet("taikun_access_profile.bar", "organization_id"),
//				),
//			},
//			{
//				ResourceName:      "taikun_access_profile.bar",
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

// testAccCheckTaikunAccessProfileTrustedRegistries verifies trusted_registry values
// regardless of list order. Used by TestAccResourceTaikunAccessProfileTrustedRegistries when uncommented.
//func testAccCheckTaikunAccessProfileTrustedRegistries(resourceName string, expected ...string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[resourceName]
//		if !ok {
//			return fmt.Errorf("resource %s not found in state", resourceName)
//		}
//
//		countAttr := rs.Primary.Attributes["trusted_registry.#"]
//		if countAttr != fmt.Sprint(len(expected)) {
//			return fmt.Errorf("%s: trusted_registry.# is %s, want %d", resourceName, countAttr, len(expected))
//		}
//
//		found := make(map[string]struct{}, len(expected))
//		for i := range expected {
//			registry := rs.Primary.Attributes[fmt.Sprintf("trusted_registry.%d.registry", i)]
//			found[registry] = struct{}{}
//		}
//
//		for _, registry := range expected {
//			if _, ok := found[registry]; !ok {
//				return fmt.Errorf("%s: trusted_registry missing %q", resourceName, registry)
//			}
//		}
//
//		return nil
//	}
//}

func testAccCheckTaikunAccessProfileExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.AccessProfilesAPI.AccessprofilesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("access profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunAccessProfileDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_access_profile" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.AccessProfilesAPI.AccessprofilesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("access profile still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("access profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
