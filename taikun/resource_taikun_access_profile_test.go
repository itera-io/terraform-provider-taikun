package taikun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/access_profiles"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_access_profile", &resource.Sweeper{
		Name: "taikun_access_profile",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := access_profiles.NewAccessProfilesListParams().WithV(ApiVersion)

			var accessProfilesList []*models.AccessProfilesListDto

			for {
				response, err := apiClient.client.AccessProfiles.AccessProfilesList(params, apiClient)
				if err != nil {
					return err
				}
				accessProfilesList = append(accessProfilesList, response.GetPayload().Data...)
				if len(accessProfilesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(accessProfilesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range accessProfilesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := access_profiles.NewAccessProfilesDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, _, err = apiClient.client.AccessProfiles.AccessProfilesDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunAccessProfileConfig = `
resource "taikun_access_profile" "foo" {
  name            = "%s"
  is_locked       = %t

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
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, firstName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
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

func TestAccResourceTaikunAccessProfileRenameAndLock(t *testing.T) {
	firstName := randomTestName()
	secondName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunAccessProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, firstName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "false"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunAccessProfileConfig, secondName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunAccessProfileExists,
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "name", secondName),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "is_locked", "true"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.0.address", "8.8.8.8"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "dns_server.1.address", "8.8.4.4"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.#", "2"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.0.address", "time.windows.com"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ntp_server.1.address", "ntp.pool.org"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.#", "1"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.name", "oui oui"),
					resource.TestCheckResourceAttr("taikun_access_profile.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
					resource.TestCheckResourceAttrSet("taikun_access_profile.foo", "organization_id"),
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
			return fmt.Errorf("access profile doesn't exist (id = %s)", rs.Primary.ID)
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
			return fmt.Errorf("access profile still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
