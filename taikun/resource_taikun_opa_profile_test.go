package taikun

import (
	"fmt"
	"github.com/itera-io/taikungoclient/client/opa_profiles"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_opa_profile", &resource.Sweeper{
		Name:         "taikun_opa_profile",
		Dependencies: []string{"taikun_project"},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion)

			var OPAProfilesList []*models.OpaProfileListDto
			for {
				response, err := apiClient.client.OpaProfiles.OpaProfilesList(params, apiClient)
				if err != nil {
					return err
				}
				OPAProfilesList = append(OPAProfilesList, response.GetPayload().Data...)
				if len(OPAProfilesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(OPAProfilesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range OPAProfilesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := opa_profiles.NewOpaProfilesDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteOpaProfileCommand{ID: e.ID})
					_, err = apiClient.client.OpaProfiles.OpaProfilesDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunOPAProfileConfig = `
resource "taikun_opa_profile" "foo" {
  name = "%s"
  lock = %t

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t
}
`

func TestAccResourceTaikunOPAProfile(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOPAProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceTaikunOPAProfileConfig,
					firstName,
					false,
					false,
					false,
					false,
					false,
					false,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOPAProfileExists,
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_node_port", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_http_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "require_probe", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_service_selector", "false"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_name"),
				),
			},
			{
				ResourceName:      "taikun_opa_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunOPAProfileLock(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOPAProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceTaikunOPAProfileConfig,
					firstName,
					false,
					true,
					true,
					true,
					true,
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOPAProfileExists,
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_http_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunOPAProfileConfig,
					firstName,
					true,
					true,
					true,
					true,
					true,
					true,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOPAProfileExists,
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "forbid_http_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_opa_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_opa_profile.foo", "organization_name"),
				),
			},
		},
	})
}

func testAccCheckTaikunOPAProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_opa_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.OpaProfiles.OpaProfilesList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("opa profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunOPAProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_opa_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := opa_profiles.NewOpaProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.OpaProfiles.OpaProfilesList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("opa profile still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
