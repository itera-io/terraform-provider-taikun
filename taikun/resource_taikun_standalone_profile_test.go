package taikun

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/stand_alone_profile"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_standalone_profile", &resource.Sweeper{
		Name:         "taikun_standalone_profile",
		Dependencies: []string{"taikun_project"},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion)

			var standaloneProfilesList []*models.StandAloneProfilesListDto
			for {
				response, err := apiClient.client.StandAloneProfile.StandAloneProfileList(params, apiClient)
				if err != nil {
					return err
				}
				standaloneProfilesList = append(standaloneProfilesList, response.GetPayload().Data...)
				if len(standaloneProfilesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(standaloneProfilesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range standaloneProfilesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					body := &models.DeleteStandAloneProfileCommand{ID: e.ID}
					params := stand_alone_profile.NewStandAloneProfileDeleteParams().WithV(ApiVersion).WithBody(body)
					_, err = apiClient.client.StandAloneProfile.StandAloneProfileDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunStandaloneProfileConfig = `
resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
    lock = %t
}
`

func TestAccResourceTaikunStandaloneProfile(t *testing.T) {
	name := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
				),
			},
			{
				ResourceName:      "taikun_standalone_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunStandaloneProfileLock(t *testing.T) {
	name := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
				),
			},
		},
	})
}

func TestAccResourceTaikunStandaloneProfileRename(t *testing.T) {
	name := randomTestName()
	newName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, newName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
				),
			},
		},
	})
}

func testAccCheckTaikunStandaloneProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_standalone_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.StandAloneProfile.StandAloneProfileList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("standalone profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunStandaloneProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_standalone_profile" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.StandAloneProfile.StandAloneProfileList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("standalone profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("standalone profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
